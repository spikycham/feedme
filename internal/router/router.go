package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/spikycham/feedme/internal/handler"
	"github.com/spikycham/feedme/internal/middleware"
	"github.com/spikycham/feedme/internal/repository"
)

type Router struct {
	mux *http.ServeMux
}

func New() *Router {
	return &Router{mux: http.NewServeMux()}
}

func (r *Router) handle(pattern string, h middleware.HandlerFunc) {
	r.mux.HandleFunc(pattern, middleware.LoggerMiddleware(h))
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ro.mux.ServeHTTP(w, r)
}

func (r *Router) Register(db *sql.DB) {
	r.handle("GET /", func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "Hey")
		return nil
	})

	r.registerAuth(db)
	r.registerUser(db)
}

// Routes registration helper functions.
func (router *Router) registerAuth(db *sql.DB) {
	r := repository.NewAuthRepository(db)
	h := handler.NewAuthHandler(r)

	router.handle("POST /api/auth/login", h.Login)
	router.handle("POST /api/auth/logout", h.Logout)
}

func (router *Router) registerUser(db *sql.DB) {
	r := repository.NewUserRepository(db)
	h := handler.NewUserHandler(r)

	router.handle("GET /api/user/me", h.Me)
	router.handle("POST /api/user/change/avatar", h.ChangeAvatarURI)
	router.handle("POST /api/user/change/username", h.ChangeUsername)
	router.handle("POST /api/user/change/password", h.ChangePassword)
}
