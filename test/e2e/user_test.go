package e2e

import (
	"testing"

	"github.com/AndresCRamos/midas-app-api/handlers"
	"github.com/AndresCRamos/midas-app-api/repository"
	"github.com/AndresCRamos/midas-app-api/services"
	"github.com/AndresCRamos/midas-app-api/utils/firebase"
)

func initUserTest(t *testing.T) *handlers.UserHandler {
	client, err := firebase.GetFireStoreClient()
	if err != nil {
		t.Fatalf("Cant initialize firestore client: %s", err)
	}

	repo := repository.NewUserRepository(client)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	return handler
}

func Test_user_CreateNewUser(t *testing.T) {

}

func Test_user_GetUserByID(t *testing.T) {

}
