package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"example.com/open-monastery/internal/category"
	"example.com/open-monastery/internal/mysql"
)

type jsonCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type jsonCategoryBook struct {
	CategoryID int64  `json:"id"`
	Title      string `json:"categoryTitle"`
	BookID     int64  `json:"bookId"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "import categories:", err)
		os.Exit(1)
	}
}

func readJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func run(ctx context.Context) error {
	var categories []jsonCategory
	if err := readJSON("categories.json", &categories); err != nil {
		return err
	}
	var links []jsonCategoryBook
	if err := readJSON("category_book.json", &links); err != nil {
		return err
	}

	sourceNames := make(map[int64]string)
	for _, c := range categories {
		if c.ID <= 0 || strings.TrimSpace(c.Name) == "" {
			return fmt.Errorf("invalid category: %+v", c)
		}
		if _, exists := sourceNames[c.ID]; exists {
			return fmt.Errorf("duplicate source category ID %d", c.ID)
		}
		sourceNames[c.ID] = c.Name
	}
	for _, link := range links {
		name, exists := sourceNames[link.CategoryID]
		if !exists || name != link.Title || link.BookID <= 0 {
			return fmt.Errorf("invalid category book mapping: %+v", link)
		}
	}

	// Use the same local database configuration as cmd/cli.
	db, err := mysql.Open("root", "Password123!", "localhost", "3306", "om")
	if err != nil {
		return err
	}
	defer db.Close()

	// Validate every referenced book before inserting any categories or links.
	checkedBooks := make(map[int64]bool)
	bookRepo := mysql.NewBookRepo(db)
	for _, link := range links {
		if checkedBooks[link.BookID] {
			continue
		}
		if _, err := bookRepo.FindById(ctx, link.BookID); err != nil {
			return fmt.Errorf("check referenced book %d: %w", link.BookID, err)
		}
		checkedBooks[link.BookID] = true
	}

	categoryService := category.NewService(mysql.NewCategoryRepo(db))
	existingCategories, err := categoryService.GetCategories(ctx)
	if err != nil {
		return err
	}
	idsByName := make(map[string]int64)
	for _, c := range existingCategories {
		if _, exists := idsByName[c.Name]; exists {
			return fmt.Errorf("multiple database categories named %q; resolve duplicates before importing", c.Name)
		}
		idsByName[c.Name] = c.Id
	}
	existingLinks, err := categoryService.GetCategoryBooks(ctx)
	if err != nil {
		return err
	}
	seenLinks := make(map[[2]int64]bool)
	for _, link := range existingLinks {
		seenLinks[[2]int64{link.CategoryId, link.BookId}] = true
	}

	databaseIDs := make(map[int64]int64)
	createdCategories, createdLinks := 0, 0
	for _, c := range categories {
		id, exists := idsByName[c.Name]
		if !exists {
			created, err := categoryService.CreateCategory(ctx, c.Name)
			if err != nil {
				return fmt.Errorf("create category %q: %w", c.Name, err)
			}
			id = created.Id
			idsByName[c.Name] = id
			createdCategories++
		}
		databaseIDs[c.ID] = id
	}
	for _, link := range links {
		categoryID := databaseIDs[link.CategoryID]
		key := [2]int64{categoryID, link.BookID}
		if seenLinks[key] {
			continue
		}
		if err := categoryService.AddBook(ctx, categoryID, link.BookID); err != nil {
			return fmt.Errorf("link category %d to book %d: %w", categoryID, link.BookID, err)
		}
		seenLinks[key] = true
		createdLinks++
	}
	fmt.Printf("Imported %d categories and %d category-book links; existing entries skipped.\n", createdCategories, createdLinks)
	return nil
}
