package main

import (
	"context"
	"log"
	"time"

	"ip-store/backend/internal/api"
	"ip-store/backend/internal/database"
	"ip-store/backend/internal/model"
)

func main() {
	// Initialize the database
	database.InitDB("ip-store.db")

	// Insert initial products
	products := []*model.Product{
		{
			Name:          "Go Programming Ebook",
			Description:   "A comprehensive guide to Go programming.",
			Price:         29.99,
			CoverImageURL: "https://example.com/go-ebook.jpg",
			FileKey:       "go-ebook.pdf",
		},
		{
			Name:          "Next.js Masterclass Video",
			Description:   "Learn Next.js from scratch to advanced.",
			Price:         99.00,
			CoverImageURL: "https://example.com/nextjs-video.jpg",
			FileKey:       "nextjs-masterclass.mp4",
		},
		{
			Name:          "Tailwind CSS UI Kit",
			Description:   "Ready-to-use UI components with Tailwind CSS.",
			Price:         49.50,
			CoverImageURL: "https://example.com/tailwind-ui-kit.jpg",
			FileKey:       "tailwind-ui-kit.zip",
		},
		{
			Name:          "DuckDB Data Analysis Course",
			Description:   "Master data analysis with DuckDB.",
			Price:         75.00,
			CoverImageURL: "https://example.com/duckdb-course.jpg",
			FileKey:       "duckdb-course.zip",
		},
		{
			Name:          "Advanced Go Concurrency Patterns",
			Description:   "Deep dive into Go's concurrency model.",
			Price:         39.99,
			CoverImageURL: "https://example.com/go-concurrency.jpg",
			FileKey:       "go-concurrency.pdf",
		},
		{
			Name:          "React Hooks Best Practices",
			Description:   "Optimize your React applications with Hooks.",			Price:         55.00,
			CoverImageURL: "https://example.com/react-hooks.jpg",
			FileKey:       "react-hooks.pdf",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database.InsertInitialProducts(ctx, products)

	r := api.SetupRouter(database.DBConn)
	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
