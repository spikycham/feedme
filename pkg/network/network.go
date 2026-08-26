package network

import (
	"encoding/json"
	"net/http"

	"github.com/spikycham/feedme/internal/constant"
	"github.com/spikycham/feedme/pkg/validator"
)

type Response struct {
	Data any `json:"data"`
}

func Read[Data any](r *http.Request, data *Data) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return constant.InvalidJSON
	}

	if err := validator.Validate(data); err != nil {
		return err
	}

	return nil
}

func Write[Data any](w http.ResponseWriter, data *Data) error {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{
		Data: data,
	}); err != nil {
		return err
	}
	return nil
}

func WriteEmpty(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func Error(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
