package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"snippetbox.kavinddd.net/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func (app *application) newTemplateData() *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// this is required using vanilla because if server get requets /someRandomShit
	// and we have no handler for that, it will use this handler
	// but using httprouter, it matches the path / exactly, so only path "/" can use this handler
	// if r.RequestURI != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData()
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	// vanilla version
	// view?id=10
	// idString := r.URL.Query().Get("id")

	// httprouter version
	// the path is changed, since we can use named parameter, not query parameter
	// view/:id
	params := httprouter.ParamsFromContext(r.Context())
	idString := params.ByName("id")

	id, err := strconv.Atoi(idString)

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {

		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}

		app.serverError(w, err)
		return
	}

	data := app.newTemplateData()
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// vanilla version, we can get rid of this since we uses method-based router from httprouter
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", "POST")
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	content := "O snail\nClimb Mount Fuji, \nBut slowly slowly!\n\n= Kobashi"
	expires := 7

	id, err := app.snippets.Insert("Test", content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("snippet %d is just created", id)
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusTemporaryRedirect)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new page..."))
}
