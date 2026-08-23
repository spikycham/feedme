package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Validate[Schema any](v Schema) error {
	value := reflect.ValueOf(v)

	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return errors.New("expected struct")
	}

	typ := value.Type()

	for i := range value.NumField() {
		tag := typ.Field(i).Tag.Get("validate")

		// We only need the required tag now...
		switch tag {
		case "required":
			if value.Field(i).IsZero() {
				return fmt.Errorf("%s is required", strings.ToLower(typ.Field(i).Name))
			}
		}
	}

	return nil
}
