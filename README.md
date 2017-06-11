# go-launchpad
A small go library for using the Novation Launchpad MK2 as an LED display and input device

## Notes

Currently, portmidi must be initialized by the user.

## Example

```go
func main() {
	err := portmidi.Initialize()
	if err != nil {
		log.Fatalln(err)
	}
	defer portmidi.Terminate()

	lp, err := NewLaunchpad("MK2")
	if err != nil {
		log.Fatalln(err)
	}
	defer lp.Close()

	lp.ScrollText("Hello World", 40, 1, false)

	lp.SetAllLEDs(100)
	time.Sleep(3 * time.Second)
	lp.SetAllLEDs(0)
}
```
