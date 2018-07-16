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

// Package device_test contains tests, benchmarks, and examples for package
// device.
package device_test

import (
	"fmt"
	"github.com/mwalto7/netconfig/device"
	"log"
	"net"
	"time"
)

func ExampleDevice_Run() {
	// Create a new client configuration.
	config, err := device.NewClientConfig("user", device.Password("password"))
	if err != nil {
		log.Fatal(err)
	}

	// Establish a client connection to a host and defer closing the connection.
	netdev, err := device.Dial(net.JoinHostPort("host", "port"), config)
	if err != nil {
		log.Fatal(err)
	}
	defer netdev.Close()

	// Run the commands and capture the session output.
	var cmds []string
	output, err := netdev.Run(cmds...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
}

func ExampleNewClientConfig() {
	// Example client configuration that specifies public key and password
	// authentication methods and allows only connections to known hosts.
	config, err := device.NewClientConfig(
		// Username for device login
		"user",

		// Only connect to hosts in known_hosts
		device.AllowKnowHosts("~/.ssh/known_hosts"),

		// Use key authentication
		device.PrivateKey("~/.ssh/id_rsa"),

		// Use password authentication as backup.
		device.Password("password"),

		// Timeout if establishing the connection exceeds the duration
		device.Timeout(5*time.Second),

		// Add additional ciphers supported by this device
		device.Ciphers("aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"),
	)
	if err != nil {
		log.Fatal(err)
	}
	device.Dial("addr", config)
}
