package copyer

import (
	"fmt"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var client *sftp.Client
var err error

func Connection(conn *ssh.Client) error {
	client, err = sftp.NewClient(conn)
	if err != nil {
		return err
	}
	return nil
}

func WriteFile(path string, body []byte) error {
	f, err := client.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(body); err != nil {
		return err
	}
	return nil
}
func Chmod(path string) error {
	err := client.Chmod(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
func CopyFile(src, dst string) error {
	file, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	wf, err := client.Create(dst)
	if err != nil {
		return err
	}
	n, err := wf.Write(file)
	if err != nil {
		return err
	}
	if n != len(file) {
		return fmt.Errorf("write %d bytes need %d", n, len(file))
	}
	wf.Close()
	return nil
}
