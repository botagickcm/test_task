// @title           User Test
// @version         1.0
// @description     CRUD управления пользователями с JWT аутентификацией

// @contact.name   akbota
// @contact.email  akbota.bolat678@gmail.com

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description     ## 📋 Описание API
// @description
// @description     Это REST API для управления пользователями с JWT аутентификацией.
// @description
// @description     ### 🔐 Аутентификация
// @description     1. Зарегистрируйтесь через `POST /api/users`
// @description     2. Войдите через `POST /api/users/login` чтобы получить JWT токен
// @description     3. Используйте токен в заголовке `Authorization: Bearer <token>`
// @description
// @description     ### 🚀 Быстрый старт
// @description     ```bash
// @description     # Регистрация
// @description     curl -X POST http://localhost:8080/api/users \
// @description       -H "Content-Type: application/json" \
// @description       -d '{"first_name":"Иван","last_name":"Петров","login":"ivan@test.com","password":"pass123"}'
// @description
// @description     # Вход
// @description     curl -X POST http://localhost:8080/api/users/login \
// @description       -H "Content-Type: application/json" \
// @description       -d '{"login":"ivan@test.com","password":"pass123"}'
// @description     ```
// @description
// @description     ### 📊 Коды ответов
// @description     - 200: Успех
// @description     - 201: Создано
// @description     - 400: Неверный запрос
// @description     - 401: Не авторизован
// @description     - 404: Не найдено
// @description     - 409: Конфликт (логин уже занят)
// @description     - 500: Внутренняя ошибка сервера

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Введите "Bearer" затем пробел и ваш JWT токен. Пример: `Bearer eyJhbGciOiJIUzI1NiIs...`
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"test_task/internal/handler"
	"test_task/internal/pkg/jwt"

	"test_task/internal/repository"
	"test_task/internal/service"

	_ "test_task/docs"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func maskDSN(dsn string) string {
	if !strings.Contains(dsn, "@") {
		return "invalid_dsn_format"
	}

	parts := strings.Split(dsn, "@")
	if len(parts) != 2 {
		return dsn
	}

	credentials := strings.TrimPrefix(parts[0], "postgresql://")
	if strings.Contains(credentials, ":") {
		user := strings.Split(credentials, ":")[0]
		return fmt.Sprintf("postgresql://%s:***@%s", user, parts[1])
	}

	return fmt.Sprintf("postgresql://***@%s", parts[1])
}

func runMigrations(databaseURL string) error {
	logger := slog.Default()

	logger.Info("Проверка миграций БД")
	if !strings.Contains(databaseURL, "sslmode=") {
		databaseURL += "?sslmode=disable"
	}

	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("Все миграции уже применены")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	logger.Info("Миграции успешно применены")
	return nil
}
func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Info("Файл info.env не найден")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	logger.Info("Запуск приложения")

	connStr := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")

	if connStr == "" {
		logger.Error("DATABASE_URL не установлен")
		os.Exit(1)
	}

	if jwtSecret == "" {
		logger.Error("JWT_SECRET не установлен")
		os.Exit(1)
	}
	jwt.Init(jwtSecret)
	logger.Info("Проверка переменных окружения",
		"database_url_masked", maskDSN(connStr),
		"jwt_secret_exists", jwtSecret != "",
		"jwt_secret_length", len(jwtSecret),
		"port", port,
	)

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
	if err := runMigrations(connStr); err != nil {
		logger.Error("Ошибка миграций БД", "error", err)
		os.Exit(1)
	}

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
	// Ping проверяет работу сервера
	// @Summary      🏓 Проверка работы сервера
	// @Description  Простой endpoint для проверки что сервер запущен и работает
	// @Tags         Здоровье
	// @Accept       json
	// @Produce      json
	// @Success      200  {object}  map[string]interface{}  "Сервер работает"
	// @Router       /ping [get]
	// @Example response 200
	// {
	//   "message": "pong",
	//   "status": "ok"
	// }

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "ok",
		})
	})
	// Health проверяет здоровье сервера
	// @Summary      ❤️ Проверка здоровья
	// @Description  Проверяет доступность базы данных и других зависимостей
	// @Tags         Здоровье
	// @Accept       json
	// @Produce      json
	// @Success      200  {object}  map[string]interface{}  "Все системы работают"
	// @Failure      503  {object}  map[string]interface{}  "Проблемы с подключением к БД"
	// @Router       /health [get]
	// @Example response 200
	// {
	//   "status": "healthy",
	//   "message": "Service is running"
	// }
	// @Example response 503
	// {
	//   "status": "error",
	//   "error": "Database unavailable"
	// }

	r.GET("/health", func(c *gin.Context) {

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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("Регистрация маршрутов пользователей")

	api := r.Group("/api")
	{

		users := api.Group("/users")
		{

			users.POST("", userHandler.CreateUser)
			users.POST("/login", userHandler.Login)

			auth := users.Group("/")
			auth.Use(handler.AuthMiddleware())
			{
				auth.GET("", userHandler.GetAllUsers)
				auth.GET("/:id", userHandler.GetUser)
				auth.PUT("/:id", userHandler.UpdateUser)
				auth.DELETE("/:id", userHandler.DeleteUser)

			}
		}
	}
	// В main() после создания роутера добавьте:
	r.GET("/debug/schema", func(c *gin.Context) {
		ctx := context.Background()

		// Получаем информацию о колонках таблицы users
		query := `
        SELECT 
            column_name, 
            data_type, 
            is_nullable,
            column_default
        FROM information_schema.columns 
        WHERE table_name = 'users' 
        ORDER BY ordinal_position
    `

		rows, err := pool.Query(ctx, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		defer rows.Close()

		var columns []map[string]interface{}
		for rows.Next() {
			var columnName, dataType, isNullable string
			var columnDefault *string
			if err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"error":  err.Error(),
				})
				return
			}

			column := map[string]interface{}{
				"name":     columnName,
				"type":     dataType,
				"nullable": isNullable == "YES",
				"default":  columnDefault,
			}
			columns = append(columns, column)
		}

		// Получаем информацию об индексах
		indexQuery := `
        SELECT 
            indexname,
            indexdef
        FROM pg_indexes 
        WHERE tablename = 'users'
    `

		indexRows, err := pool.Query(ctx, indexQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		defer indexRows.Close()

		var indexes []map[string]interface{}
		for indexRows.Next() {
			var indexName, indexDef string
			if err := indexRows.Scan(&indexName, &indexDef); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"error":  err.Error(),
				})
				return
			}
			indexes = append(indexes, map[string]interface{}{
				"name":       indexName,
				"definition": indexDef,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status":            "success",
			"table":             "users",
			"migration_version": 2,
			"columns":           columns,
			"indexes":           indexes,
			"total_columns":     len(columns),
			"total_indexes":     len(indexes),
		})
	})

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
	logger.Info("Запуск сервера", "port", port)
	if err := r.Run(":" + port); err != nil {
		logger.Error("Ошибка запуска сервера", "error", err)
		os.Exit(1)
	}
}
