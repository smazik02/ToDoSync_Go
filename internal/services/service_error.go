package services

import (
	"encoding/json"
	"fmt"
)

type ResponseSource = string

const (
	AUTH ResponseSource = "AUTH"
	T    ResponseSource = "T"
	TL   ResponseSource = "TL"
)

type ServiceError struct {
	message string
	source  ResponseSource
}

func (e ServiceError) Error() string {
	body := map[string]string{"message": e.message, "source": e.source}
	bodyJSON, _ := json.Marshal(body)
	return fmt.Sprintf("FAIL\n%s\n\n", string(bodyJSON))
}
