package main

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

func BindAndValidate(c echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "failed to bind request")
	}

	found := false

	typeOf := reflect.TypeOf(req)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		tag := field.Tag
		if tag.Get("validate") != "" {
			found = true
			break
		}
	}

	if found {
		if err := c.Validate(req); err != nil {
			return err
		}
	}

	return nil
}
