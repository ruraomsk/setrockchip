package command

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

var conn *ssh.Client

func Connection(c *ssh.Client) {
	conn = c
}
func AnyCommand(command string) error {
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	return session.Run(command)
}
func KillProc(name string) error {
	res := make([]string, 0)
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run(fmt.Sprintf("ps -e | grep %s", name))
	scanner := bufio.NewScanner(bytes.NewReader(b.Bytes()))
	for scanner.Scan() {
		s := scanner.Text()
		s = strings.TrimLeft(s, " ")
		bs := strings.Split(s, " ")
		res = append(res, bs[0])
	}
	for _, v := range res {
		kill(v)
	}
	return nil
}
func kill(kp string) {
	session, err := conn.NewSession()
	if err != nil {
		return
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run(fmt.Sprintf("kill -9 %s\n", kp))
}
func DeleteFile(path string) error {
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run(fmt.Sprintf("rm %s", path))
	return nil
}
func CreateDir(path string) error {
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(fmt.Sprintf("mkdir -p %s", path)); err != nil {
		return err
	}
	return nil
}
func Permisson(path string, file string) error {
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(fmt.Sprintf("cd %s ; chmod 755 %s ", path, file)); err != nil {
		return err
	}
	return nil
}

func DeleteDir(path string) error {
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run(fmt.Sprintf("rm -r %s", path))
	return nil
}

func GetSystem() (string, error) {
	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("uname -a"); err != nil {
		return "", err
	}
	bs := b.String()
	bbs := strings.Split(bs, " ")
	return bbs[1], nil
}
