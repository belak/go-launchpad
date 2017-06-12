package launchpad

import (
	"errors"

	"github.com/rakyll/portmidi"
)

const (
	inputBufferSize  = 8
	outputBufferSize = 8
)

type midiDevice struct {
	input  *portmidi.Stream
	output *portmidi.Stream
}

func getMidiDevice(name string) (*midiDevice, error) {
	var inputID portmidi.DeviceID = -1
	var outputID portmidi.DeviceID = -1

	for x := 0; x < portmidi.CountDevices(); x++ {
		id := portmidi.DeviceID(x)
		info := portmidi.Info(id)

		if info.Name != name {
			continue
		}

		if info.IsInputAvailable {
			inputID = id
		}
		if info.IsOutputAvailable {
			outputID = id
		}
	}

	if inputID == -1 {
		return nil, errors.New("Missing input stream")
	}
	if outputID == -1 {
		return nil, errors.New("Missing output stream")
	}

	return newMidiDevice(inputID, outputID)
}

func newMidiDevice(inputID, outputID portmidi.DeviceID) (*midiDevice, error) {
	var err error
	ret := &midiDevice{}
	ret.input, err = portmidi.NewInputStream(inputID, inputBufferSize)
	if err != nil {
		return nil, err
	}

	ret.output, err = portmidi.NewOutputStream(outputID, outputBufferSize, 0)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *midiDevice) WriteShort(status, data1, data2 int64) error {
	return m.output.WriteShort(status, data1, data2)
}

func (m *midiDevice) Read(count int) ([]portmidi.Event, error) {
	return m.input.Read(count)
}

func (m *midiDevice) Close() error {
	err := m.input.Close()
	if err != nil {
		return err
	}

	return m.output.Close()
}
