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

// Package device_tests contains tests and examples for the device package.
package device_test

import (
	"github.com/mwalto7/device/device"
	"log"
	"fmt"
	"io/ioutil"
)

func ExampleConfiguration() {
	// Establish an SSH connection to a device and defer
	// closing the connection.
	netDev, err := device.Dial("127.0.0.1", "22", "user", "password")
	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}
	defer netDev.Close()

	// Send configuration commands and capture the output.
	output, err := netDev.SendCmds("conf t", "int Gi1/0/1", "description hello_world")
	if err != nil {
		log.Fatalf("Failed to run: %v\n", err)
	}

	// Read all the contents of the session.
	contents, _ := ioutil.ReadAll(output)
	fmt.Println(string(contents))
}
