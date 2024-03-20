package routes

import (
	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/handlers"
	"github.com/AndresCRamos/midas-app-api/repository"
	"github.com/AndresCRamos/midas-app-api/services"
	"github.com/gin-gonic/gin"
)

func addMovementRoutes(client *firestore.Client, r *gin.Engine) {
	repo := repository.NewMovementRepository(client)
	service := services.NewMovementService(repo)
	handler := handlers.NewMovementHandler(service)

	movementGroup := r.Group("/movement")
	{
		movementGroup.GET("/:id", handler.GetMovementByID)
	}

}
