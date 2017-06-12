package launchpad

// Launchpad exposes a common interface between different launchpad devices.
type Launchpad interface {
	// SetLED will set the given LED to the given color. Note that the
	// coordinates start in the bottom left to be opposite the "missing" led in
	// the opposite corner.
	SetLED(x, y int64, color int64) error

	// SetAll sets all LEDs to the given color.
	SetAll(color int64) error

	// Clear will clear all LEDs.
	Clear() error

	// GetEvent polls for an event and returns it along with a bool to show if
	// there was an actual event.
	GetEvent() (*Event, bool)

	// Close will clean up any resources being used by the midi device.
	Close() error
}

// Event represents a key press or release. X and Y have an origin in the bottom
// left, similar to SetLED on the Launchpad interface.
type Event struct {
	X, Y     int64
	Velocity int64
}
