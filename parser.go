package main

import (
	"errors"
	"strings"
)

type ResourceMethod = string

const (
	AUTH_LOGIN ResourceMethod = "AUTH_LOGIN"
	T_GET_ALL  ResourceMethod = "T_GET_ALL"
	T_CREATE   ResourceMethod = "T_CREATE"
	T_DELETE   ResourceMethod = "T_DELETE"
	TL_GET_ALL ResourceMethod = "TL_GET_ALL"
	TL_CREATE  ResourceMethod = "TL_CREATE"
	TL_DELETE  ResourceMethod = "TL_DELETE"
)

type ParserOutput struct {
	resourceMethod ResourceMethod
}

func ProcessRequest(request string) ([]string, error) {
	lines := strings.Split(request, "\n")

	if len(lines) != 2 {
		return nil, errors.New("Invalid request form")
	}

	return lines, nil
}
