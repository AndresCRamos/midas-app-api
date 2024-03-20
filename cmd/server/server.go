package server

import (
	"log"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"github.com/AndresCRamos/midas-app-api/utils/validations"
	"github.com/gin-gonic/gin"
)

type Server struct {
	FirestoreClient    *firestore.Client
	FirebaseAuthClient *auth.Client
	Router             *gin.Engine
}

func NewServer(firestoreClient *firestore.Client, firebaseAuthClient *auth.Client) (*Server, error) {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	err := validations.AddCustomValidations()
	if err != nil {
		return nil, err
	}

	return &Server{
		FirestoreClient:    firestoreClient,
		FirebaseAuthClient: firebaseAuthClient,
		Router:             r,
	}, nil
}

func (s *Server) Run() {
	r := s.Router
	if err := r.Run(); err != nil {
		log.Println(err)
		return
	}
}
