package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"test_task/internal/handler"

	"test_task/internal/repository"
	"test_task/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	logger.Info("Запуск приложения")

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://akbota:4912@localhost:5432/test_db"
	}

	logger.Info("Подключение к базе данных", "dsn", connStr)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		logger.Error("Ошибка парсинга конфигурации подключения", "error", err)
		os.Exit(1)
	}

	poolConfig.MinConns = 2
	poolConfig.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Error("Ошибка подключения к базе данных", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("Ошибка ping базы данных", "error", err)
		os.Exit(1)
	}

	logger.Info("Успешное подключение к базе данных")

	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP запрос",
			"method", param.Method,
			"path", param.Path,
			"status", param.StatusCode,
			"latency", param.Latency,
			"clientIP", param.ClientIP,
		)
		return ""
	}))

	r.Use(gin.Recovery())

	r.RedirectTrailingSlash = false

	r.GET("/routes", func(c *gin.Context) {
		routes := r.Routes()
		c.JSON(http.StatusOK, gin.H{
			"routes": routes,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		logger.Info("Ping endpoint called")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "ok",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		logger.Info("Health check called")
		if err := pool.Ping(context.Background()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"error":  "Database unavailable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Service is running",
		})
	})

	logger.Info("Регистрация маршрутов пользователей")

	api := r.Group("/api")
	{

		users := api.Group("/users")
		{

			users.POST("", func(c *gin.Context) {
				logger.Info("POST /api/users called - регистрация")
				userHandler.CreateUser(c)
			})

			users.POST("/login", func(c *gin.Context) {
				logger.Info("POST /api/users/login called - вход")
				userHandler.Login(c)
			})
			protected := users.Group("")
			protected.Use(handler.AuthMiddleware())
			{

				protected.GET("", func(c *gin.Context) {
					logger.Info("GET /api/users called (protected)")
					userHandler.GetAllUsers(c)
				})

				protected.GET("/:id", func(c *gin.Context) {
					logger.Info("GET /api/users/:id called", "id", c.Param("id"))
					userHandler.GetUser(c)
				})

				protected.PUT("/:id", func(c *gin.Context) {
					logger.Info("PUT /api/users/:id called", "id", c.Param("id"))
					userHandler.UpdateUser(c)
				})

				protected.DELETE("/:id", func(c *gin.Context) {
					logger.Info("DELETE /api/users/:id called", "id", c.Param("id"))
					userHandler.DeleteUser(c)
				})
			}
		}
	}
	r.NoRoute(func(c *gin.Context) {
		logger.Warn("Маршрут не найден",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Endpoint not found",
		})
	})

	for _, route := range r.Routes() {
		logger.Info("Зарегистрирован маршрут",
			"method", route.Method,
			"path", route.Path,
		)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Запуск сервера", "port", port)
	if err := r.Run(":" + port); err != nil {
		logger.Error("Ошибка запуска сервера", "error", err)
		os.Exit(1)
	}
}
