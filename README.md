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
    "net"
    "time"
)

func main() {
    // Create a new client configuration.
    config, err := device.NewClientConfig(
    	"user",
    	device.PrivateKey("~/.ssh/id_rsa"),
    	device.Password("password"),
    	device.AllowKnowHosts("~/.ssh/known_hosts"),
    	device.Timeout(5 * time.Second),
    )
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
```
