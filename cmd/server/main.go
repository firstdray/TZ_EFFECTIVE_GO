package main

import (
	"effective/internal/enrichment"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"

	_ "effective/docs"
	"effective/internal/config"
	"effective/internal/handler"
	"effective/internal/repository"
	"effective/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.LoadConfig()

	repo, err := repository.NewPeopleRepository(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	if err := repository.RunMigrations(cfg.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	enricher := enrichment.NewEnricher()
	svc := service.NewPeopleService(repo, enricher)

	h := handler.NewPeopleHandler(svc)

	router := gin.Default()
	h.RegisterRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
