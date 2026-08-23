package middleware

import (
	"net/http"

	"github.com/spikycham/feedme/pkg/logger"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func LoggerMiddleware(h HandlerFunc) http.HandlerFunc {
	log := logger.New()

	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log.Error(err)
		}
	}
}
