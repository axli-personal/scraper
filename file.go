package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
)

func WriteComments() {
	var (
		file *os.File
		err  error
	)

	err = os.Mkdir("comments", 0666)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		log.Fatal("Fail to create comments directory:", err)
	}

	for id, bookComments := range comments {
		file, err = os.Create(fmt.Sprintf("comments/%v.txt", id))
		if err != nil {
			log.Fatal("Fail to create comments file:", err)
		}
		writeCommentHeader(file)
		for _, comment := range bookComments {
			writeComment(file, comment)
		}
		file.Close()
	}
}

func writeCommentHeader(file *os.File) {
	_, err := file.WriteString("User,Rate,Text\n")
	if err != nil {
		log.Print("Fail to write comment header:", err)
	}
}

func writeComment(file *os.File, comment Comment) {
	commentInfo := fmt.Sprintf("%v,%v,%v\n", comment.User, comment.Rate, comment.Text)
	_, err := file.WriteString(commentInfo)
	if err != nil {
		log.Print("Fail to write comment:", err)
	}
}

func WriteBooks() {
	file, err := os.Create("books.txt")
	if err != nil {
		log.Fatal("Fail to create books file:", err)
	}
	writeBookHeader(file)
	for _, book := range books {
		writeBook(file, book)
	}
	file.Close()
}

func writeBookHeader(file *os.File) {
	_, err := file.WriteString("Name,Author,Rate,Quote,Link\n")
	if err != nil {
		log.Print("Fail to write book header:", err)
	}
}

func writeBook(file *os.File, book Book) {
	bookInfo := fmt.Sprintf("%v,%v,%v,%v,%v\n", book.Name, book.Author, book.Rate, book.Quote, book.Link)
	_, err := file.WriteString(bookInfo)
	if err != nil {
		log.Print("Fail to write book:", err)
	}
}
