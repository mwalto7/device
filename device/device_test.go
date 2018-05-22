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
	"testing"
	"log"
	"fmt"
	"github.com/mwalto7/device/device"
)

func TestDial(t *testing.T) {
	// TODO: Test the Dial function.
}

func TestDevice_SendCmds(t *testing.T) {
	// TODO: Test the Device.SendCmds method.
}

func TestDevice_Close(t *testing.T) {
	// TODO: Test the Device.Close method.
}

func BenchmarkDial(b *testing.B) {
	// TODO: Benchmark Dial function.
}

func BenchmarkDevice_SendCmds(b *testing.B) {
	// TODO: Benchmark Device.SendCmds method.
}

func ExampleDevice_SendCmds() {
	// Establish an SSH connection to a network device.
	netDev, err := device.Dial("127.0.0.1", "22", "user", "password")
	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}
	defer netDev.Close()

	// Send configuration commands and capture the output.
	cmds := []string{"conf t", "int Gi1/0/1", "description hello_world", "exit", "exit"}
	output, err := netDev.SendCmds(cmds...)
	if err != nil {
		log.Fatalf("Failed to run: %v\n", err)
	}
	fmt.Println(string(output))
}
