package main

/*
#include <stdlib.h>
#include <ftdi.h>
#include <libusb.h>

#cgo pkg-config: libftdi1 libusb-1.0 opencv4
#cgo linux LDFLAGS: -pthread
*/
import "C"

import (
	"bytes"
	"fmt"
	_ "image/png"
	"os"
	"runtime"
	"time"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/notifai/ftdi"
	"gocv.io/x/gocv"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	version string
)

func main() {
	fmt.Printf("version: %s\n", version)

	fmt.Println("detecting ftdi devices")

	var devs []*ftdi.Device

	l, err := ftdi.FindAll(0x0403, 0x6010)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		goto camera
	}

	if len(l) == 0 {
		fmt.Println("no usb devices detected with vid: 0x0403 pid: 0x6010")
	} else {
		buf := &bytes.Buffer{}
		_, _ = fmt.Fprintf(buf, "detected %d device(s) with vid: 0x0403 pid: 0x6010:\n", len(l))
		for i, u := range l {
			_, _ = fmt.Fprintf(buf, "\t[%d] serial: %s", i, u.Serial)
			var dev *ftdi.Device
			if dev, err = ftdi.OpenUSBDev(u, 0); err != nil {
				fmt.Printf("couldn't open channel 0 of device: %s\n", u.Serial)
			} else {
				devs = append(devs, dev)
			}
		}

		fmt.Println(buf.String())
	}

	defer func() {
		for _, d := range devs {
			_ = d.Close()
		}
	}()

camera:
	fmt.Println("detecting camera using opencv...")
	var cam *gocv.VideoCapture
	cam, err = gocv.OpenVideoCapture(0)
	if err != nil {
		return
	}
	defer func() {
		_ = cam.Close()
	}()

	fmt.Println("show your awesome face and smile 😺")

	fmt.Printf("capturing in ")

	wait:
	for i := 3; i > 0; i-- {
		select {
		case <-time.After(1 * time.Second):
			fmt.Printf("%d", i)
			if i != 1 {
				fmt.Printf("..")
			} else {
				break wait
			}
		}
	}

	fmt.Printf("\n")

	// prepare image matrix
	cvImg := gocv.NewMat()
	defer func() {
		_ = cvImg.Close()
	}()

	_ = cam.Read(&cvImg)

	var png []byte
	if png, err = gocv.IMEncode(".png", cvImg); err != nil {
		fmt.Println(err.Error())
		return
	}

	// get terminal size
	var tx int
	var ty int

	if tx, ty, err = getTerminalSize(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// get scale mode from flag
	sm := ansimage.ScaleMode(0)

	// get dithering mode from flag
	dm := ansimage.DitheringMode(0)

	// set image scale factor for ANSIPixel grid
	sfy, sfx := ansimage.BlockSizeY, ansimage.BlockSizeX // 8x4 --> with dithering
	if ansimage.DitheringMode(0) == ansimage.NoDithering {
		sfy, sfx = 2, 1 // 2x1 --> without dithering
	}

	var mc colorful.Color
	if mc, err = colorful.Hex("#000000"); err != nil { // RGB color from Hex format
		fmt.Println(err.Error())
		return
	}
	var pix *ansimage.ANSImage
	pix, err = ansimage.NewScaledFromReader(bytes.NewBuffer(png), sfy*ty, sfx*tx, mc, sm, dm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// draw ANSImage to terminal
	if isTerminal() {
		ansimage.ClearTerminal()
	}

	pix.SetMaxProcs(runtime.NumCPU()) // maximum number of parallel goroutines!
	pix.DrawExt(false, false)

	fmt.Println("you've been captured!!!")
}

func isTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func getTerminalSize() (width, height int, err error) {
	if isTerminal() {
		return terminal.GetSize(int(os.Stdout.Fd()))
	}

	// fallback when piping to a file!
	return 80, 24, nil // VT100 terminal size
}
