package routes

import "github.com/AndresCRamos/midas-app-api/cmd/server"

func AddRoutes(server *server.Server) {
	addMovementRoutes(server.FirestoreClient, server.Router)
	addUserRoutes(server.FirestoreClient, server.Router)
}
