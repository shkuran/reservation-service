package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library-microservices/reservation-service/middleware"
	"github.com/shkuran/go-library-microservices/reservation-service/reservation"
)

func RegisterRoutes(server *gin.Engine, reservation reservation.Handler) {
	server.GET("/reservations", reservation.GetReservations)
	server.GET("/reservations/books", reservation.ShowAllBooks)
	server.POST("/reservations", middleware.Authenticate, reservation.AddReservation)
	server.POST("/reservations/:id", middleware.Authenticate, reservation.CompleteReservation)

}
