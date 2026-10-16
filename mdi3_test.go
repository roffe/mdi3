package mdi3

import (
	"image/color"
	"math"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func TestHCT(t *testing.T) {
	// CAM16 reference values from material-color-utilities' tests.
	for _, tc := range []struct {
		c           color.NRGBA
		hue, chroma float64
	}{
		{color.NRGBA{0xff, 0, 0, 0xff}, 27.408, 113.357},
		{color.NRGBA{0, 0xff, 0, 0xff}, 142.139, 108.410},
		{color.NRGBA{0, 0, 0xff, 0xff}, 282.788, 87.230},
		{color.NRGBA{0xff, 0xff, 0xff, 0xff}, 209.492, 2.869},
	} {
		h := hctOf(tc.c)
		if math.Abs(h.h-tc.hue) > 0.001 || math.Abs(h.c-tc.chroma) > 0.001 {
			t.Errorf("%v: got hue %.3f chroma %.3f, want %.3f %.3f", tc.c, h.h, h.c, tc.hue, tc.chroma)
		}
	}

	// Every in-gamut color must survive sRGB -> HCT -> sRGB.
	for r := 0; r < 256; r += 15 {
		for g := 0; g < 256; g += 15 {
			for b := 0; b < 256; b += 15 {
				c := color.NRGBA{uint8(r), uint8(g), uint8(b), 0xff}
				if got := hctOf(c).color(); got != c {
					t.Fatalf("round trip %v -> %v", c, got)
				}
			}
		}
	}

	// Out of gamut requests clip chroma but keep tone.
	got := hctOf(hct{120, 200, 50}.color())
	if math.Abs(got.t-50) > 0.5 || got.c > 200 {
		t.Errorf("gamut map: %+v", got)
	}
}

func TestTheme(t *testing.T) {
	th := New(nil)
	if c := th.Color(ColorNameSurfaceContainerLowest, theme.VariantLight); c != (color.NRGBA{0xff, 0xff, 0xff, 0xff}) {
		t.Errorf("light surfaceContainerLowest = %v, want white", c)
	}
	for _, tc := range []struct {
		n    fyne.ThemeColorName
		v    fyne.ThemeVariant
		tone float64
	}{
		{theme.ColorNamePrimary, theme.VariantLight, 40},
		{theme.ColorNamePrimary, theme.VariantDark, 80},
		{theme.ColorNameBackground, theme.VariantDark, 6},
		{theme.ColorNameButton, theme.VariantLight, 90},
	} {
		if got := hctOf(th.Color(tc.n, tc.v)).t; math.Abs(got-tc.tone) > 0.5 {
			t.Errorf("%s variant %d: tone %.2f, want %v", tc.n, tc.v, got, tc.tone)
		}
	}
	if got := th.Icon(theme.IconNameWindowClose).Name(); got != "foreground_m3_close.svg" {
		t.Errorf("WindowClose icon = %s, want the Material Symbol", got)
	}
	if got, want := th.Icon(theme.IconNameDragCornerIndicator).Name(), theme.DefaultTheme().Icon(theme.IconNameDragCornerIndicator).Name(); got != want {
		t.Errorf("unmapped icon = %s, want Fyne fallback %s", got, want)
	}
}
