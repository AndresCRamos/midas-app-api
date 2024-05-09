package integration

import (
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/AndresCRamos/midas-app-api/handlers"
	"github.com/AndresCRamos/midas-app-api/repository"
	"github.com/AndresCRamos/midas-app-api/services"
	test_utils "github.com/AndresCRamos/midas-app-api/utils/test"
)

func initSourceTest(t *testing.T) (*firestore.Client, *handlers.SourceHandler, *handlers.SourceHandler) {
	client := test_utils.InitTestingFireStore(t)
	failClient := test_utils.InitTestingFireStoreFail(t)

	repo := repository.NewSourceRepository(client)
	service := services.NewSourceService(repo)
	handler := handlers.NewSourceHandler(service)

	failRepo := repository.NewSourceRepository(failClient)
	failService := services.NewSourceService(failRepo)
	failHandler := handlers.NewSourceHandler(failService)

	return client, failHandler, handler
}
