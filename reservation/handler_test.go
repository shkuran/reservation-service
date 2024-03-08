package reservation

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupTestEnv(reservationsInDB []Reservation) TestEnv {
	resRepo := NewMockReservationRepo(reservationsInDB)
	resHandler := NewHandler(resRepo)

	return TestEnv{
		ReservationRepo:    resRepo,
		ReservationHandler: resHandler,
	}
}

type TestEnv struct {
	ReservationRepo    *MockReservationRepo
	ReservationHandler Handler
}

func TestGetReservations(t *testing.T) {

	testCases := []struct {
		testName             string
		// booksInDB            []book.Book
		reservationsInDB     []Reservation
		expectedCode         int
		expectedReservations []Reservation
		expectedErrorMsg     string
	}{
		// Case 1: GetReservation returns []Reservation
		{
			testName:             "Return reservations",
			// booksInDB:            []book.Book{},
			reservationsInDB:     []Reservation{{ID: 1, BookId: 1, UserId: 1}, {ID: 2, BookId: 2, UserId: 2}},
			expectedCode:         http.StatusOK,
			expectedReservations: []Reservation{{ID: 1, BookId: 1, UserId: 1}, {ID: 2, BookId: 2, UserId: 2}},
			expectedErrorMsg:     "",
		},
		// Case 2: GetReservation returns an error
		{
			testName:             "Return an error",
			// booksInDB:            []book.Book{},
			reservationsInDB:     []Reservation{},
			expectedCode:         http.StatusInternalServerError,
			expectedReservations: nil,
			expectedErrorMsg:     "Could not fetch reservations!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// Set up the test environment
			env := setupTestEnv(tc.reservationsInDB)

			router := gin.Default()

			router.GET("/reservations", env.ReservationHandler.GetReservations)

			// Perform a test request
			req, err := http.NewRequest("GET", "/reservations", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder to record the response
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("Expected status %d; got %d", tc.expectedCode, w.Code)
			}

			if tc.expectedErrorMsg != "" {
				// Check if the response contains the expected error message
				var response map[string]string
				err = json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				if response["message"] != tc.expectedErrorMsg {
					t.Errorf("Expected error message '%s'; got '%s'", tc.expectedErrorMsg, response["message"])
				}
			} else {
				// Check if the response matches the expected reservations
				var responseReservations []Reservation
				err = json.Unmarshal(w.Body.Bytes(), &responseReservations)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(tc.expectedReservations, responseReservations) {
					t.Errorf("Expected %+v; got %+v", tc.expectedReservations, responseReservations)
				}
			}
		})
	}

}

func TestAddReservation(t *testing.T) {
	testCases := []struct {
		testName         string
		// booksInDB        []book.Book
		reservationsInDB []Reservation
		requestBody      string
		expectedCode     int
		expectedErrorMsg string
	}{
		// Case 1: AddReservation adds new reservation and update AvailableCopies
		{
			testName:         "Successfully added reservation",
			// booksInDB:        []book.Book{{ID: 1, Title: "Book_1", AvailableCopies: 1}},
			reservationsInDB: []Reservation{},
			requestBody:      `{"book_id": 1}`,
			expectedCode:     http.StatusCreated,
			expectedErrorMsg: "",
		},
		// Case 2: AddReservation returns a bad request
		{
			testName:         "Bad request",
			// booksInDB:        []book.Book{},
			reservationsInDB: []Reservation{},
			requestBody:      `{"book_id": 1a}`,
			expectedCode:     http.StatusBadRequest,
			expectedErrorMsg: "Could not parse request data!",
		},
		// Case 3: AddReservation could not fetch book! Returns InternalServerError
		{
			testName:         "No books",
			// booksInDB:        []book.Book{},
			reservationsInDB: []Reservation{},
			requestBody:      `{"book_id": 18}`,
			expectedCode:     http.StatusInternalServerError,
			expectedErrorMsg: "Could not fetch book!",
		},
		// Case 4: AddReservation returns a bad request. The book is not available!
		{
			testName:         "AvailableCopies is 0",
			// booksInDB:        []book.Book{{ID: 1, Title: "Book_1", AvailableCopies: 0}},
			reservationsInDB: []Reservation{},
			requestBody:      `{"book_id": 1}`,
			expectedCode:     http.StatusBadRequest,
			expectedErrorMsg: "The book is not available!",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// Set up the test environment
			env := setupTestEnv(tc.reservationsInDB)

			// HTTP request
			req := httptest.NewRequest(http.MethodPost, "/reservations", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Gin context
			gin.SetMode(gin.TestMode)
			context, _ := gin.CreateTestContext(w)
			context.Request = req
			context.Set("userId", int64(1))

			// Perform the request
			env.ReservationHandler.AddReservation(context)

			if w.Code != tc.expectedCode {
				t.Errorf("Expected status %d; got %d", tc.expectedCode, w.Code)
			}

			if tc.expectedErrorMsg == "" {
				// Check if AvailableCopies of book with id:1 was updated(was: 1, should be: 0)
				// reservedBook, err := env.BookRepo.GetById(1)
				// if err != nil {
				// 	t.Errorf("Could not fetch book! error: %d", err)
				// }
				// if reservedBook.AvailableCopies != 0 {
				// 	t.Errorf("Expected AvailableCopies %d; got %d", 0, reservedBook.AvailableCopies)
				// }

				// Check if reservation was added
				expRes := Reservation{ID: 1, BookId: 1, UserId: 1}
				gotedRes, err := env.ReservationRepo.GetById(1)
				if err != nil {
					t.Errorf("Could not fetch book! error: %d", err)
				}
				if !reflect.DeepEqual(gotedRes, expRes) {
					t.Errorf("Expected new rservation id:%d, book_id:%d, user_id:%d; got id:%d, book_id:%d, user_id:%d",
						expRes.ID, expRes.BookId, expRes.UserId, gotedRes.ID, gotedRes.BookId, gotedRes.UserId)
				}
			} else {
				// Check if the response contains the expected error message
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				if response["message"] != tc.expectedErrorMsg {
					t.Errorf("Expected error message '%s'; got '%s'", tc.expectedErrorMsg, response["message"])
				}

			}
		})
	}

}

func TestCompleteReservation(t *testing.T) {
	testCases := []struct {
		testName         string
		// booksInDB        []book.Book
		reservationsInDB []Reservation
		reservationId    string
		expectedCode     int
		expectedErrorMsg string
	}{
		// Case 1: CopleteReservation add return date for reservation and update AvailableCopies for book
		{
			testName:         "Successfully completed reservation",
			// booksInDB:        []book.Book{{ID: 1, Title: "Book_1", AvailableCopies: 1}},
			reservationsInDB: []Reservation{{ID: 1, BookId: 1, UserId: 1, ReturnDate: nil}},
			reservationId:    "1",
			expectedCode:     http.StatusOK,
			expectedErrorMsg: "",
		},
		// Case 2: CopleteReservation returns a bad request
		{
			testName:         "Bad request",
			// booksInDB:        []book.Book{},
			reservationsInDB: []Reservation{},
			reservationId:    "a",
			expectedCode:     http.StatusBadRequest,
			expectedErrorMsg: "Could not parse reservationId!",
		},
		// Case 3: CopleteReservation could not fetch reservation! Returns InternalServerError
		{
			testName:         "No resrvation with this id",
			// booksInDB:        []book.Book{{ID: 1, Title: "Book_1", AvailableCopies: 1}},
			reservationsInDB: []Reservation{{ID: 1, BookId: 1, UserId: 1, ReturnDate: nil}},
			reservationId:    "2",
			expectedCode:     http.StatusInternalServerError,
			expectedErrorMsg: "Could not fetch reservation!",
		},
		// Case 4: CopleteReservation returns a StatusUnauthorized. User1 cannot complete reservation of user2
		{
			testName:         "No access to reservation",
			// booksInDB:        []book.Book{{ID: 1, Title: "Book_1", AvailableCopies: 1}},
			reservationsInDB: []Reservation{{ID: 1, BookId: 1, UserId: 2, ReturnDate: nil}},
			reservationId:    "1",
			expectedCode:     http.StatusUnauthorized,
			expectedErrorMsg: "Not access to copmlete reservation!",
		},
		// Case 5: Cannot complete reservation if returnDate is not nil
		{
			testName:         "Rreservation is completed already",
			// booksInDB:        []book.Book{{ID: 1, Title: "Book_1", AvailableCopies: 1}},
			reservationsInDB: []Reservation{{ID: 1, BookId: 1, UserId: 1, ReturnDate: &time.Time{}}},
			reservationId:    "1",
			expectedCode:     http.StatusBadRequest,
			expectedErrorMsg: "The reservation is copleted already!",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// Set up the test environment
			env := setupTestEnv(tc.reservationsInDB)

			// HTTP request
			req := httptest.NewRequest(http.MethodPost, "/reservations", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Gin context
			gin.SetMode(gin.TestMode)
			context, _ := gin.CreateTestContext(w)
			context.Request = req
			context.Set("userId", int64(1))
			context.AddParam("id", tc.reservationId)

			// Perform the request
			env.ReservationHandler.CompleteReservation(context)

			if w.Code != tc.expectedCode {
				t.Errorf("Expected status %d; got %d", tc.expectedCode, w.Code)
			}

			if tc.expectedErrorMsg == "" {
				// Check if AvailableCopies of book with id:1 was updated(was: 1, should be: 2)
				// reservedBook, err := env.BookRepo.GetById(1)
				// if err != nil {
				// 	t.Errorf("Could not fetch book! error: %v", err)
				// }
				// if reservedBook.AvailableCopies != 2 {
				// 	t.Errorf("Expected AvailableCopies %d; got %d", 2, reservedBook.AvailableCopies)
				// }
				
				// Check if reservation was completed
				gotRes, err := env.ReservationRepo.GetById(1)
				if err != nil {
					t.Errorf("Could not fetch reservation! error: %v", err)
				}
				if gotRes.ReturnDate == nil {
					t.Errorf("Rservation ReturnDate was not updated: %+v;", gotRes)
				}
			} else {
				// Check if the response contains the expected error message
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				if response["message"] != tc.expectedErrorMsg {
					t.Errorf("Expected error message '%s'; got '%s'", tc.expectedErrorMsg, response["message"])
				}

			}
		})
	}

}
