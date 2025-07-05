package model

import "time"

type Product struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	CoverImageURL string    `json:"cover_image_url"`
	FileKey       string    `json:"file_key"`
	CreatedAt     time.Time `json:"created_at"`
}
