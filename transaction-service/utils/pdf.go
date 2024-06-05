package utils

import (
    "fmt"
    "time"
    "transaction-service/models"

    "github.com/jung-kurt/gofpdf"
)

func GenerateReceiptPDF(transaction models.Transaction) (string, error) {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)

    pdf.Cell(40, 10, "Transaction Receipt")
    pdf.Ln(12)
    pdf.SetFont("Arial", "", 12)
    pdf.Cell(40, 10, fmt.Sprintf("Transaction ID: %d", transaction.ID))
    pdf.Ln(8)
    pdf.Cell(40, 10, fmt.Sprintf("Date: %s", time.Now().Format(time.RFC1123)))
    pdf.Ln(8)
    pdf.Cell(40, 10, fmt.Sprintf("Customer: %s", transaction.Customer.Name))
    pdf.Ln(8)
    pdf.Cell(40, 10, fmt.Sprintf("Email: %s", transaction.Customer.Email))
    pdf.Ln(12)
    
    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(40, 10, "Items")
    pdf.Ln(10)
    pdf.SetFont("Arial", "", 12)
    total := 0.0
    for _, item := range transaction.CartItems {
        pdf.Cell(40, 10, fmt.Sprintf("%s - %.2f", item.Name, item.Price))
        total += item.Price
        pdf.Ln(6)
    }
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Total: %.2f", total))
    pdf.Ln(12)
    pdf.Cell(40, 10, "Payment Status: Paid")

    pdfPath := fmt.Sprintf("receipt_%d.pdf", transaction.ID)
    err := pdf.OutputFileAndClose(pdfPath)
    if err != nil {
        return "", err
    }

    return pdfPath, nil
}
