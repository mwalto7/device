// Copyright © 2018 Mason Walton <dev.mwalto7@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package device implements an SSH client for configuring network devices.
// It defines a type, Device, with methods for establishing an SSH connection,
// sending configuration commands, and retrieving the contents of the device's
// standard output and standard error.
package device

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"time"
)

// Device represents a network device that can be
// configured over SSH.
type Device struct {
	*ssh.Client // the underlying SSH client
}

// Dial returns an SSH client connection to the specified
// host using password authentication. If a connection
// is not established within 5 seconds, Dial will timeout.
func Dial(host, port, user, password string) (*Device, error) {
	// Initialize the client configuration.
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// Set the defaults for the client configuration and
	// append supported ciphers for older Cisco and HP
	// devices.
	config.SetDefaults()
	config.Ciphers = append(config.Ciphers, "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc")

	// Establish an SSH connection to the device.
	client, err := ssh.Dial("tcp", net.JoinHostPort(host, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	return &Device{client}, nil
}

// SendCmds starts an SSH session, writes `cmds` to its
// standard input, waits for the remote commands to exit,
// then returns the contents of standard output and
// standard error.  If the remote commands wait for more
// than 10 seconds, the session will be terminated.
// SendCmds is safe to use concurrently.
func (d *Device) SendCmds(cmds ...string) ([]byte, error) {
	// Create a new session
	session, err := d.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Create pipes to stdin, stdout, and stderr
	stdin, stdout, stderr, err := setIO(session)
	if err != nil {
		return nil, fmt.Errorf("failed to setup IO: %v", err)
	}
	defer stdin.Close()

	// Start the remote shell
	if err := startShell(session); err != nil {
		return nil, err
	}

	// Write the commands to stdin
	for _, cmd := range cmds {
		if _, err := stdin.Write([]byte(cmd + "\n")); err != nil {
			return nil, fmt.Errorf("failed to run: %v", err)
		}
	}

	// Wait for remote commands to exit or timeout.
	exit := make(chan error, 1)
	go func(exit chan<- error) {
		exit <- session.Wait()
	}(exit)
	timeout := time.After(30 * time.Second)
	for {
		select {
		case <-exit:
			// TODO: Handle error value of exit channel.
			return ioutil.ReadAll(io.MultiReader(stdout, stderr))
		case <-timeout:
			return nil, fmt.Errorf("session timed out")
		}
	}
}

// Close closes the SSH client connection.
func (d *Device) Close() error {
	return d.Client.Close()
}

// setIO returns pipes connected to the remote shell's standard
// input, standard output, and standard error.
func setIO(session *ssh.Session) (stdin io.WriteCloser, stdout, stderr io.Reader, err error) {
	stdin, err = session.StdinPipe()
	if err != nil {
		return
	}
	stdout, err = session.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err = session.StderrPipe()
	if err != nil {
		return
	}
	return
}

// startShell requests a pseudo terminal and starts the remote shell.
func startShell(session *ssh.Session) error {
	if err := session.RequestPty("vt100", 0, 0, ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed: 14.4k baud
		ssh.TTY_OP_OSPEED: 14400, // output speed: 14.4k baud
	}); err != nil {
		return fmt.Errorf("failed to request pseudo terminal: %v", err)
	}
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start remote shell: %v", err)
	}
	return nil
}
