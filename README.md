# Brief
Usage example for [goreleaser-cross](https://github.com/goreleaser/goreleaser-cross)

## Whats inside
Demonstrates how to cross-compile Go project with CGO dependencies.
Project depends on C/C++ libraries `libftdi1, libusb-1.0, opencv4`.
- Executable will try to detect FTDI 2232 boards with VID: 0x0403 and PID: 0x6010 and open channel 0.
- Then it lists serial number of all detected boards. 
- Afterwards it will try to open camera with OpenCV and take test picture into your terminal

Example of output
```
version: 0.0.4
detecting ftdi devices
detected 1 device(s) with vid: 0x0403 pid: 0x6010:
    [0] serial: FT4AYSJC
detecting camera using opencv...
show your awesome face and smile 😺
capturing in 3..2..1
<your picture is here>
you've been captured!!!!
```
Cross-compiles for targets:
 - Darwin/amd64
 - Linux/armhf (Raspberry Pi 4 aka RPI4)
 - Linux/amd64 (should not be any issue to add this example)

This example is based on real working solution where Macbook is dev machine and RPI4 is end target. Cross-compilation is set for both darwin and linux to ensure nothing is broken with new commits
 
## Sysroot
Sysroots are located in [separate repo](https://github.com/goreleaser/goreleaser-cross-example-sysroot) and added as submodule to current repo 
- **Darwin** sysroot is created by simply copying libraries from `homebrew`
- **Linux** sysroot is created from RPI4 using [this script](https://github.com/goreleaser/goreleaser-cross/blob/master/scripts/sysroot-rsync.sh)
