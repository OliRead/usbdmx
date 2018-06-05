package mock

import (
	"testing"

	"github.com/oliread/usbdmx"
)

var config = usbdmx.ControllerConfig{}

func TestGetSetChannel(t *testing.T) {
	c := NewDMXController(config)

	if err := c.SetChannel(-1, 0); err == nil {
		t.Errorf("Expected error when setting channel at index -1, got nil")
	}

	if err := c.SetChannel(513, 0); err == nil {
		t.Errorf("Expected error when setting channel at index 513, got nil")
	}

	if err := c.SetChannel(0, 255); err != nil {
		t.Error(err)
	}

	if data, err := c.GetChannel(1); err != nil {
		t.Error(err)
	} else if data != 255 {
		t.Errorf("Expected channel 1 value to be 255, got %d", data)
	}
}
