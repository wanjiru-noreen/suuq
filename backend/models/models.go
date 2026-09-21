package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Business struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID         int       `json:"id"`
	BusinessID int       `json:"business_id"`
	Name       string    `json:"name"`
	Quantity   int       `json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
}

type StockUpdate struct {
	ID             int       `json:"id"`
	ProductID      int       `json:"product_id"`
	QuantityChange int       `json:"quantity_change"`
	UpdateType     string    `json:"update_type"`
	CreatedAt      time.Time `json:"created_at"`
}

type Debtor struct {
	ID         int       `json:"id"`
	BusinessID int       `json:"business_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	AmountOwed float64   `json:"amount_owed"`
	CreatedAt  time.Time `json:"created_at"`
}

type Creditor struct {
	ID         int       `json:"id"`
	BusinessID int       `json:"business_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	AmountOwed float64   `json:"amount_owed"`
	CreatedAt  time.Time `json:"created_at"`
}

type Transaction struct {
	ID         int       `json:"id"`
	BusinessID int       `json:"business_id"`
	DebtorID   *int      `json:"debtor_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type TransactionItem struct {
	ID            int `json:"id"`
	TransactionID int `json:"transaction_id"`
	ProductID     int `json:"product_id"`
	Quantity      int `json:"quantity"`
}
