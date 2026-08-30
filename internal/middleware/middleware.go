package middleware

import "net/http"

func Register(h http.Handler) http.Handler {
	// The later ones are the outside ones.
	// TODO: remember to restore the auth middleware.
	// return chain(h, AuthMiddleware, CorsMiddleware)
	return chain(h, CorsMiddleware)
}

func chain(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}
