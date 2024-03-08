package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleBadRequest(context *gin.Context, message string, err error) {
	context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message})
	if err != nil {
		log.Println(err)
	}
}

func HandleInternalServerError(context *gin.Context, message string, err error) {
	context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": message})
	if err != nil {
		log.Println(err)
	}
}

func HandleStatusUnauthorized(context *gin.Context, message string, err error) {
	context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": message})
	if err != nil {
		log.Println(err)
	}
}

func HandleStatusCreated(context *gin.Context, message string) {
	context.JSON(http.StatusCreated, gin.H{"message": message})
}
