package router

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
	"github.com/Full-finger/OIDC/internal/handler"
	"github.com/Full-finger/OIDC/internal/service"
	"github.com/Full-finger/OIDC/internal/repository"
	"github.com/Full-finger/OIDC/internal/helper"
	"github.com/Full-finger/OIDC/internal/util"
	"github.com/Full-finger/OIDC/internal/middleware"
	"github.com/Full-finger/OIDC/internal/mapper"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 初始化数据库连接
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// 如果数据库连接失败，使用内存模式继续运行
		fmt.Printf("警告: 无法连接到数据库: %v\n", err)
		db = nil
	}

	// 初始化依赖
	var userRepo repository.UserRepository
	var userMapper mapper.UserMapper
	
	if db != nil {
		userMapper = mapper.NewUserMapper(db)
		userRepo = repository.NewUserRepository(userMapper)
	} else {
		// 使用内存存储
		userRepo = repository.NewUserRepository(nil)
	}
	
	userHelper := helper.NewUserHelper()
	tokenRepo := repository.NewVerificationTokenRepository()
	emailQueue := util.NewSimpleEmailQueue()
	
	userService := service.NewUserService(userRepo, userHelper, tokenRepo, emailQueue)
	userHandler := handler.NewUserHandler(userService)
	verificationHandler := handler.NewVerificationHandler(userService)

	// 初始化OAuth依赖
	oauthService := service.NewOAuthService()
	oauthHandler := handler.NewOAuthHandler(oauthService)

	// 初始化番剧收藏依赖
	animeRepo := repository.NewAnimeRepository()
	animeService := service.NewAnimeService(animeRepo)
	animeHandler := handler.NewAnimeHandler(animeService)

	// 初始化收藏依赖
	collectionRepo := repository.NewCollectionRepository()
	collectionService := service.NewCollectionService(collectionRepo, animeRepo)
	collectionHandler := handler.NewCollectionHandler(collectionService)

	// 初始化Bangumi依赖
	bangumiRepo := repository.NewBangumiRepository()
	bangumiService := service.NewBangumiService(bangumiRepo, animeRepo, collectionRepo)
	bangumiHandler := handler.NewBangumiHandler(bangumiService)

	// 初始化中间件
	rateLimiter := middleware.NewRateLimiter()

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 用户相关路由
		v1.POST("/register", rateLimiter.LimitByIP(), userHandler.Register)
		v1.POST("/resend-verification", rateLimiter.LimitByUser(), userHandler.ResendVerificationEmail)
		v1.POST("/login", userHandler.Login)
		// 邮箱验证路由
		v1.GET("/verify", verificationHandler.VerifyEmail)
		
		// 番剧收藏路由
		anime := v1.Group("/anime")
		{
			anime.GET("/:id", animeHandler.GetAnimeByIDHandler)
			anime.GET("/search", animeHandler.SearchAnimesHandler)
			anime.GET("/list", animeHandler.ListAnimesHandler)
			anime.GET("/status", animeHandler.ListAnimesByStatusHandler)
			// 添加创建、更新和删除番剧的路由
			anime.POST("/", animeHandler.CreateAnimeHandler)
			anime.PUT("/:id", animeHandler.UpdateAnimeHandler)
			anime.DELETE("/:id", animeHandler.DeleteAnimeHandler)
		}
		
		collection := v1.Group("/collection")
		{
			collection.Use(middleware.JWTAuthMiddleware())
			collection.POST("/", collectionHandler.AddToCollectionHandler)
			collection.GET("/:anime_id", collectionHandler.GetCollectionHandler)
			collection.PUT("/:anime_id", collectionHandler.UpdateCollectionHandler)
			collection.DELETE("/:anime_id", collectionHandler.RemoveFromCollectionHandler)
			collection.GET("/", collectionHandler.ListUserCollectionsHandler)
			collection.GET("/status", collectionHandler.ListUserCollectionsByStatusHandler)
			collection.GET("/favorites", collectionHandler.ListUserFavoritesHandler)
		}
		
		// Bangumi绑定路由
		bangumi := v1.Group("/bangumi")
		{
			bangumi.Use(middleware.JWTAuthMiddleware())
			bangumi.GET("/authorize", bangumiHandler.AuthorizeHandler)
			bangumi.GET("/callback", bangumiHandler.CallbackHandler)
			bangumi.DELETE("/unbind", bangumiHandler.UnbindHandler)
			bangumi.GET("/account", bangumiHandler.GetBoundAccountHandler)
			bangumi.POST("/sync", bangumiHandler.SyncCollectionHandler)
		}
	}

	// OIDC Discovery端点
	r.GET("/.well-known/openid-configuration", oauthHandler.DiscoveryHandler)
	r.GET("/jwks.json", oauthHandler.JWKSHandler)

	// OAuth 2.0 路由
	oauth := r.Group("/oauth")
	{
		// 授权端点
		oauth.GET("/authorize", oauthHandler.AuthorizeHandler)
		// 令牌端点
		oauth.POST("/token", oauthHandler.TokenHandler)
		// 用户信息端点
		oauth.GET("/userinfo", oauthHandler.UserInfoHandler)
	}

	return r
}