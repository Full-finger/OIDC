package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTUtil JWT工具接口
type JWTUtil interface {
	// GenerateIDToken 生成ID Token
	GenerateIDToken(claims *IDTokenClaims) (string, error)
	
	// GenerateAccessToken 生成Access Token
	GenerateAccessToken(claims *AccessTokenClaims) (string, error)
	
	// ParseIDToken 解析ID Token
	ParseIDToken(tokenString string) (*IDTokenClaims, error)
	
	// ParseAccessToken 解析Access Token
	ParseAccessToken(tokenString string) (*AccessTokenClaims, error)
}

// jwtUtil JWT工具实现
type jwtUtil struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
}

// IDTokenClaims ID Token声明
type IDTokenClaims struct {
	jwt.RegisteredClaims
	Nonce   string `json:"nonce,omitempty"`
	Profile string `json:"profile,omitempty"`
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
}

// AccessTokenClaims Access Token声明
type AccessTokenClaims struct {
	jwt.RegisteredClaims
	Scope string `json:"scope,omitempty"`
}

// NewJWTUtil 创建JWT工具实例
func NewJWTUtil() (JWTUtil, error) {
	// 从环境变量获取JWT配置
	issuer := getEnv("JWT_ISSUER", "OIDC")
	
	// 读取私钥文件
	privateKeyPath := getEnv("JWT_PRIVATE_KEY_PATH", "config/private_key.pem")
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	
	// 解析私钥
	privateKey, err := parsePrivateKey(privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	
	// 读取公钥文件
	publicKeyPath := getEnv("JWT_PUBLIC_KEY_PATH", "config/public_key.pem")
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}
	
	// 解析公钥
	publicKey, err := parsePublicKey(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	
	return &jwtUtil{
		privateKey: privateKey,
		publicKey:  publicKey,
		issuer:     issuer,
	}, nil
}

// GenerateIDToken 生成ID Token
func (j *jwtUtil) GenerateIDToken(claims *IDTokenClaims) (string, error) {
	// 设置标准声明
	if claims.Issuer == "" {
		claims.Issuer = j.issuer
	}
	
	if claims.IssuedAt.IsZero() {
		claims.IssuedAt = jwt.NewNumericDate(time.Now())
	}
	
	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour)) // 1小时过期
	}
	
	if claims.Subject == "" {
		claims.Subject = fmt.Sprintf("user:%d", 1) // 示例用户ID
	}
	
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	
	// 签名并生成token字符串
	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign ID token: %w", err)
	}
	
	return tokenString, nil
}

// GenerateAccessToken 生成Access Token
func (j *jwtUtil) GenerateAccessToken(claims *AccessTokenClaims) (string, error) {
	// 设置标准声明
	if claims.Issuer == "" {
		claims.Issuer = j.issuer
	}
	
	if claims.IssuedAt.IsZero() {
		claims.IssuedAt = jwt.NewNumericDate(time.Now())
	}
	
	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour)) // 1小时过期
	}
	
	if claims.Subject == "" {
		claims.Subject = fmt.Sprintf("user:%d", 1) // 示例用户ID
	}
	
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	
	// 签名并生成token字符串
	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	
	return tokenString, nil
}

// ParseIDToken 解析ID Token
func (j *jwtUtil) ParseIDToken(tokenString string) (*IDTokenClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &IDTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID token: %w", err)
	}
	
	// 验证token
	if claims, ok := token.Claims.(*IDTokenClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, fmt.Errorf("invalid ID token")
}

// ParseAccessToken 解析Access Token
func (j *jwtUtil) ParseAccessToken(tokenString string) (*AccessTokenClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}
	
	// 验证token
	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, fmt.Errorf("invalid access token")
}

// parsePrivateKey 解析私钥
func parsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}
	
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// 尝试解析PKCS8格式
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA private key")
		}
		
		return rsaKey, nil
	}
	
	return privateKey, nil
}

// parsePublicKey 解析公钥
func parsePublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		// 尝试解析PKIX格式
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %w", err)
		}
		
		rsaKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA public key")
		}
		
		return rsaKey, nil
	}
	
	return publicKey, nil
}