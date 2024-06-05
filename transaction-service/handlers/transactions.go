package handlers

import (
    "encoding/json"
    "fmt"
    "html/template"
    "net/http"
    "os"
    "strconv"
    "time"
    "transaction-service/models"
    "transaction-service/storage"
    "transaction-service/utils"

    "github.com/gorilla/mux"
    "gopkg.in/gomail.v2"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
    var transaction models.Transaction
    if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    transaction.Status = "pending"
    transaction.CreatedAt = time.Now()
    transaction.UpdatedAt = time.Now()

    if err := storage.CreateTransaction(transaction); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(transaction)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    transaction, err := storage.GetTransaction(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(transaction)
}

func PayTransaction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var paymentForm models.PaymentForm
    if err := json.NewDecoder(r.Body).Decode(&paymentForm); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Assume payment is always successful
    if err := storage.UpdateTransactionStatus(id, "paid"); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    transaction, err := storage.GetTransaction(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    pdfPath, err := utils.GenerateReceiptPDF(transaction)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Send receipt via email
    if err := sendEmail(transaction.Customer.Email, "Transaction Receipt", "Here is your receipt.", pdfPath); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Payment successful and receipt sent",
        "receipt": pdfPath,
    })
}

func sendEmail(to, subject, body, attachmentPath string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", os.Getenv("FROM_EMAIL"))
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)
    m.Attach(attachmentPath)

    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
    if err != nil {
        return fmt.Errorf("invalid SMTP_PORT: %v", err)
    }
    smtpUser := os.Getenv("SMTP_USER")
    smtpPassword := os.Getenv("SMTP_PASSWORD")

    d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)

    return d.DialAndSend(m)
}

// HTML Handlers
func RenderCreateTransactionPage(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("create_transaction.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

func RenderPayTransactionPage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    transaction, err := storage.GetTransaction(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    tmpl, err := template.ParseFiles("pay_transaction.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, transaction)
}
