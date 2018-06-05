package usbdmx

import (
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/google/gousb"
)

// Controller Generic interface for all USB DMX controllers
type Controller interface {
	Connect() (err error)
	SetChannel(channel int16, value byte) error
	GetChannel(channel int16) (byte, error)
	Render() error
}

// ControllerConfig configuration for controlling device
type ControllerConfig struct {
	VID               uint16 `toml:"vid"`
	PID               uint16 `toml:"pid"`
	OutputInterfaceID int    `toml:"output_interface_id"`
	DebugLevel        int    `toml:"debug_level"`

	Context *gousb.Context
}

// ReadConfigFile reads device configuration information from file
func ReadConfigFile(path string) (ControllerConfig, error) {
	type raw struct {
		VID               string `toml:"VID"`
		PID               string `toml:"PID"`
		OutputInterfaceID int    `toml:"outputInterfaceID"`
		DebugLevel        int    `toml:"debugLevel"`
	}
	rawConf := raw{}
	conf := ControllerConfig{}

	if _, err := toml.DecodeFile(path, &rawConf); err != nil {
		return conf, err
	}

	vid, err := strconv.ParseUint(rawConf.VID, 16, 16)
	if err != nil {
		return conf, err
	}

	pid, err := strconv.ParseUint(rawConf.PID, 16, 16)
	if err != nil {
		return conf, err
	}

	conf.VID = uint16(vid)
	conf.PID = uint16(pid)
	conf.OutputInterfaceID = rawConf.OutputInterfaceID
	conf.DebugLevel = rawConf.DebugLevel

	return conf, nil
}

// NewConfig helper function for creating a new ControllerConfig
func NewConfig(vid, pid uint16, outputInterfaceID, debugLevel int) ControllerConfig {
	return ControllerConfig{
		VID:               vid,
		PID:               pid,
		OutputInterfaceID: outputInterfaceID,
		DebugLevel:        debugLevel,
	}
}

// GetUSBContext gets a gousb/context for a given configuration
func (c *ControllerConfig) GetUSBContext() {
	c.Context = gousb.NewContext()
	c.Context.Debug(c.DebugLevel)
}
