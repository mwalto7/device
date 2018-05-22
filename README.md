# device

Device is a pure Go library for configuring network devices.

Documentation: [GoDoc](https://godoc.org/github.com/mwalto7/device/device)

Import: `go get -u github.com/mwalto7/device/device`

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
    // Establish an SSH connection to a network device.
    netdev, err := device.Dial("127.0.0.1", "22", "user", "password")
    if err != nil {
        log.Fatalf("Failed to connect: %v\n", err)
    }
    defer netdev.Close()

    // Send configuration commands and capture the output.
    cmds := []string{"conf t", "int Gi1/0/1", "description hello_world", "exit", "exit"}
    output, err := netdev.SendCmds(cmds...)
    if err != nil {
        log.Fatalf("Failed to run: %v\n", err)
    }
    fmt.Println(string(output))
}
```
