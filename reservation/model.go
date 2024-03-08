package reservation

import "time"

type Reservation struct {
	ID           int64      `json:"id" db:"id"`
	BookId       int64      `json:"book_id" db:"book_id"`
	UserId       int64      `json:"user_id" db:"user_id"`
	CheckoutDate time.Time  `json:"checkout_date" db:"checkout_date"`
	ReturnDate   *time.Time `json:"return_date" db:"return_date"`
}
