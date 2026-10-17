package main

import (
	"flag"
	"image/png"
	"os"
	"testing"

	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/test"
)

var update = flag.Bool("update", false, "rewrite ../screenshots")

// go test ./example -update
func TestSnapshots(t *testing.T) {
	if !*update {
		t.Skip("run with -update to regenerate the README screenshots")
	}
	a := test.NewTempApp(t)
	for file, s := range map[string]state{
		"baseline-light":          {seed: "Baseline"},
		"baseline-dark":           {seed: "Baseline", dark: true},
		"teal-light":              {seed: "Teal"},
		"deep-blue-fidelity-dark": {seed: "Deep blue", dark: true, fidelity: true},
	} {
		a.Settings().SetTheme(s.theme()) // before demo(), the swatches read colors at build time
		c := software.NewCanvas()
		c.SetScale(2)
		c.SetContent(demo(s, func(state) {}))
		f, err := os.Create("../screenshots/" + file + ".png")
		if err != nil {
			t.Fatal(err)
		}
		if err := png.Encode(f, c.Capture()); err != nil {
			t.Fatal(err)
		}
		f.Close()
	}
}
