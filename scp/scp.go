package scp

import (
	"bytes"
	"context"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

var client scp.Client
var err error

func Connection(conn *ssh.Client) error {
	client, err = scp.NewClientBySSH(conn)
	if err != nil {
		return err
	}
	return nil
}

func WriteFile(path string, body []byte, exec bool) error {
	bs := bytes.NewReader(body)
	permissons := "0664"
	if exec {
		permissons = "0755"
	}
	err := client.CopyFile(context.Background(), bs, path, permissons)
	return err
}
func CopyFile(src, dst string, exec bool) error {
	permissons := "0600"
	if exec {
		permissons = "0755"
	}
	file, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	bs := bytes.NewReader(file)
	err = client.CopyFile(context.Background(), bs, dst, permissons)
	return err
}
