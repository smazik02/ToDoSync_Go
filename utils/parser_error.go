package utils

import (
	"encoding/json"
	"fmt"
)

type ParserError struct {
	message string
}

func (e ParserError) Error() string {
	body := map[string]string{"body": e.message}
	bodyJSON, _ := json.Marshal(body)
	return fmt.Sprintf("FAIL\n%s\n\n", string(bodyJSON))
}
