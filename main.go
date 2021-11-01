package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// Page struct describes how page data will be stored in memory.
type Page struct {
	Title string
	Body  []byte
}

// This method will save the Page's Body to a text file with .page extension.
// For simplicity, we will use the Title as the file name.
// All pages will be stored in 'pages' directory.
func (p *Page) save() error {
	filename := p.Title + ".page"
	filepath := "pages/" + filename
	return ioutil.WriteFile(filepath, p.Body, 0600)
}

// loadPage constructs the file name from the title parameter, reads the file's contents into a new variable body,
// and returns a pointer to a Page literal constructed with the proper title and body values.
func loadPage(title string) (*Page, error) {
	filename := title + ".page"
	filepath := "pages/" + filename
	body, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// renderTemplate parses and executes the template, then writing generated HTML to http.ResponseWriter.
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	// The function template.ParseFiles will read the contents of tmpl and return a *template.Template.
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		// Log error to the standard logger (stdout).
		log.Println("Func renderTemplate: ", err)
		// The http.Error function sends a specified HTTP response code (in this case "Internal Server Error")
		// and error message to http.ResponseWriter.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// The method t.Execute executes the template, writing the generated HTML to the http.ResponseWriter.
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// viewHandler will allow users to view a wiki page.
// It will handle URLs prefixed with "/view/".
func viewHandler(w http.ResponseWriter, r *http.Request) {
	// [len("/view/"):]	means that everything after "/view/" will be taken.
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		// err != nil means that page file not exsist. Print the log and redirect to editHandler route.
		log.Printf("Func viewHandler: %v; Redirect to /edit/%s\n", err, title)
		// The http.Redirect function adds an HTTP status code of http.StatusFound (302)
		// and a Location header to the HTTP response.
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view.html", p)
}

// editHandler loads the page (or, if it doesn't exist, create an empty Page struct), and displays an HTML form.
func editHandler(w http.ResponseWriter, r *http.Request) {
	// [len("/edit/"):]	means that everything after "/edit/" will be taken.
	title := r.URL.Path[len("/edit/"):]

	p, err := loadPage(title)
	if err != nil {
		// Return empty Page struct when err != nil.
		p = &Page{Title: title}
	}

	// Display an HTML form.
	renderTemplate(w, "edit.html", p)
}

// saveHandler will handle the submission of forms located on the edit pages.
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]

	// Get body value from the POST request (<textarea name="body">, edit.html template).
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}

	err := p.save()
	if err != nil {
		log.Println("Func saveHandler: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	// Simple routing.
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/view/Welcome", http.StatusFound)
	})

	port := "8080"
	fmt.Printf("Listening on port: %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
