package dmxusbpromk2

import (
	"fmt"

	"github.com/google/gousb"
	"github.com/oliread/usbdmx"
)

// DMXController a real world Enttec DMX USB Pro Mk2 device to handle comms
type DMXController struct {
	vid      uint16
	pid      uint16
	channels []byte
	packet   []byte

	ctx               *gousb.Context
	device            *gousb.Device
	output            *gousb.OutEndpoint
	input             *gousb.InEndpoint
	outputInterfaceID int
	inputInterfaceID  int
}

// NewDMXController helper function for creating a new DMX USB PRO Mk2 controller
func NewDMXController(conf usbdmx.ControllerConfig) DMXController {
	d := DMXController{}

	d.channels = make([]byte, 512)
	d.packet = make([]byte, 518)

	d.vid = conf.VID
	d.pid = conf.PID
	d.outputInterfaceID = conf.OutputInterfaceID
	d.inputInterfaceID = conf.InputInterfaceID
	d.ctx = conf.Context

	return d
}

// GetProduct returns a device product name
func (d *DMXController) GetProduct() (info string, err error) {
	info, err = d.device.Product()
	return info, err
}

// GetSerial returns a device serial number
func (d *DMXController) GetSerial() (info string, err error) {
	info, err = d.device.SerialNumber()
	return info, err
}

// Connect handles connection to a Enttec DMX USB Pro Mk2 controller
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
	cfg, err := d.device.Config(1)
	if err != nil {
		return err
	}

	intf, err := cfg.Interface(0, 0)
	if err != nil {
		return err
	}

	d.output, err = intf.OutEndpoint(d.outputInterfaceID)
	if err != nil {
		return err
	}

	commsInterface, _, err := d.device.DefaultInterface()
	if err != nil {
		return err
	}

	d.output, err = commsInterface.OutEndpoint(d.outputInterfaceID)
	if err != nil {
		return err
	}

	// d.input, err = commsInterface.InEndpoint(d.inputInterfaceID)
	// if err != nil {
	// 	return err
	// }

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

	d.output.Write([]byte{0x7E, 0x0A, 0x00, 0x00, 0xE7})

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

// Render sends channel data to fixtures
func (d *DMXController) Render() error {
	// ENTTEC USB DMX PRO Start Message
	d.packet[0] = 0x7E

	// Set our protocol
	d.packet[1] = 0x06
	d.packet[2] = 0x01

	// Wat?
	d.packet[3] = 0x02
	d.packet[4] = 0x00

	// Set DMX Data
	for i := 0; i < 512; i++ {
		d.packet[i+5] = d.channels[i]
	}

	// ENTTEC USB DMX PRO End Message
	d.packet[517] = 0xE7

	if _, err := d.output.Write(d.packet); err != nil {
		return err
	}

	return nil
}
