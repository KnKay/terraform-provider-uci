package ssh_helper

import (
	"bytes"
	"log"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	host   string
	config *ssh.ClientConfig
	client *ssh.Client
}

func NewClient(config *ssh.ClientConfig, host string) (t *SshClient, err error) {
	t = &SshClient{
		config: config,
		host:   host,
	}
	t.client, err = ssh.Dial("tcp", t.host, t.config)
	return
}

func (t *SshClient) RunCommand(config *ssh.ClientConfig, command string) (reply string, err error) {
	t.client, err = ssh.Dial("tcp", t.host, t.config)
	session, err := t.client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	var cmd_er bytes.Buffer
	session.Stderr = &cmd_er
	var body bytes.Buffer
	session.Stdout = &body
	err = session.Run(command)
	if err != nil {
		log.Fatalln("Unable to run command: " + err.Error())
	}
	error_out := cmd_er.String()
	if error_out != "" {
		log.Println(error_out)
	}
	reply = body.String()
	return
}
