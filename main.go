package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

// Page struct describes how page data will be stored in memory.
type Page struct {
	Title string
	Body  []byte
}

// This method will save the Page's Body to a text file. For simplicity, we will use the Title as the file name.
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func main() {
	// Saving a new Page.
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()

	// Reading the Page from file.
	p2, err := loadPage("TestPage")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(p2.Body))
}

// loadPage constructs the file name from the title parameter, reads the file's contents into a new variable body,
// and returns a pointer to a Page literal constructed with the proper title and body values.
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
