package main

import (
	"flag"
	"log"
	"time"

	"github.com/oliread/usbdmx"
	"github.com/oliread/usbdmx/enttec/dmxusbpromk2"
)

func main() {
	// Constants, these should really be defined in the module and will be
	// as of the next release
	vid := uint16(0x0403)
	pid := uint16(0x6001)
	outputInterfaceID := flag.Int("output-id", 0, "Output interface ID for device")
	inputInterfaceID := flag.Int("input-id", 0, "Input interface ID for device")
	debugLevel := flag.Int("debug", 0, "Debug level for USB context")
	flag.Parse()

	// Create a configuration from our flags
	config := usbdmx.NewConfig(vid, pid, *outputInterfaceID, *inputInterfaceID, *debugLevel)

	// Get a usb context for our configuration
	config.GetUSBContext()

	// Create a controller and connect to it
	controller := dmxusbpromk2.NewDMXController(config)
	if err := controller.Connect(); err != nil {
		log.Fatalf("Failed to connect DMX Controller: %s", err)
	}

	// Set first three channels to zero, this assumes that our 3 channel RGB
	// fixture, like a par can, will not be showing any light. We're ignoring
	// errors but the SetChannel function will return an error if it fails to
	// write to the array
	controller.SetChannel(1, 1, 0)
	controller.SetChannel(1, 2, 0)
	controller.SetChannel(1, 3, 0)

	// Create an array of colours for our fixture to switch between
	colours := [][]byte{
		[]byte{255, 0, 0},
		[]byte{255, 255, 0},
		[]byte{0, 255, 0},
		[]byte{0, 255, 255},
		[]byte{0, 0, 255},
		[]byte{255, 0, 255},
	}

	// Create a go routine that will ensure our controller keeps sending data
	// to our fixture with a short delay. No delay, or too much delay, may cause
	// flickering in fixtures. Check the specification of your fixtures and controller
	go func(c *dmxusbpromk2.DMXController) {
		for {
			if err := controller.Render(); err != nil {
				log.Fatalf("Failed to render output: %s", err)
			}

			time.Sleep(30 * time.Millisecond)
		}
	}(&controller)

	// Create a loop that will cycle through all of the colours defined in the "colours"
	// array and set the channels on our controller. Once the channels have been set their
	// values are ouptut to stdout. Wait 2 seconds between updating our new channels
	for i := 0; true; i++ {
		colour := colours[i%len(colours)]
		controller.SetChannel(1, 1, colour[0])
		controller.SetChannel(1, 2, colour[1])
		controller.SetChannel(1, 3, colour[2])

		r, _ := controller.GetChannel(1, 1)
		g, _ := controller.GetChannel(1, 2)
		b, _ := controller.GetChannel(1, 3)

		log.Printf("Ch1: %d \t Ch2: %d \t Ch3: %d", r, g, b)
		time.Sleep(time.Second * 2)
	}
}
