package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
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
		tags := typ.Field(i).Tag.Get("validate")
		field := value.Field(i)
		name := strings.ToLower(typ.Field(i).Name)

		for tag := range strings.SplitSeq(tags, ",") {
			switch {
			case tag == "required":
				if field.Kind() == reflect.Pointer && field.IsNil() {
					return fmt.Errorf("%s is required", strings.ToLower(typ.Field(i).Name))
				}
				if field.IsZero() {
					return fmt.Errorf("%s is required", strings.ToLower(typ.Field(i).Name))
				}
			case strings.HasPrefix(tag, "min="):
				min, err := strconv.Atoi(strings.TrimPrefix(tag, "min="))
				if err != nil {
					return fmt.Errorf("invalid min length for %s", name)
				}
				if field.Kind() == reflect.String && field.Len() < min {
					return fmt.Errorf("%s must be at least %d length", name, min)
				}
			case strings.HasPrefix(tag, "max="):
				max, err := strconv.Atoi(strings.TrimPrefix(tag, "max="))
				if err != nil {
					return fmt.Errorf("invalid max length for %s", name)
				}
				if field.Kind() == reflect.String && field.Len() > max {
					return fmt.Errorf("%s must be at most %d length", name, max)
				}
			}
		}

	}

	return nil
}
