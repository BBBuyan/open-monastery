package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"example.com/open-monastery/internal/book"
	"example.com/open-monastery/internal/mysql"
)

type JsonBook struct {
	Id          string
	Name        string
	Description string
	Price       string
}

func main() {
	fmt.Println("Hello")
	db, err := mysql.Open("root", "Password123!", "localhost", "3306", "om")
	if err != nil {
		panic(err)
	}
	bookRepository := mysql.NewBookRepo(db)
	data, err := os.ReadFile("monastery_books.json")

	if err != nil {
		panic(err)
	}
	var jBooks []JsonBook

	err = json.Unmarshal(data, &jBooks)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	for _, jbook := range jBooks {
		id, err := strconv.ParseInt(jbook.Id, 10, 64)
		if err != nil {
			fmt.Println(jbook.Id, "failed")
			continue
		}
		price, err := strconv.ParseInt(jbook.Price, 10, 32)
		if err != nil {
			fmt.Println(jbook.Id, "failed")
			continue
		}
		var book *book.Book = &book.Book{
			Id:          id,
			Name:        jbook.Name,
			Description: jbook.Description,
			Price:       uint32(price),
		}
		err = bookRepository.Add(ctx, book)
		if err != nil {
			panic(err)
		}
		fmt.Println(book.Id, "done")
	}
}
