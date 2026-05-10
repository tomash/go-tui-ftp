package ftp

import (
	"fmt"
	"sync"

	ftplib "github.com/jlaffaye/ftp"
)

// Client wraps the FTP connection logic
type Client struct {
	mu   sync.Mutex
	conn *ftplib.ServerConn
}

func NewClient() *Client {
	return &Client{}
}

// Connect handles authentication. Address always includes an explicit port for net.Dial.
func (c *Client) Connect(host string, port int, user, pass string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		_ = c.conn.Quit()
		c.conn = nil
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := ftplib.Dial(addr, ftplib.DialWithTimeout(ftplib.DefaultDialTimeout))
	if err != nil {
		return err
	}

	if err := conn.Login(user, pass); err != nil {
		_ = conn.Quit()
		return err
	}

	c.conn = conn
	return nil
}

// ListDir returns entry names in the given path.
func (c *Client) ListDir(path string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	entries, err := c.conn.List(path)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e != nil && e.Name != "" {
			names = append(names, e.Name)
		}
	}
	return names, nil
}

func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil
	}
	err := c.conn.Quit()
	c.conn = nil
	return err
}
