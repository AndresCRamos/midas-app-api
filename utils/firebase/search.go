package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckDocumentExists(collection *firestore.CollectionRef, id string) (bool, error) {
	docRef, err := collection.Doc(id).Get(context.Background())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, err
	}
	if docRef == nil && !docRef.Exists() {
		return false, nil
	}

	return true, nil

}
