package utils

import (
	"encoding/json"
	"fmt"
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
	ResourceMethod ResourceMethod
	Payload        []byte
}

type ParserError struct {
	message string
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
	case "AUTH_LOGIN":
		return AUTH_LOGIN, nil
	case "T_GET_ALL":
		return T_GET_ALL, nil
	case "T_CREATE":
		return T_CREATE, nil
	case "T_DELETE":
		return T_DELETE, nil
	case "TL_GET_ALL":
		return TL_GET_ALL, nil
	case "TL_CREATE":
		return TL_CREATE, nil
	case "TL_DELETE":
		return TL_DELETE, nil
	default:
		return "", ParserError{"Method unknown"}
	}
}

func (e ParserError) Error() string {
	body := map[string]string{"body": e.message}
	bodyJSON, _ := json.Marshal(body)
	return fmt.Sprintf("FAIL\n%s\n\n", string(bodyJSON))
}
