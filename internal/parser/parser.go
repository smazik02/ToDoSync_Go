package parser

import (
	"encoding/json"
	"strings"
)

type ResourceMethod = string

const (
	AuthLogin      ResourceMethod = "AUTH|LOGIN"
	TaskGetAll     ResourceMethod = "T|GET_ALL"
	TaskCreate     ResourceMethod = "T|CREATE"
	TaskDelete     ResourceMethod = "T|DELETE"
	TaskListGetAll ResourceMethod = "TL|GET_ALL"
	TaskListCreate ResourceMethod = "TL|CREATE"
	TaskListDelete ResourceMethod = "TL|DELETE"
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

	if err := determineMethod(lines[0]); err != nil {
		return nil, err
	}

	method := ResourceMethod(lines[0])
	payload := []byte(lines[1])
	if !json.Valid(payload) {
		return nil, ParserError{"Invalid request body"}
	}

	return &ParserOutput{
		ResourceMethod: method,
		Payload:        payload,
	}, nil
}

func determineMethod(methodString string) error {
	switch methodString {
	case AuthLogin,
		TaskGetAll,
		TaskCreate,
		TaskDelete,
		TaskListGetAll,
		TaskListCreate,
		TaskListDelete:
		return nil
	default:
		return ParserError{"Method unknown"}
	}
}
