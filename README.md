# go-launchpad
A small go library for using the Novation Launchpad MK2 as an LED display and input device

## Notes

Currently, portmidi must be initialized by the user.

## Example

```go
package main

import (
	"fmt"
	"log"

	launchpad "github.com/belak/go-launchpad"
	"github.com/rakyll/portmidi"
)

func main() {
	err := portmidi.Initialize()
	if err != nil {
		log.Fatalln(err)
	}
	defer portmidi.Terminate()

	lp, err := launchpad.NewLaunchpadMk1()
	if err != nil {
		log.Fatalln(err)
	}
	defer lp.Close()

	for {
		event, ok, err := lp.GetEvent()
		if err != nil {
			log.Fatalln(err)
		}
		if !ok {
			continue
		}

		fmt.Printf("%+v\n", event)
		if event.Velocity == 0 {
			lp.SetLED(event.X, event.Y, 0)
		} else {
			lp.SetLED(event.X, event.Y, 43)
		}
	}
}
```
