package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf := new(bytes.Buffer) // write the template to the buffer, instead straight to the http.ResponseWriter to ensure that there is no error the whold template file

	err := ts.ExecuteTemplate(buf, "base", data)

	if err != nil {
		app.serverError(w, err)
		return
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) decodePostForm(r *http.Request, dst interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)

	decoderError := form.InvalidDecoderError{}

	if err != nil {
		if errors.As(err, &decoderError) {
			panic(err)
		}
		return err
	}

	return nil

}
