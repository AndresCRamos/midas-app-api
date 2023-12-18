package main

import (
	"fmt"
	"log"

	"github.com/AndresCRamos/midas-app-api/cmd/server"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	firebase_utils "github.com/AndresCRamos/midas-app-api/utils/firebase"
)

func main() {
	firestoreClient, err := firebase_utils.GetFireStoreClient()
	if err != nil {
		final_err := fmt.Errorf(error_utils.INITIALIZE_APP_ERROR, err)
		log.Println(final_err)
		return
	}
	authClient, err := firebase_utils.GetFirebaseAuthClient()
	if err != nil {
		final_err := fmt.Errorf(error_utils.INITIALIZE_APP_ERROR, err)
		log.Println(final_err)
		return
	}

	server := server.NewServer(firestoreClient, authClient)

	server.Run()
}
