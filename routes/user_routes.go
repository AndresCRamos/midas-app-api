package routes

import (
	"github.com/AndresCRamos/midas-app-api/cmd/server"
	"github.com/AndresCRamos/midas-app-api/handlers"
	"github.com/AndresCRamos/midas-app-api/middleware"
	"github.com/AndresCRamos/midas-app-api/repository"
	"github.com/AndresCRamos/midas-app-api/services"
)

func addUserRoutes(server *server.Server) {

	client := server.FirestoreClient
	r := server.Router

	repo := repository.NewUserRepository(client)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	userGroup := r.Group("/user")
	userGroup.Use(middleware.VerifyToken(server.FirebaseAuthClient))
	{
		userGroup.POST("/", handler.CreateNewUser)
		userGroup.GET("/:id", handler.GetUserByID)
	}

}
