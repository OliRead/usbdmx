package mock

import (
	"fmt"

	"github.com/oliread/usbdmx"
)

// DMXController A mock controller for sending DMX commands to
type DMXController struct {
	channels []byte
}

// NewDMXController helper function for creating a new mock controller
func NewDMXController(conf usbdmx.ControllerConfig) DMXController {
	d := DMXController{}

	d.channels = make([]byte, 512)

	return d
}

// Connect handles connection to a mock DMX controller
func (d *DMXController) Connect() error {
	return nil
}

// SetChannel sets a single DMX channel value
func (d *DMXController) SetChannel(index int16, data byte) error {
	if index < 1 || index > 512 {
		return fmt.Errorf("Index %d out of range, must be between 1 and 512", index)
	}

	d.channels[index-1] = data

	return nil
}

// GetChannel returns the value of a single DMX channel
func (d *DMXController) GetChannel(index int16) (byte, error) {
	if index < 1 || index > 512 {
		return 0, fmt.Errorf("Index %d out of range, must be between 1 and 512", index)
	}

	return d.channels[index-1], nil
}

// Render sends channel data to fixtures, in this case prints it to stdout
func (d *DMXController) Render() error {
	fmt.Print(d.channels)

	return nil
}
