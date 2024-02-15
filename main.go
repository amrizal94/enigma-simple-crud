package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var dir = "data/"
var fileName = "books.csv"
var filePath = dir + fileName

type book struct {
	id          int
	title       string
	author      string
	releaseYear string
	pages       int
}

var books []book

func main() {
	var input int
menu:
	for {
		fmt.Print(strings.Repeat("=", 14))
		fmt.Print(" Book Data Management ")
		fmt.Println(strings.Repeat("=", 14))
		fmt.Println("1. View All Books")
		fmt.Println("2. Add New Book")
		fmt.Println("3. Update Book")
		fmt.Println("4. Delete Book")
		fmt.Println("5. Exit")
		fmt.Print("Eneter your choice : ")
		fmt.Scanln(&input)
		switch input {
		case 1:
			err := viewAllbooks()
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			err := addNewBook()
			if err != nil {
				fmt.Println(err)
			}
		case 3:
			err := updateBook()
			if err != nil {
				fmt.Println(err)
			}
		case 4:
			err := deleteBook()
			if err != nil {
				fmt.Println(err)
			}
		case 5:
			break menu
		}
	}
}

func addNewBook() error {
	var input string
	newBook := book{}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Create new book")
	fmt.Print("id book = ")
	scanner.Scan()
	newBook.id, _ = strconv.Atoi(scanner.Text())
	fmt.Print("title book = ")
	scanner.Scan()
	newBook.title = scanner.Text()
	fmt.Print("author book = ")
	scanner.Scan()
	newBook.author = scanner.Text()
	fmt.Print("release year book = ")
	scanner.Scan()
	newBook.releaseYear = scanner.Text()
	fmt.Print("pages book = ")
	scanner.Scan()
	newBook.pages, _ = strconv.Atoi(scanner.Text())

	_, err := findBookById(newBook.id)
	if err == nil {
		return fmt.Errorf("book with id %d already exists", newBook.id)
	}
	fmt.Print("are you sure you want to add this book (y/n)? ")
	fmt.Scanln(&input)
	input = strings.ToLower(input)
	if input[0] != 'y' {
		return fmt.Errorf("adding book is cancelled")
	}
	err = loadDataFromCSV()
	if err != nil {
		return err
	}
	books = append(books, newBook)
	err = saveDataToCSV()
	if err != nil {
		books = books[:len(books)-1]
		return err
	}
	fmt.Println("book added successfully")
	return nil
}

func deleteBook() error {
	var input string
	fmt.Print("id book = ")
	fmt.Scan(&input)
	id, _ := strconv.Atoi(input)
	_, err := findBookById(id)
	if err != nil {
		return err
	}
	fmt.Print("are you sure you want to delete this book (y/n)? ")
	fmt.Scanln(&input)
	input = strings.ToLower(input)
	if input[0] != 'y' {
		return fmt.Errorf("deleting is cancelled")
	}
	temp := books
	for i, book := range books {
		if book.id == id {
			books = append(books[:i], books[i+1:]...)
		}
	}
	err = saveDataToCSV()
	if err != nil {
		books = temp
		return err
	}
	fmt.Println("book deleted successfully")
	return nil
}

func loadDataFromCSV() (err error) {
	books = []book{}
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := strings.Split(scanner.Text(), ",")
		id, _ := strconv.Atoi(record[0])
		pages, _ := strconv.Atoi(record[4])
		book := book{
			id:          id,
			title:       record[1],
			author:      record[2],
			releaseYear: record[3],
			pages:       pages,
		}
		books = append(books, book)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error opening csv file: %w", err)
	}
	return
}

func updateBook() (err error) {
	var input string
	fmt.Print("id book = ")
	fmt.Scan(&input)
	id, _ := strconv.Atoi(input)
	updateBook, err := findBookById(id)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Update book ID: ", updateBook.id)

	fmt.Print("title book = ")
	scanner.Scan()
	if scanner.Text() != "" {
		updateBook.title = scanner.Text()
	}

	fmt.Print("author book = ")
	scanner.Scan()
	if scanner.Text() != "" {
		updateBook.author = scanner.Text()
	}

	fmt.Print("release year book = ")
	scanner.Scan()
	if scanner.Text() != "" {
		updateBook.releaseYear = scanner.Text()
	}

	fmt.Print("pages book = ")
	scanner.Scan()
	if scanner.Text() != "" {
		updateBook.pages, _ = strconv.Atoi(scanner.Text())
	}

	for i, book := range books {
		if book.id == updateBook.id {
			books[i] = updateBook
		}
	}
	fmt.Print("are you sure you want to update this book (y/n)? ")
	fmt.Scanln(&input)
	input = strings.ToLower(input)
	if input[0] != 'y' {
		return fmt.Errorf("updating book is cancelled")
	}
	err = saveDataToCSV()
	if err != nil {
		return
	}
	fmt.Println("updated book successfully")
	return
}

func viewAllbooks() error {
	err := loadDataFromCSV()
	if err != nil {
		return err
	}

	if len(books) == 0 {
		return fmt.Errorf("no book available")
	}

	for i, book := range books {
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("Book -", i+1)
		fmt.Println("Book ID:", book.id)
		fmt.Println("Book Title:", book.title)
		fmt.Println("Book Author:", book.author)
		fmt.Println("Book Release Year:", book.releaseYear)
		fmt.Println("Book Pages:", book.pages)
		fmt.Println(strings.Repeat("=", 50))
	}
	return nil
}

func saveDataToCSV() (err error) {
	err = createFile()
	if err != nil {
		return
	}
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	for _, book := range books {
		row := strconv.Itoa(book.id) + "," + book.title + "," + book.author + "," + book.releaseYear + "," + strconv.Itoa(book.pages) + "\n"
		file.WriteString(row)
	}
	return
}

func findBookById(id int) (book, error) {
	err := loadDataFromCSV()
	if err != nil {
		return book{}, err
	}
	for _, book := range books {
		if book.id == id {
			return book, nil
		}
	}
	return book{}, fmt.Errorf("book id: %d not found", id)
}

func createFile() error {
	// check and make a directory if not exists
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return fmt.Errorf("error making directory: %w", err)
		}
	}

	// check and make a file if not exists
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("error opening csv file: %w", err)
		}
		defer file.Close()
		fmt.Println("File", filePath, "telah dibuat")
	}
	return nil
}
