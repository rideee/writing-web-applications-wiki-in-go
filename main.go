package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// Parsing all templates into a single *Template.
// Then we can use the ExecuteTemplate method to render a specific template.
var parsedTemplates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))

// The function regexp.MustCompile will parse and compile the regular expression, and return a regexp.Regexp.
// MustCompile is distinct from Compile in that it will panic if the expression compilation fails, while
// Compile returns an error as a second parameter.
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

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
	err := parsedTemplates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		// Log error to the standard logger (stdout).
		log.Println("Func renderTemplate: ", err)
		// The http.Error function sends a specified HTTP response code (in this case "Internal Server Error")
		// and error message to http.ResponseWriter.
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// make Handler is a wrapper function that takes a function of the above type, and returns a function
// of type http.HandlerFunc (suitable to be passed to the function http.HandleFunc)
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// viewHandler will allow users to view a wiki page.
// It will handle URLs prefixed with "/view/".
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
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
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		// Return empty Page struct when err != nil.
		p = &Page{Title: title}
	}

	// Display an HTML form.
	renderTemplate(w, "edit.html", p)
}

// saveHandler will handle the submission of forms located on the edit pages.
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
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
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/view/Welcome", http.StatusFound)
	})

	port := "8080"
	fmt.Printf("Listening on port: %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
