package shared

import (
	"net"
	"strings"
)

type Client struct {
	Connection net.Conn
	Buffer     strings.Builder
	UserId     int
}

func (c Client) IsLoggedIn() bool {
	return c.UserId != -1
}
