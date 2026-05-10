package ftp

import "fmt"

// Client wraps FTP operations. Will be fleshed out in Phase 2.
type Client struct {
	Host string
	Port int
}

func NewClient(host string, port int) *Client {
	return &Client{Host: host, Port: port}
}

func (c *Client) Connect(user, pass string) error {
	return fmt.Errorf("not implemented yet")
}

func (c *Client) ListDir(path string) ([]string, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// Upload/Download stubs for future implementation
func (c *Client) Upload(localPath, remoteName string) error   { return nil }
func (c *Client) Download(remotePath, localDest string) error { return nil }
