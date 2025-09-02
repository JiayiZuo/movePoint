package main

import (
	"log"

	"movePoint/internal/database"
	"movePoint/internal/handlers"
	"movePoint/internal/services"
	"movePoint/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// 初始化数据库
	//dsn := "user:password@tcp(127.0.0.1:3306)/climbing_app?charset=utf8mb4&parseTime=True&loc=Local"
	//err = database.InitDB(dsn)
	//if err != nil {
	//	log.Fatal("Failed to initialize database:", err)
	//}

	// 初始化服务
	climbingService := services.NewClimbingService(database.DB)
	analysisService := services.NewAnalysisService(database.DB)
	userService := services.NewUserService(database.DB)

	// 初始化处理器
	climbingHandler := handlers.NewClimbingHandler(climbingService)
	analysisHandler := handlers.NewAnalysisHandler(analysisService)
	userHandler := handlers.NewUserHandler(userService)

	// 设置路由
	router := gin.Default()

	// 公开路由
	router.POST("/login", func(c *gin.Context) {
		// 登录逻辑
	})
	router.POST("/register", func(c *gin.Context) {
		// 注册逻辑
	})

	// 需要认证的路由组
	auth := router.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// 攀岩记录路由
		auth.POST("/records", climbingHandler.CreateRecord)
		auth.GET("/records", climbingHandler.GetRecords)
		auth.GET("/records/:id", climbingHandler.GetRecord)
		auth.PUT("/records/:id", climbingHandler.UpdateRecord)
		auth.DELETE("/records/:id", climbingHandler.DeleteRecord)

		// 分析路由
		auth.GET("/analysis/climbing", analysisHandler.GetClimbingAnalysis)

		// 用户路由 (个人主页)
		auth.GET("/profile", userHandler.GetProfile)
		auth.PUT("/profile", userHandler.UpdateProfile)
		auth.GET("/profile/stats", userHandler.GetStats)
		auth.GET("/profile/achievements", userHandler.GetAchievements)
		auth.POST("/profile/check-achievements", userHandler.CheckAchievements)
	}

	// 启动服务器
	log.Println("Server starting on :8080")
	router.Run(":8080")
}
