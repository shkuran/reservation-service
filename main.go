package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library-microservices/reservation-service/config"
	"github.com/shkuran/go-library-microservices/reservation-service/db"
	"github.com/shkuran/go-library-microservices/reservation-service/reservation"
	"github.com/shkuran/go-library-microservices/reservation-service/routes"
)

func main() {
	conf := config.LoadConfig()
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = conf.Database.Host
	}
	log.Printf("db_host: %s", host)
	port := conf.Database.Port
	user := conf.Database.User
	pass := conf.Database.Password
	dbName := conf.Database.DbName
	sslMode := conf.Database.SslMode
	driverName := conf.Database.DriverName
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbName, sslMode)

	varDb, err := db.InitDB(driverName, connStr)
	if err != nil {
		log.Fatal(err)
		return
	}

	//db.CreateTables(varDb)

	server := gin.Default()

	reservationRepo := reservation.NewRepo(varDb)
	reservationHandler := reservation.NewHandler(reservationRepo)

	routes.RegisterRoutes(server, reservationHandler)

	err = server.Run(":" + conf.Server.Port)
	if err != nil {
		log.Fatal(err)
		return
	}
}
