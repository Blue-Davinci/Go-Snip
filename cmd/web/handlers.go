package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/blue-davinci/gosnip/internal/models"
	"github.com/go-chi/chi/v5"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Use the template.ParseFiles() function to read the templates into a
	// template set. If there's an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data  and add the snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the render helper.
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError() helper to
	// send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an integer. So we need to manually covert
	// the form data to an integer using strconv.Atoi(), and we send a 400 Bad
	// Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a map to hold all possible errors:
	fieldErrors := make(map[string]string)
	// Check title is not empty or more than 100 characters
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be empty"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}
	// Check that the content value is not empty
	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}
	// Check the ecpires value matches one of the permitted values i.e.
	// a day, week, month or year
	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must equal a day, week or year"
	}
	// If there are any errors, dump them in a plain text HTTP response and
	// return from the handler.
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

	// validation complete, pass it to the DB
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect to the view showing the new snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// Add a new snippetCreate handler, which for now returns a placeholder
// response. We'll update this shortly to show a HTML form.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.notFound(w) // Use the notFound() helper.
		return
	}
	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Use the render helper.
	app.render(w, http.StatusOK, "view.tmpl", data)
}
