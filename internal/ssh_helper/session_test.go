package ssh_helper_test

import (
	"testing"
	"time"

	"github.com/KnKay/terraform-provider-uci/internal/ssh_helper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

const (
	username = "root"
	password = "test123"
	host     = "192.168.1.199:22"
)

func TestCommand(t *testing.T) {
	conf := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// Non-production only
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	command := "echo OpenWrt"
	client, err := ssh_helper.NewClient(conf, host)
	if err != nil {
		t.Error(err.Error())
	}
	reply, err := client.RunCommand(conf, command)
	if err != nil {
		t.Error(err.Error())
	}
	assert := assert.New(t)
	assert.Equal("OpenWrt\n", reply)
}
