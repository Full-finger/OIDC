package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	// 创建config目录（如果不存在）
	if err := os.MkdirAll("config", 0755); err != nil {
		fmt.Printf("Failed to create config directory: %v\n", err)
		return
	}

	// 生成RSA私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Failed to generate private key: %v\n", err)
		return
	}

	// 创建私钥文件
	privateKeyFile, err := os.Create("config/private_key.pem")
	if err != nil {
		fmt.Printf("Failed to create private key file: %v\n", err)
		return
	}
	defer privateKeyFile.Close()

	// 编码私钥为PEM格式
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// 写入私钥文件
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		fmt.Printf("Failed to write private key: %v\n", err)
		return
	}

	// 获取公钥
	publicKey := &privateKey.PublicKey

	// 创建公钥文件
	publicKeyFile, err := os.Create("config/public_key.pem")
	if err != nil {
		fmt.Printf("Failed to create public key file: %v\n", err)
		return
	}
	defer publicKeyFile.Close()

	// 编码公钥为PEM格式
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)

	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	// 写入公钥文件
	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		fmt.Printf("Failed to write public key: %v\n", err)
		return
	}

	fmt.Println("JWT RSA key pair generated successfully!")
	fmt.Println("Private key saved to: config/private_key.pem")
	fmt.Println("Public key saved to: config/public_key.pem")
}