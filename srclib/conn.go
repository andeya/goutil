package srclib

import (
	"net"

	"github.com/henrylee2cn/srcpool"
)

// Resource is a resource that can be stored in the Pool.
type Net struct {
	net.Conn
	avatar *srcpool.Avatar
}

var _ srcpool.Resource = new(Net)

// New wraps a net.Conn
func NewNet(conn net.Conn) *Net {
	return &Net{
		Conn: conn,
	}
}

// SetAvatar stores the contact with pool
// Do not call it yourself, it is only called by (*Pool).get, and will only be called once
func (n *Net) SetAvatar(avatar *srcpool.Avatar) {
	n.avatar = avatar
}

// GetAvatar gets the contact with pool
// Do not call it yourself, it is only called by (*Pool).Put
func (n *Net) GetAvatar() *srcpool.Avatar {
	return n.avatar
}

// Close closes the original source
// No need to call it yourself, it is only called by (*Avatar).close
func (n *Net) Close() error {
	return n.Conn.Close()
}
