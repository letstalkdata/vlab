package models

import (
	"context"
	"fmt"
	"time"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"strings"
)

type Item struct {
	ID 				uuid.UUID 	`json:"id"`
	CreatedAt 		time.Time 	`json:"_"`
	UpdatedAt 		time.Time 	`json:"_"`
	Title 			string 		`json:"title"`
	Notes 			string 		`json:"notes"`
	SellerID 		uuid.UUID 	`json:"seller"`
	PriceInCents 	int64 		`json:"price_in_cents"`
}

func (i *Item) Create(conn *pgx.Conn, userID string) error {
	i.Title = strings.Trim(i.Title, " ")
	if len(i.Title) < 1 {
		return fmt.Errorf("Title must not be empty")
	}
	if i.PriceInCents < 0 {
		i.PriceInCents = 0
	}
	now := time.Now()
	row := conn.QueryRow(context.Background(), "INSERT INTO item (title, notes, seller_id, price_in_cents, created_at, updated_at) VALUES ($1, $2,$3, $4, $5, $6) RETURNING id, seller_id", i.Title, i.Notes, userID, i.PriceInCents, now, now)
	err := row.Scan(&i.ID, &i.SellerID)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("There was an error creating the item")
	}

	return nil
}