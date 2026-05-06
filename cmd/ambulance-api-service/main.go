package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xtruhlar/dt26-ambulance-webapi/api"
	"github.com/xtruhlar/dt26-ambulance-webapi/internal/ambulance_ufe"
	"github.com/xtruhlar/dt26-ambulance-webapi/internal/db_service"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("AMBULANCE_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("AMBULANCE_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	dbService := db_service.NewMongoService[ambulance_ufe.Ambulance](db_service.MongoServiceConfig{})
	defer dbService.Disconnect(context.Background())
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service", dbService)
		ctx.Next()
	})

	handleFunctions := &ambulance_ufe.ApiHandleFunctions{
		AmbulanceRemoteConsultationAPI: ambulance_ufe.NewAmbulanceRemoteConsultationApi(),
	}
	ambulance_ufe.NewRouterWithGinEngine(engine, *handleFunctions)

	ambulancesAPI := ambulance_ufe.NewAmbulancesApi()
	engine.POST("/api/ambulance", ambulancesAPI.CreateAmbulance)
	engine.DELETE("/api/ambulance/:ambulanceId", ambulancesAPI.DeleteAmbulance)

	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
