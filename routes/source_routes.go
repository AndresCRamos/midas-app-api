package routes

import (
	"github.com/AndresCRamos/midas-app-api/cmd/server"
	"github.com/AndresCRamos/midas-app-api/handlers"
	"github.com/AndresCRamos/midas-app-api/middleware"
	"github.com/AndresCRamos/midas-app-api/repository"
	"github.com/AndresCRamos/midas-app-api/services"
)

func addSourceRoutes(server *server.Server) {

	client := server.FirestoreClient
	r := server.Router

	repo := repository.NewSourceRepository(client)
	service := services.NewSourceService(repo)
	handler := handlers.NewSourceHandler(service)

	sourceGroup := r.Group("/source")
	sourceGroup.Use(middleware.VerifyToken(server.FirebaseAuthClient))
	{
		sourceGroup.GET("/", handler.GetSourcesByUser)
		sourceGroup.POST("/", handler.CreateNewSource)
		sourceGroup.GET("/:id", handler.GetSourceByID)
		sourceGroup.PUT("/:id", handler.UpdateSource)
		sourceGroup.DELETE("/:id", handler.DeleteSource)
	}

}
