package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SeizenPass/play2gether/pkg/forms"
	"github.com/SeizenPass/play2gether/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.tmpl", &templateData{
	})
}

func (app *application) showHub(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "hub.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) showGame(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.games.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	ows, err := app.gamesOwnerships.GetByGameID(id)
	users := []*models.User{}
	for _, v := range ows {
		u, err := app.users.Get(v.UserID)
		if err == nil {
			users = append(users, u)
		}
	}
	userID := app.session.GetInt(r, "userID")
	ow, err := app.gamesOwnerships.GetByUserIDAndGameID(userID, id)
	app.render(w, r, "game.page.tmpl", &templateData{
		Game: s,
		Users: users,
		Ownership: ow,
	})
}

func (app *application) addOwnership(w http.ResponseWriter, r *http.Request)  {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	userID := app.session.GetInt(r, "userID")
	_, err = app.gamesOwnerships.Insert(id, userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// add a flash message to the session to indicate to the user success
	app.session.Put(r, "flash", "Game was added to your library!")

	http.Redirect(w, r, fmt.Sprintf("/game/%d", id), http.StatusSeeOther)
}

func (app *application) removeOwnership(w http.ResponseWriter, r *http.Request)  {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	ow, err := app.gamesOwnerships.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	err = app.gamesOwnerships.Remove(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// add a flash message to the session to indicate to the user success
	app.session.Put(r, "flash", "Game was removed from your library!")

	http.Redirect(w, r, fmt.Sprintf("/game/%d", ow.GameID), http.StatusSeeOther)
}

func (app *application) showListOfGames(w http.ResponseWriter, r *http.Request) {
	s, err := app.games.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "games.page.tmpl", &templateData{
		Games: s,
	})
}

func (app *application) showUsers(w http.ResponseWriter, r *http.Request) {
	s, err := app.users.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "users.page.tmpl", &templateData{
		Users: s,
	})
}

func (app *application) showChats(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "userID")
	s, err := app.chatMessages.GetAllByUserID(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	chats := []*models.Chat{}
	chatMap :=	make(map[int][]*models.ChatMessage)
	for _, v := range s {
		var resID int
		if v.SenderID == userID {
			resID = v.ReceiverID
		} else {
			resID = v.SenderID
		}
		if chatMap[resID] == nil {
			chatMap[resID] = []*models.ChatMessage{}
		}
		chatMap[resID] = append(chatMap[resID], v)
	}

	for k, v := range chatMap {
		chat := &models.Chat{}
		u, err := app.users.Get(k)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}
		chat.Companion = u
		chat.Messages = v
		unread := 0
		for _, ur := range chat.Messages {
			if !ur.IsRead {
				unread++
			}
		}
		chat.Unread = unread
		chats = append(chats, chat)
	}

	app.render(w, r, "chats.page.tmpl", &templateData{
		Chats: chats,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) addReviewForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	user, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "add.review.page.tmpl", &templateData{
		User: user,
		Form: forms.New(nil),
	})
}

func (app *application) addReview(w http.ResponseWriter, r *http.Request)  {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	// parse POST data into PostForm map
	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("text")
	form.MaxLength("text", 1500)

	if !form.Valid() {
		user, err := app.users.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}
		app.render(w, r, "add.game.page.tmpl", &templateData{
			Form: form,
			User: user,
		})
		return
	}
	userID := app.session.GetInt(r,"userID")
	_, err = app.reviews.Insert(form.Get("text"), userID, id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// add a flash message to the session to indicate to the user success
	app.session.Put(r, "flash", "Review successfully added!")

	http.Redirect(w, r, fmt.Sprintf("/user/%d", id), http.StatusSeeOther)
}

func (app *application) addGameForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "add.game.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) addGame(w http.ResponseWriter, r *http.Request)  {
	// parse POST data into PostForm map
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title")
	form.MaxLength("title", 255)
	form.MaxLength("description", 2000)
	form.MaxLength("image_link", 500)

	if !form.Valid() {
		app.render(w, r, "add.game.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.games.Insert(form.Get("title"), form.Get("image_link"), form.Get("description"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// add a flash message to the session to indicate to the user success
	app.session.Put(r, "flash", "Game successfully added!")

	http.Redirect(w, r, fmt.Sprintf("/game/%d", id), http.StatusSeeOther)
}

func (app *application) showUser(w http.ResponseWriter, r *http.Request)  {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	ows, err := app.gamesOwnerships.GetByUserID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	games := []*models.Game{}
	for _, v := range ows{
		g, err := app.games.Get(v.GameID)
		if err == nil {
			games = append(games, g)
		}
	}
	app.render(w, r, "profile.page.tmpl", &templateData{
		User: s,
		Games: games,
	})
}

func (app *application) showReviews(w http.ResponseWriter, r *http.Request)  {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	reviews, err := app.reviews.GetByReviewedID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "reviews.page.tmpl", &templateData{
		User: s,
		Reviews: reviews,
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// parse POST data into PostForm map
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// add a flash message to the session to indicate to the user success
	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password", "image_link")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"),
		form.Get("password"), form.Get("image_link"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("bio")
	form.MatchesPattern("email", forms.EmailRX)
	form.MaxLength("bio", 1500)
	// TODO: revert if form is not valid

	id := app.session.GetInt(r, "userID")
	err = app.users.Update(id, form.Get("bio"))
	app.session.Put(r, "flash", "Your bio was updated.")

	http.Redirect(w, r, fmt.Sprintf("/user/%d", id), http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{
			Form: form,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userID", id)

	http.Redirect(w, r, fmt.Sprintf("/user/%d", id), http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	app.session.Put(r, "flash", "You've been logged out successfully")
	http.Redirect(w, r, "/", 303)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) profilePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "profile.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

