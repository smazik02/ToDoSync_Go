package shared

import (
	"bytes"
	"net"
)

type Client struct {
	Connection net.Conn
	Buffer     bytes.Buffer
	UserId     int
}

func (c Client) IsLoggedIn() bool {
	return c.UserId != -1
}
