package launchpad

import (
	"bytes"
	"strings"

	"github.com/pkg/errors"
	"github.com/rakyll/portmidi"
)

// TODO: Fix all the TODOs
// TODO: Add DeviceInquiry, VersionInquiry, and SetToBootloader
// TODO: Finalize API
// TODO: Add examples
// TODO: Handle input
// TODO: Support more devices

// This package allows for control of the Launchpad Mk2 via SysEx messages.
// There is an alternate method of controlling the Launchpad using notes, but
// there doesn't seem to be a large advantage to using it, as not all operations
// are available using notes, but all operations are available with SysEx
// messages.

const (
	devicePrefix     = "Launchpad "
	inputBufferSize  = 8
	outputBufferSize = 8
)

// FaderType is used when setting fader information.
type FaderType byte

// There are two main fader types, Volume and Pan. As far as I can tell, these
// only change how the LEDs display.
const (
	VolumeFader FaderType = 1
	PanFader              = 2
)

// LayoutType can be used to change how LEDs are addressed.
type LayoutType byte

// SessionLayout has the one of the simplest LED addressing strategies
// as moving to the right increases the LED ID by 1 and moving up
// increases by 10. Additionally, the AbletonLayout is reserved and
// should not be used.
const (
	SessionLayout LayoutType = 0
	User1Layout              = 1
	User2Layout              = 2
	AbletonLayout            = 3
	VolumeLayout             = 4
	PanLayout                = 5
)

// It's not possible to make const byte slices, so we fall back to just making
// them variables and not messing with them.
var (
	sysExHeader = []byte{240, 00, 32, 41, 2, 24}
	sysExFooter = []byte{247}

	// Not currently used. Also note that the 127 can be replaced with
	// a specific device ID (0 through F) to ensure a response will
	// only happen if the ID matches.
	deviceInquiryPayload = []byte{240, 126, 127, 6, 1, 247}

	// Not currently used.
	versionInquiryPayload = []byte{240, 0, 32, 41, 0, 112, 247}

	// Not currently used.
	setToBootloaderPayload = []byte{240, 0, 32, 41, 0, 113, 0, 105, 247}
)

// Launchpad is an abstraction around the Launchpad as a midi device, adding
// simple functions so users don't need to know the magic under the hood.
type Launchpad struct {
	input  *portmidi.Stream
	output *portmidi.Stream
}

// NewLaunchpad loops for a device having a name with the prefix "Launchpad "
// and the suffix given. It will return an error if anything fails to open.
func NewLaunchpad(suffix string) (*Launchpad, error) {
	var err error
	ret := &Launchpad{}

	for x := 0; x < portmidi.CountDevices(); x++ {
		id := portmidi.DeviceID(x)
		info := portmidi.Info(id)
		if !strings.HasPrefix(info.Name, devicePrefix) {
			continue
		}

		if info.Name[len(devicePrefix):] != suffix {
			continue
		}

		if info.IsInputAvailable {
			ret.input, err = portmidi.NewInputStream(id, inputBufferSize)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to open input stream")
			}
		}
		if info.IsOutputAvailable {
			ret.output, err = portmidi.NewOutputStream(id, outputBufferSize, 0)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to open output stream")
			}
		}
	}

	if ret.input == nil {
		return nil, errors.New("Missing input stream")
	}
	if ret.output == nil {
		return nil, errors.New("Missing output stream")
	}

	return ret, nil
}

func (lp *Launchpad) SetLED(ledID byte, color byte) error {
	// TODO: This can be repeated up to 80 times
	return lp.writeSysExBytes([]byte{10, ledID, color})
}

func (lp *Launchpad) SetLEDRGB(ledID byte, r, g, b byte) error {
	// TODO: This can be repeated up to 80 times
	return lp.writeSysExBytes([]byte{11, ledID, r, g, b})
}

func (lp *Launchpad) SetLEDColumn(colID byte, color byte) error {
	// TODO: This can be repeated up to 9 times
	return lp.writeSysExBytes([]byte{12, colID, color})
}

func (lp *Launchpad) SetLEDRow(rowID byte, color byte) error {
	// TODO: This can be repeated up to 9 times
	return lp.writeSysExBytes([]byte{13, rowID, color})
}

func (lp *Launchpad) SetAllLEDs(color byte) error {
	return lp.writeSysExBytes([]byte{14, color})
}

func (lp *Launchpad) ScrollText(text string, color byte, speed byte, loop bool) error {
	// TODO: Determine how to handle UTF8 text
	var loopByte byte
	if loop {
		loopByte = 1
	}
	return lp.writeSysExBytes(append([]byte{20, color, loopByte, speed}, []byte(text)...))
}

func (lp *Launchpad) SelectLayout(layout LayoutType) error {
	return lp.writeSysExBytes([]byte{34, byte(layout)})
}

func (lp *Launchpad) FlashLED(ledID byte, color byte) error {
	// TODO: This can be repeated up to 80 times
	return lp.writeSysExBytes([]byte{35, 0, ledID, color})
}

func (lp *Launchpad) PulseLED(ledID byte, color byte) error {
	// TODO: This can be repeated up to 80 times
	return lp.writeSysExBytes([]byte{40, 0, ledID, color})
}

func (lp *Launchpad) FaderSetup(faderID byte, faderType FaderType, color, initialValue byte) error {
	// TODO: This can be repeated up to 8 times
	return lp.writeSysExBytes([]byte{43, faderID, byte(faderType), color, initialValue})
}

func (lp *Launchpad) Close() error {
	err := lp.input.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close input stream")
	}

	err = lp.output.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close output stream")
	}

	return nil
}

func (lp *Launchpad) writeSysExBytes(data []byte) error {
	out := bytes.NewBuffer(sysExHeader)
	out.Write(data)
	out.Write(sysExFooter)
	return lp.output.WriteSysExBytes(portmidi.Time(), out.Bytes())
}
