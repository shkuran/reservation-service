package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library-microservices/reservation-service/reservation"
)

func RegisterRoutes(server *gin.Engine, reservation reservation.Handler) {
	server.GET("/reservations", reservation.GetReservations)
	server.POST("/reservations", reservation.AddReservation)
	server.POST("/reservations/:id", reservation.CompleteReservation)

}
