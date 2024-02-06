package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippet.devlake.xyz/internal/models"
	"snippet.devlake.xyz/internal/validator"

	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	validator.Validator `       form:"-"`
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
}

type userSignupForm struct {
	validator.Validator `       form:"-"`
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
}

// Base Handlers

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippedView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

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

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

// Snippet Creation Handlers

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 1095,
	}

	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Parse and decode form values
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validate form values
	form.CheckField(validator.NotBlank(form.Title), "title", "Title cannot be blank")
	form.CheckField(
		validator.MaxChars(form.Title, 100),
		"title",
		"Title cannot be longer than 100 characters",
	)

	form.CheckField(validator.NotBlank(form.Content), "content", "Content cannot be blank")
	form.CheckField(
		validator.PermittedInt(form.Expires, 1, 7, 365, 1095),
		"expires",
		"Expires must be equal to 1, 7, 365, or 1095 days",
	)

	// if there are any validation errors re-render create snippet template
	// with user values and validation errors
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// add flash message data to requesting users session
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	// return response
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// User Handlers

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Parse form in the request
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate form
	form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be blank")

	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	form.CheckField(
		validator.MatchesRegex(form.Email, validator.EmailRX),
		"email",
		"Email address is invalid",
	)

	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.CheckField(
		validator.MinChars(form.Password, 8),
		"password",
		"Password must be at least 8 characters long",
	)

	// If form validation failed return errorLog
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	// handle create new user erorrs
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please login.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
