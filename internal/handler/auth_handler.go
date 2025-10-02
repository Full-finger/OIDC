// internal/handler/auth_handler.go

package handler

import (
	"html/template"
	"net/http"
	"net/url"

	"github.com/Full-finger/OIDC/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler 定义认证相关处理函数
type AuthHandler struct {
	oauthService service.OAuthService
}

// NewAuthHandler 创建一个新的AuthHandler实例
func NewAuthHandler(oauthService service.OAuthService) *AuthHandler {
	return &AuthHandler{
		oauthService: oauthService,
	}
}

// AuthorizeHandler 处理GET /oauth/authorize请求
func (h *AuthHandler) AuthorizeHandler(c *gin.Context) {
	// 获取查询参数
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	scope := c.Query("scope")
	responseType := c.Query("response_type")
	state := c.Query("state")
	codeChallenge := c.Query("code_challenge")
	codeChallengeMethod := c.Query("code_challenge_method")

	// 验证必需参数
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id is required"})
		return
	}

	if redirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "redirect_uri is required"})
		return
	}

	if responseType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "response_type is required"})
		return
	}

	// 验证response_type是否为"code"
	if responseType != "code" {
		redirectWithError(c, redirectURI, "unsupported_response_type", state)
		return
	}

	// 验证客户端是否存在
	client, err := h.oauthService.FindClientByClientID(c.Request.Context(), clientID)
	if err != nil {
		redirectWithError(c, redirectURI, "invalid_client", state)
		return
	}

	// 验证重定向URI是否在允许列表中
	if !isValidRedirectURI(redirectURI, client.RedirectURIs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid redirect_uri"})
		return
	}

	// 验证scope是否被客户端允许
	scopes := []string{} // 默认scope
	if scope != "" {
		scopes = append(scopes, scope)
	}

	// 检查用户是否登录（这里简化实现，实际项目中应该检查session或JWT token）
	// 在这个示例中，我们假设用户未登录，重定向到登录页面
	// 在实际应用中，你需要根据你的认证机制来检查用户是否已登录
	userID, isLoggedIn := h.checkUserLoggedIn(c)
	if !isLoggedIn {
		// 重定向到登录页面，登录成功后再返回授权页面
		loginURL := "/login?redirect=/oauth/authorize?" + c.Request.URL.RawQuery
		c.Redirect(http.StatusFound, loginURL)
		return
	}

	// 用户已登录，渲染授权页面
	// 这里我们直接返回一个简单的HTML页面
	html := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>授权确认</title>
</head>
<body>
    <h2>授权请求</h2>
    <p>应用 <strong>{{.ClientName}}</strong> 请求访问您的账户。</p>
    <p>请求的权限: {{.Scope}}</p>
    
    <form method="POST" action="/oauth/authorize">
        <input type="hidden" name="client_id" value="{{.ClientID}}">
        <input type="hidden" name="redirect_uri" value="{{.RedirectURI}}">
        <input type="hidden" name="scope" value="{{.Scope}}">
        <input type="hidden" name="response_type" value="{{.ResponseType}}">
        <input type="hidden" name="state" value="{{.State}}">
        <input type="hidden" name="code_challenge" value="{{.CodeChallenge}}">
        <input type="hidden" name="code_challenge_method" value="{{.CodeChallengeMethod}}">
        <input type="hidden" name="user_id" value="{{.UserID}}">
        
        <button type="submit" name="approve" value="yes">同意</button>
        <button type="submit" name="approve" value="no">拒绝</button>
    </form>
</body>
</html>
`
	
	tmpl, err := template.New("authorize").Parse(html)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to render authorization page"})
		return
	}
	
	data := map[string]interface{}{
		"ClientID":             clientID,
		"ClientName":           client.Name,
		"RedirectURI":          redirectURI,
		"Scope":                scope,
		"ResponseType":         responseType,
		"State":                state,
		"CodeChallenge":        codeChallenge,
		"CodeChallengeMethod":  codeChallengeMethod,
		"UserID":               userID,
	}
	
	c.Header("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(c.Writer, data)
}

// AuthorizePostHandler 处理POST /oauth/authorize请求
func (h *AuthHandler) AuthorizePostHandler(c *gin.Context) {
	// 获取表单参数
	clientID := c.PostForm("client_id")
	redirectURI := c.PostForm("redirect_uri")
	scope := c.PostForm("scope")
	responseType := c.PostForm("response_type")
	state := c.PostForm("state")
	approve := c.PostForm("approve")
	codeChallenge := c.PostForm("code_challenge")
	codeChallengeMethod := c.PostForm("code_challenge_method")
	
	// 从上下文获取用户ID（实际应用中应通过认证中间件获取）
	userID, exists := c.Get("userID")
	var realUserID int64
	
	if exists {
		realUserID = userID.(int64)
	} else {
		// 如果上下文中没有用户ID，使用表单中的用户ID（仅用于演示）
		// userID := c.PostForm("user_id")
		// 在实际应用中，这里应该返回错误，因为用户未正确登录
		realUserID = 1 // 默认用户ID，仅用于演示
	}

	// 验证必需参数
	if clientID == "" || redirectURI == "" || responseType == "" {
		redirectWithError(c, redirectURI, "invalid_request", state)
		return
	}

	// 验证用户是否拒绝授权
	if approve != "yes" {
		redirectWithError(c, redirectURI, "access_denied", state)
		return
	}

	// 验证客户端是否存在
	client, err := h.oauthService.FindClientByClientID(c.Request.Context(), clientID)
	if err != nil {
		redirectWithError(c, redirectURI, "invalid_client", state)
		return
	}

	// 验证重定向URI是否在允许列表中
	if !isValidRedirectURI(redirectURI, client.RedirectURIs) {
		redirectWithError(c, redirectURI, "invalid_request", state)
		return
	}

	// 生成授权码
	scopes := []string{}
	if scope != "" {
		scopes = append(scopes, scope)
	}
	
	// 处理PKCE参数
	var codeChallengePtr *string
	var codeChallengeMethodPtr *string
	
	if codeChallenge != "" {
		codeChallengePtr = &codeChallenge
	}
	
	if codeChallengeMethod != "" {
		codeChallengeMethodPtr = &codeChallengeMethod
	}

	authCode, err := h.oauthService.GenerateAuthorizationCode(
		c.Request.Context(),
		clientID,
		realUserID,
		redirectURI,
		scopes,
		codeChallengePtr,
		codeChallengeMethodPtr,
	)

	if err != nil {
		redirectWithError(c, redirectURI, "server_error", state)
		return
	}

	// 重定向回客户端，携带授权码
	redirectWithCode(c, redirectURI, authCode.Code, state)
}

// checkUserLoggedIn 检查用户是否已登录（简化实现）
func (h *AuthHandler) checkUserLoggedIn(c *gin.Context) (int64, bool) {
	// 在实际应用中，你应该检查session或JWT token来确定用户是否已登录
	// 这里我们简化实现，假设用户未登录
	// 你可以根据实际的认证机制来实现这个方法
	
	// 示例实现（如果你使用JWT中间件）：
	// userID, exists := c.Get("userID")
	// if !exists {
	//     return 0, false
	// }
	// return userID.(int64), true
	
	return 0, false // 简化实现，始终返回未登录
}

// isValidRedirectURI 验证重定向URI是否在允许列表中
func isValidRedirectURI(redirectURI string, allowedRedirectURIs []string) bool {
	for _, uri := range allowedRedirectURIs {
		if uri == redirectURI {
			return true
		}
	}
	return false
}

// redirectWithError 重定向并携带错误信息
func redirectWithError(c *gin.Context, redirectURI, errorType, state string) {
	u, err := url.Parse(redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid redirect_uri"})
		return
	}

	params := u.Query()
	params.Set("error", errorType)
	if state != "" {
		params.Set("state", state)
	}
	u.RawQuery = params.Encode()

	c.Redirect(http.StatusFound, u.String())
}

// redirectWithCode 重定向并携带授权码
func redirectWithCode(c *gin.Context, redirectURI, code, state string) {
	u, err := url.Parse(redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid redirect_uri"})
		return
	}

	params := u.Query()
	params.Set("code", code)
	if state != "" {
		params.Set("state", state)
	}
	u.RawQuery = params.Encode()

	c.Redirect(http.StatusFound, u.String())
}