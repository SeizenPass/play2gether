package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes(staticDir string) http.Handler {

	// create a middleware chain
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	/* Play2GetHer */
	mux.Get("/hub", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showHub))
	mux.Get("/game/add", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.addGameForm))
	mux.Post("/game/add", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.addGame))
	mux.Get("/game/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showGame))
	mux.Get("/game", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showListOfGames))

	mux.Get("/ownership/add/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.addOwnership))
	mux.Get("/ownership/remove/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.removeOwnership))

	mux.Get("/review/add/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.addReviewForm))
	mux.Post("/review/add/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.addReview))
	mux.Get("/review/show/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showReviews))

	//mux by id
	mux.Get("/chat", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showChats))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logoutUser))
	mux.Get("/user/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showUser))
	mux.Post("/user/update", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.updateUser))
	mux.Get("/users", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showUsers))

	fileserver := http.FileServer(http.Dir(staticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileserver))

	return standardMiddleware.Then(mux)
}
