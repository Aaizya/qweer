package models

import (
    "time"
)

type CartItem struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

type Customer struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type Transaction struct {
    ID        int       `json:"id"`
    CartItems []CartItem `json:"cartItems"`
    Customer  Customer  `json:"customer"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

type PaymentForm struct {
    CardNumber     string `json:"cardNumber"`
    ExpirationDate string `json:"expirationDate"`
    CVV            string `json:"cvv"`
    Name           string `json:"name"`
    Address        string `json:"address"`
}
