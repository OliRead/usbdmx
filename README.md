# USBDMX

USBDMX, a versatile USB DMX library written in Go for programatic show
control and special effects.

## Supported Controllers

* DIY FT232 controller *[Steven Breuls Tutorial](https://stevenbreuls.com/2013/05/diy-usb-dmx-dongle-interface-for-under-10/)*

Planned for support:

* Enttec DMX USB Pro Mk2

If you don't see your controller here, but would like it to be added, create an issue 
with the name of the device and a link to any more information.

## Quick Start

This will create a mock controller, set channels 1,2 and 3 to 255 and render
the output to stdout once. If you have a real world controller that is supported
you can create a controller in the same way.

``` Go
config := usbdmx.NewConfig(vid, pid, outputInterfaceID, debugLevel)
config.GetUSBContext()

controller := mock.NewDMXController(config)
if err := controller.Connect(); err != nil {
    log.Fatal(err)
}

// set a fixture to bright white
controller.SetChannel(1, 255)
controller.SetChannel(2, 255)
controller.SetChannel(3, 255)

if err := controller.Render(); err != nil {
    log.Fatal(err)
}
```

## Contribute

Documentation is a big part of any project, but writing it is very time consuming.
Any documentation is greatly appreciated, please adhere to the Go standards for
writing documentation. Translations of existing documentation to other languages
will also be required.

Contributing code should have tests associated with it and all tests should pass.
Code will not be merged into the master branch until it has been proven to work
and has passing tests.

## Contributors

Oliver Read (Twitter: @oli_read)

## Support

If you would like to donate hardware please contact me either via Twitter.
Alternatively you can email me, my email can be found in my profile.

## License

This project is released under a GNU GPLv3 license, more information can be
found at: *[GNU GPLv3 Website](https://www.gnu.org/licenses/gpl-3.0.en.html)*

**To summarise**: You may copy, distribute and modify the software as long as you
track changes/dates in source files. Any modifications to or software including (via compiler)
GPL-licensed code must also be made available under the GPL along with build & install
instructions. If you wish to discuss any of these you can email me, which can be found
in my profile.
