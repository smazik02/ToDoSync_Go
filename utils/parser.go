package utils

import (
	"encoding/json"
	"strings"
)

type ResourceMethod = string

const (
	AUTH_LOGIN ResourceMethod = "AUTH|LOGIN"
	T_GET_ALL  ResourceMethod = "T|GET_ALL"
	T_CREATE   ResourceMethod = "T|CREATE"
	T_DELETE   ResourceMethod = "T|DELETE"
	TL_GET_ALL ResourceMethod = "TL|GET_ALL"
	TL_CREATE  ResourceMethod = "TL|CREATE"
	TL_DELETE  ResourceMethod = "TL|DELETE"
)

type ParserOutput struct {
	ResourceMethod ResourceMethod
	Payload        []byte
}

func ProcessRequest(request string) (*ParserOutput, error) {
	lines := strings.Split(request, "\n")

	if len(lines) != 2 {
		return nil, ParserError{"Invalid request form"}
	}

	method, err := determineMethod(lines[0])
	if err != nil {
		return nil, err
	}

	payload := []byte(lines[1])
	if !json.Valid(payload) {
		return nil, ParserError{"Invalid request body"}
	}

	return &ParserOutput{
		ResourceMethod: method,
		Payload:        payload,
	}, nil
}

func determineMethod(methodString string) (ResourceMethod, error) {
	switch methodString {
	case "AUTH|LOGIN":
		return AUTH_LOGIN, nil
	case "T|GET_ALL":
		return T_GET_ALL, nil
	case "T|CREATE":
		return T_CREATE, nil
	case "T|DELETE":
		return T_DELETE, nil
	case "TL|GET_ALL":
		return TL_GET_ALL, nil
	case "TL|CREATE":
		return TL_CREATE, nil
	case "TL|DELETE":
		return TL_DELETE, nil
	default:
		return "", ParserError{"Method unknown"}
	}
}
