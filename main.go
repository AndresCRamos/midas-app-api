package main

import (
	"log"

	"github.com/AndresCRamos/midas-app-api/cmd/server"
	"github.com/AndresCRamos/midas-app-api/routes"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	firebase_utils "github.com/AndresCRamos/midas-app-api/utils/firebase"
)

func main() {
	firestoreClient, err := firebase_utils.GetFireStoreClient()
	if err != nil {
		final_err := error_utils.InitializeAppError{}
		final_err.Wrap(err)
		log.Println(final_err)
		return
	}
	authClient, err := firebase_utils.GetFirebaseAuthClient()
	if err != nil {
		final_err := error_utils.InitializeAppError{}
		final_err.Wrap(err)
		log.Println(final_err)
		return
	}

	server := server.NewServer(firestoreClient, authClient)
	routes.AddRoutes(server)

	server.Run()
}
