//       ___  _____  ____
//      / _ \/  _/ |/_/ /____ ______ _
//     / ___// /_>  </ __/ -_) __/  ' \
//    /_/  /___/_/|_|\__/\__/_/ /_/_/_/
//
//    Copyright 2017 Eliuk Blau
//
//    This Source Code Form is subject to the terms of the Mozilla Public
//    License, v. 2.0. If a copy of the MPL was not distributed with this
//    file, You can obtain one at https://mozilla.org/MPL/2.0/.

package render

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/crypto/ssh/terminal"
)

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

func Render(r io.Reader) error {
	var (
		pix *ansimage.ANSImage
		err error
	)

	// get terminal size
	tx, ty, err := getTerminalSize()
	if err != nil {
		return err
	}

	// get scale mode from flag
	sm := ansimage.ScaleMode(0)

	// get dithering mode from flag
	dm := ansimage.DitheringMode(0)

	// set image scale factor for ANSIPixel grid
	sfy, sfx := 2, 1

	mc, err := colorful.Hex("#000000") // RGB color from Hex format
	if err != nil {
		return err
	}

	pix, err = ansimage.NewScaledFromReader(r, sfy*ty, sfx*tx, mc, sm, dm)
	if err != nil {
		return err
	}

	// draw ANSImage to terminal
	if isTerminal() {
		ansimage.ClearTerminal()
	}
	pix.SetMaxProcs(runtime.NumCPU()) // maximum number of parallel goroutines!
	pix.DrawExt(false, false)
	if isTerminal() {
		fmt.Println()
	}
	return nil
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // use paralelism for goroutines!
}
