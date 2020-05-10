package ft232

import (
	"fmt"

	"github.com/google/gousb"
	"github.com/oliread/usbdmx"
)

// DMXController a real world FT232 DMX Controller to handle comms
type DMXController struct {
	vid      uint16
	pid      uint16
	channels []byte
	packet   []byte

	ctx               *gousb.Context
	device            *gousb.Device
	output            *gousb.OutEndpoint
	outputInterfaceID int
}

// NewDMXController helper function for creating a new ft232 controller
func NewDMXController(conf usbdmx.ControllerConfig) DMXController {
	d := DMXController{}

	d.channels = make([]byte, 512)
	d.packet = make([]byte, 513)

	d.vid = conf.VID
	d.pid = conf.PID
	d.outputInterfaceID = conf.OutputInterfaceID
	d.ctx = conf.Context

	return d
}

// Connect handles connectio to a mock DMX controller
func (d *DMXController) Connect() error {
	// try to connect to device
	device, err := d.ctx.OpenDeviceWithVIDPID(gousb.ID(d.vid), gousb.ID(d.pid))
	if err != nil {
		return err
	}
	d.device = device

	// make this device ours, even if it is being used elsewhere
	if err := d.device.SetAutoDetach(true); err != nil {
		return err
	}

	// open communications
	commsInterface, _, err := d.device.DefaultInterface()
	if err != nil {
		return err
	}

	d.output, err = commsInterface.OutEndpoint(d.outputInterfaceID)
	if err != nil {
		return err
	}

	// Send our control headers for this device
	d.device.Control(0x40, 0x00, 0x00, 0x00, nil)
	d.device.Control(0x40, 0x03, 0x4138, 0x00, nil)
	d.device.Control(0x40, 0x00, 0x00, 0x00, nil)
	d.device.Control(0x40, 0x04, 0x1008, 0x00, nil)
	d.device.Control(0x40, 0x02, 0x00, 0x00, nil)
	d.device.Control(0x40, 0x03, 0x000c, 0x00, nil)
	d.device.Control(0x40, 0x00, 0x0001, 0x00, nil)
	d.device.Control(0x40, 0x00, 0x0002, 0x00, nil)
	d.device.Control(0x40, 0x01, 0x0200, 0x00, nil)
	d.device.Control(0x40, 0x04, 0x5008, 0x00, nil)
	d.device.Control(0x40, 0x00, 0x0002, 0x00, nil)
	d.device.Control(0x40, 0x04, 0x1008, 0x00, nil)

	return nil
}

func (d *DMXController) Close() error {
	device := d.device
	if device == nil {
		return fmt.Errorf("not connected")
	}
	return device.Close()
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
	for i := 0; i < 512; i++ {
		d.packet[i+1] = d.channels[i]
	}

	d.device.Control(0x40, 0x04, 0x5008, 0x00, nil)
	d.device.Control(0x40, 0x04, 0x1008, 0x00, nil)
	if _, err := d.output.Write(d.packet); err != nil {
		return err
	}

	return nil
}
