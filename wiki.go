package main 

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "html/template"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
    filename := p.Title + ".txt"
    // 3rd paramater 0600 indicates UNIX file permissions
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, err := template.ParseFiles(tmpl)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = t.Execute(w, p)
    if err != nil {
        // func Error(w ResponseWriter, error string, code int)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    // Get everything after /view/
    title := r.URL.Path[len("/view/"):]
    p, err := loadPage(title)
    if err != nil {
        // No template redirecting to edit page in order to create
        http.Redirect(w, r, "/edit/" + title, http.StatusFound)
    } else {
       renderTemplate(w, "view_page.html", p)    
   }
}


func editHandler(w http.ResponseWriter, r *http.Request) {
    // Get everything after /edit/
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        fmt.Println("No Page exists creating new one")
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit.html", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    p.save()
    http.Redirect(w, r, "/view/" + title, http.StatusFound)
}


func main() {
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    http.ListenAndServe(":8080", nil)
}

