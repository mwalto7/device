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
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"os"
	"time"
)

var TimeoutError = errors.New("session timed out")

// Device represents an SSH client.
type Device struct {
	*ssh.Client
}

// Dial creates a client connection to a remote device.
func Dial(addr string, config *ssh.ClientConfig) (*Device, error) {
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial")
	}
	return &Device{client}, nil
}

// Run creates a new session, starts a remote shell, and runs the
// specified commands. The combined output of the remote shell's standard
// output and standard error is returned.
func (d *Device) Run(cmds ...string) ([]byte, error) {
	session, err := d.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}
	defer session.Close()

	stdin, stdout, stderr, err := pipeIO(session)
	if err != nil {
		return nil, err
	}
	defer stdin.Close()

	if err := session.Shell(); err != nil {
		return nil, errors.Wrap(err, "failed to start remote shell")
	}
	for _, cmd := range cmds {
		if _, err := io.WriteString(stdin, fmt.Sprintf("%s\n", cmd)); err != nil {
			return nil, errors.Wrapf(err, "failed to run %q", cmd)
		}
	}
	wait := make(chan error, 1)
	go func(wait chan<- error) {
		wait <- session.Wait()
	}(wait)
	select {
	case <-wait:
		// TODO: Handle error value returned from `wait`.
		// TODO: Consider returning the output of stdout and stderr if an error occurs.
		//
		// if waitErr != nil {
		//     switch exitErr := waitErr.(type) {
		//	   case *ssh.ExitError:
		//         // TODO: Handle exit error.
		//	   case *ssh.ExitMissingError:
		//         // TODO: Handle missing exit error.
		//	   default:
		//		   return nil, exitErr
		//     }
		// }
		output, err := ioutil.ReadAll(io.MultiReader(stdout, stderr))
		if err != nil {
			return nil, errors.Wrap(err, "failed to read stdout and stderr")
		}
		return output, nil
	case <-time.After(5 * time.Second):
		return nil, TimeoutError
	}
}

// pipeIO creates pipes a remote shell's standard input, standard output,
// and standard error.
func pipeIO(session *ssh.Session) (stdin io.WriteCloser, stdout, stderr io.Reader, err error) {
	stdin, err = session.StdinPipe()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to create pipe to stdin")
	}
	stdout, err = session.StdoutPipe()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to create pipe to stdout")
	}
	stderr, err = session.StderrPipe()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to create pipe to stderr")
	}
	return
}

var NoAuthMethodsError = errors.New("no authentication methods specified")

// NewClientConfig is a convenience function for configuring the SSH client.
// At least one authentication method must be specified. By default the
// configuration accepts all host connections, but it is recommended to use
// the `AllowKnownHosts` Option.
func NewClientConfig(user string, opts ...Option) (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, err
		}
	}
	if len(config.Auth) == 0 {
		return nil, NoAuthMethodsError
	}
	return config, nil
}

// Option defines a function used to set the fields of a client configuration.
type Option func(*ssh.ClientConfig) error

// Password adds password authentication method to a client configuration.
func Password(password string) Option {
	return func(config *ssh.ClientConfig) error {
		if password == "" {
			fmt.Fprint(os.Stderr, "Password: ")
			pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr)
			password = string(pass)
		}
		config.Auth = append(config.Auth, ssh.Password(password))
		return nil
	}
}

// PrivateKey adds public key authentication method to a client configuration.
func PrivateKey(privateKeys ...string) Option {
	return func(config *ssh.ClientConfig) error {
		var signers []ssh.Signer
		for _, privateKey := range privateKeys {
			key, err := ioutil.ReadFile(privateKey)
			if err != nil {
				return errors.Wrap(err, "unable to read private key")
			}
			signer, err := ssh.ParsePrivateKey(key)
			if err != nil {
				return errors.Wrap(err, "unable to parse private key")
			}
			signers = append(signers, signer)
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signers...))
		return nil
	}
}

// AllowKnowHosts allows connecting only to hosts in the local known_hosts file.
func AllowKnowHosts(knownHosts string) Option {
	return func(config *ssh.ClientConfig) error {
		callback, err := knownhosts.New(knownHosts)
		if err != nil {
			return err
		}
		config.HostKeyCallback = callback
		return nil
	}
}

// Timeout sets the timeout duration for connecting to a remote host.
func Timeout(d time.Duration) Option {
	return func(config *ssh.ClientConfig) error {
		config.Timeout = d
		return nil
	}
}

// Ciphers appends the specified ciphers to a client configuration.
func Ciphers(ciphers ...string) Option {
	return func(config *ssh.ClientConfig) error {
		config.SetDefaults()
		config.Ciphers = append(config.Ciphers, ciphers...)
		return nil
	}
}
