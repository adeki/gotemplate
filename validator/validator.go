package validator

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

type Validator interface {
	Validate(s interface{}) error
	GetMessagesByError(err error) []string
	GetMessage(f, t string) string
}

type myValidator struct {
	validator *validator.Validate
}

func New() Validator {
	v := validator.New()
	v.RegisterTagNameFunc(fieldTagNameToLowerCase)
	v.RegisterValidation("password", passwordFunc)

	return &myValidator{validator: v}
}

// ref. https://godoc.org/gopkg.in/go-playground/validator.v9#Validate.RegisterTagNameFunc
func fieldTagNameToLowerCase(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func (v *myValidator) Validate(s interface{}) error {
	return v.validator.Struct(s)
}

func (v *myValidator) GetMessagesByError(err error) []string {
	errs := err.(validator.ValidationErrors)
	messages := make([]string, 0, len(errs))
	for _, err := range errs {
		msg := getErrorMessage(err.Field(), err.ActualTag())
		messages = append(messages, msg)
	}
	return messages
}

func (v *myValidator) GetMessage(f, t string) string {
	return getErrorMessage(f, t)
}

func getErrorMessage(f, t string) string {
	msg := msgMap["messages"][f+"."+t]
	field := msgMap["fields"][f]
	tag := msgMap["tags"][t]
	if msg != "" {
		return fmt.Sprintf(msg, field)
	} else if tag != "" && field != "" {
		return fmt.Sprintf(tag, field)
	}
	return fmt.Sprintf("%sが不正です。", field)
}

/*
  CUSOM VALIDATE FUMCTIONS
*/

func passwordFunc(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) > 8
}
