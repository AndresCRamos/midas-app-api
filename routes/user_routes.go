package routes

import (
	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/handlers"
	"github.com/AndresCRamos/midas-app-api/repository"
	"github.com/AndresCRamos/midas-app-api/services"
	"github.com/gin-gonic/gin"
)

func addUserRoutes(client *firestore.Client, r *gin.Engine) {
	repo := repository.NewUserRepository(client)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	userGroup := r.Group("/user")
	{
		userGroup.POST("/", handler.CreateNewUser)
		userGroup.GET("/:id", handler.GetUserByID)
	}

}
