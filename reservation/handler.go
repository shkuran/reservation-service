package reservation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library-microservices/reservation-service/utils"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) Handler {
	return Handler{repo: repo}
}

func (h Handler) GetReservations(context *gin.Context) {
	reservations, err := h.repo.GetAll()
	if err != nil {
		utils.HandleInternalServerError(context, "Could not fetch reservations!", err)
		return
	}

	context.JSON(http.StatusOK, reservations)
}

func (h Handler) AddReservation(context *gin.Context) {
	var reservation Reservation
	err := context.ShouldBindJSON(&reservation)
	if err != nil {
		utils.HandleBadRequest(context, "Could not parse request data!", err)
		return
	}

	// url := "http://localhost:8081/books/"
	url := "http://" + getBookServiceHost() + "/books/"
	book, err := h.getBookById(url, reservation.BookId)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not fetch book!", err)
		return
	}

	numberOfBookCopies := book.AvailableCopies
	if numberOfBookCopies < 1 {
		utils.HandleBadRequest(context, "The book is not available!", nil)
		return
	}

	userId := context.GetInt64("userId")
	reservation.UserId = userId

	err = h.repo.Save(reservation)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not add reservation!", err)
		return
	}

	err = h.updateAvailableCopies(url, book.ID, book.AvailableCopies-1)
    if err != nil {
        utils.HandleInternalServerError(context, "Failed to update the number of book copies in book service", err)
        return
    }

	utils.HandleStatusCreated(context, "Reservation added!")
}

func (h Handler) CompleteReservation(context *gin.Context) {
	reservationId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		utils.HandleBadRequest(context, "Could not parse reservationId!", err)
		return
	}

	reservation, err := h.repo.GetById(reservationId)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not fetch reservation!", err)
		return
	}

	userId := context.GetInt64("userId")
	if reservation.UserId != userId {
		utils.HandleStatusUnauthorized(context, "Not access to copmlete reservation!", nil)
		return
	}

	if reservation.ReturnDate != nil {
		utils.HandleBadRequest(context, "The reservation is copleted already!", nil)
		return
	}

	err = h.repo.UpdateReturnDate(reservationId)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not copmlete reservation!", err)
		return
	}

	url := "http://" + getBookServiceHost() + "/books/"
	book, err := h.getBookById(url, reservation.BookId)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not fetch book!", err)
		return
	}

	err = h.updateAvailableCopies(url, book.ID, book.AvailableCopies+1)
    if err != nil {
        utils.HandleInternalServerError(context, "Failed to update the number of book copies in book service", err)
        return
    }

	context.JSON(http.StatusOK, gin.H{"message": "Reservation copmleted!"})
}

func (h Handler) ShowAllBooks(context *gin.Context) {
	url := "http://" + getBookServiceHost() + "/books/"
	books, err := h.getBooks(url)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not get books!", err)
		return
	}

	context.JSON(http.StatusOK, books)
}

func getBookServiceHost() string {
	host := os.Getenv("BOOK_SERVICE_HOST")
	if host == "" {
		host = "localhost:8081"
	}
	return host
}

func (h Handler) getBooks(url string) ([]Book, error) {
	response, err := http.Get(url)
	if err != nil {
		return []Book{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []Book{}, fmt.Errorf("failed to get book. Status code: %d", response.StatusCode)
	}

	var books []Book
	err = json.NewDecoder(response.Body).Decode(&books)
	if err != nil {
		return []Book{}, err
	}

	return books, nil
}

func (h Handler) getBookById(bookServiceURL string, bookID int64) (Book, error) {
	url := bookServiceURL + fmt.Sprint(bookID)

	response, err := http.Get(url)
	if err != nil {
		return Book{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Book{}, fmt.Errorf("failed to get book. Status code: %d", response.StatusCode)
	}

	var bookInfo Book
	err = json.NewDecoder(response.Body).Decode(&bookInfo)
	if err != nil {
		return Book{}, err
	}

	return bookInfo, nil
}

func (h Handler) updateAvailableCopies(bookServiceURL string, bookID, availableCopies int64) error {
	updateInfo := struct {
		BookID          int64 `json:"book_id"`
		AvailableCopies int64 `json:"available_copies"`
	}{
		BookID:          bookID,
		AvailableCopies: availableCopies,
	}

	updateInfoJSON, err := json.Marshal(updateInfo)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, bookServiceURL, bytes.NewBuffer(updateInfoJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	
	client := http.Client{}
    response, err := client.Do(req)
    if err != nil {
        return err
    }
    defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update availableCopies. Status code: %d", response.StatusCode)
	}

	return nil
}
