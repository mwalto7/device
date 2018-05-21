# device

Device is a pure Go library for configuring network devices.

Full documentation: [GoDoc](https://godoc.org/github.com/mwalto7/device/device)

`go get "github.com/mwalto7/device/device"`

Example:

```go
package main

import (
    "fmt"
    "github.com/mwalto7/device/device"
    "io/ioutil"
    "log"
)

func main() {
    // Establish an SSH connection to a device.
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
```
