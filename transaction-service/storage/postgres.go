package storage

import (
    "context"
    "log"
    "time"
    "transaction-service/models"

    "github.com/jackc/pgx/v4"
)

var db *pgx.Conn

func InitDB(dataSourceName string) {
    config, err := pgx.ParseConfig(dataSourceName)
    if err != nil {
        log.Fatalf("Cannot parse database config: %v", err)
    }

    conn, err := pgx.ConnectConfig(context.Background(), config)
    if err != nil {
        log.Fatalf("Cannot connect to database: %v", err)
    }

    db = conn

    createTableQuery := `
    CREATE TABLE IF NOT EXISTS transactions (
        id SERIAL PRIMARY KEY,
        cart_items JSONB NOT NULL,
        customer JSONB NOT NULL,
        status VARCHAR(50) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL
    );
    `
    _, err = db.Exec(context.Background(), createTableQuery)
    if err != nil {
        log.Fatalf("Cannot create table: %v", err)
    }
}

func CreateTransaction(transaction models.Transaction) error {
    ctx := context.Background()
    _, err := db.Exec(ctx,
        "INSERT INTO transactions (cart_items, customer, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
        transaction.CartItems, transaction.Customer, transaction.Status, transaction.CreatedAt, transaction.UpdatedAt,
    )
    return err
}

func GetTransaction(id int) (models.Transaction, error) {
    ctx := context.Background()
    var transaction models.Transaction
    err := db.QueryRow(ctx,
        "SELECT id, cart_items, customer, status, created_at, updated_at FROM transactions WHERE id=$1", id,
    ).Scan(&transaction.ID, &transaction.CartItems, &transaction.Customer, &transaction.Status, &transaction.CreatedAt, &transaction.UpdatedAt)
    return transaction, err
}

func UpdateTransactionStatus(id int, status string) error {
    ctx := context.Background()
    _, err := db.Exec(ctx,
        "UPDATE transactions SET status=$1, updated_at=$2 WHERE id=$3", status, time.Now(), id,
    )
    return err
}
