package launchpad

import "errors"

type LaunchpadMk1 struct {
	midi *midiDevice
}

func NewLaunchpadMk1() (*LaunchpadMk1, error) {
	var err error
	lp := &LaunchpadMk1{}
	lp.midi, err = getMidiDevice("Launchpad Mini")
	if err != nil {
		return nil, err
	}
	return lp, nil
}

func (lp *LaunchpadMk1) SetLED(x, y int64, color int64) error {
	// For the top row, there's a special case
	if y == 8 {
		return lp.midi.WriteShort(176, x+104, color)
	}

	ledID := (16 * (7 - y)) + x
	return lp.midi.WriteShort(144, ledID, color)
}

func (lp *LaunchpadMk1) SetAll(color int64) error {
	// TODO: Make this work with all colors and/or check bounds of color
	return lp.midi.WriteShort(176, 0, color)
}

func (lp *LaunchpadMk1) Clear() error {
	return lp.midi.WriteShort(176, 0, 0)
}

func (lp *LaunchpadMk1) GetEvent() (*Event, bool, error) {
	events, err := lp.midi.Read(1)
	if err != nil {
		return nil, false, err
	}

	if len(events) != 1 {
		return nil, false, nil
	}

	event := &Event{
		Velocity: events[0].Data2,
	}

	if events[0].Status == 176 {
		event.X = events[0].Data1 - 104
		event.Y = 8
		return event, true, nil
	} else if events[0].Status == 144 {
		event.X = events[0].Data1 % 16
		event.Y = 7 - (events[0].Data1 / 16)
		return event, true, nil
	}

	return nil, false, errors.New("Unhandled event")
}

func (lp *LaunchpadMk1) Close() error {
	return lp.midi.Close()
}
