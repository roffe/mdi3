// Package mdi3 is a Material Design 3 theme for Fyne: tonal spot light and
// dark schemes generated from a seed color (HCT), every M3 color role, Roboto
// typography with the M3 type scale, the M3 shape scale and Material Symbols
// icons.
//
//	a := app.New()
//	a.Settings().SetTheme(mdi3.New(nil)) // or mdi3.New(myBrandColor)
package mdi3

//go:generate go run gen.go

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

// M3 type scale, e.g. canvas.Text.TextSize = theme.Size(mdi3.SizeNameTitleLarge).
const (
	SizeNameDisplayLarge   fyne.ThemeSizeName = "displayLarge"
	SizeNameDisplayMedium  fyne.ThemeSizeName = "displayMedium"
	SizeNameDisplaySmall   fyne.ThemeSizeName = "displaySmall"
	SizeNameHeadlineLarge  fyne.ThemeSizeName = "headlineLarge"
	SizeNameHeadlineMedium fyne.ThemeSizeName = "headlineMedium"
	SizeNameHeadlineSmall  fyne.ThemeSizeName = "headlineSmall"
	SizeNameTitleLarge     fyne.ThemeSizeName = "titleLarge"
	SizeNameTitleMedium    fyne.ThemeSizeName = "titleMedium"
	SizeNameTitleSmall     fyne.ThemeSizeName = "titleSmall"
	SizeNameBodyLarge      fyne.ThemeSizeName = "bodyLarge"
	SizeNameBodyMedium     fyne.ThemeSizeName = "bodyMedium"
	SizeNameBodySmall      fyne.ThemeSizeName = "bodySmall"
	SizeNameLabelLarge     fyne.ThemeSizeName = "labelLarge"
	SizeNameLabelMedium    fyne.ThemeSizeName = "labelMedium"
	SizeNameLabelSmall     fyne.ThemeSizeName = "labelSmall"
)

// M3 shape scale (corner radii), e.g. rect.CornerRadius = theme.Size(mdi3.SizeNameCornerLarge).
const (
	SizeNameCornerExtraSmall          fyne.ThemeSizeName = "cornerExtraSmall"
	SizeNameCornerSmall               fyne.ThemeSizeName = "cornerSmall"
	SizeNameCornerMedium              fyne.ThemeSizeName = "cornerMedium"
	SizeNameCornerLarge               fyne.ThemeSizeName = "cornerLarge"
	SizeNameCornerLargeIncreased      fyne.ThemeSizeName = "cornerLargeIncreased"
	SizeNameCornerExtraLarge          fyne.ThemeSizeName = "cornerExtraLarge"
	SizeNameCornerExtraLargeIncreased fyne.ThemeSizeName = "cornerExtraLargeIncreased"
	SizeNameCornerExtraExtraLarge     fyne.ThemeSizeName = "cornerExtraExtraLarge"
	SizeNameCornerFull                fyne.ThemeSizeName = "cornerFull"
)

var sizes = map[fyne.ThemeSizeName]float32{
	SizeNameDisplayLarge:   57,
	SizeNameDisplayMedium:  45,
	SizeNameDisplaySmall:   36,
	SizeNameHeadlineLarge:  32,
	SizeNameHeadlineMedium: 28,
	SizeNameHeadlineSmall:  24,
	SizeNameTitleLarge:     22,
	SizeNameTitleMedium:    16,
	SizeNameTitleSmall:     14,
	SizeNameBodyLarge:      16,
	SizeNameBodyMedium:     14,
	SizeNameBodySmall:      12,
	SizeNameLabelLarge:     14,
	SizeNameLabelMedium:    12,
	SizeNameLabelSmall:     11,

	SizeNameCornerExtraSmall:          4,
	SizeNameCornerSmall:               8,
	SizeNameCornerMedium:              12,
	SizeNameCornerLarge:               16,
	SizeNameCornerLargeIncreased:      20,
	SizeNameCornerExtraLarge:          28,
	SizeNameCornerExtraLargeIncreased: 32,
	SizeNameCornerExtraExtraLarge:     48,
	SizeNameCornerFull:                canvas.RadiusMaximum,

	theme.SizeNameText:              14,                   // body medium
	theme.SizeNameHeadingText:       24,                   // headline small
	theme.SizeNameSubHeadingText:    16,                   // title medium
	theme.SizeNameCaptionText:       12,                   // body small, text field supporting text
	theme.SizeNameButtonRadius:      canvas.RadiusMaximum, // full
	theme.SizeNameInputRadius:       4,                    // extra small, text fields
	theme.SizeNameSelectionRadius:   4,                    // extra small
	theme.SizeNameMenuRadius:        4,                    // extra small
	theme.SizeNamePopupRadius:       4,                    // extra small, menus
	theme.SizeNameCardRadius:        12,                   // medium
	theme.SizeNameInnerWindowRadius: 12,                   // medium
	theme.SizeNameDialogRadius:      28,                   // extra large
	theme.SizeNameScrollBarRadius:   canvas.RadiusMaximum, // full
}

var (
	//go:embed fonts/Roboto-Regular.ttf
	robotoRegular []byte
	//go:embed fonts/Roboto-Medium.ttf
	robotoMedium []byte
	//go:embed fonts/Roboto-Italic.ttf
	robotoItalic []byte
	//go:embed fonts/Roboto-MediumItalic.ttf
	robotoMediumItalic []byte

	regular      = fyne.NewStaticResource("Roboto-Regular.ttf", robotoRegular)
	medium       = fyne.NewStaticResource("Roboto-Medium.ttf", robotoMedium)
	italic       = fyne.NewStaticResource("Roboto-Italic.ttf", robotoItalic)
	mediumItalic = fyne.NewStaticResource("Roboto-MediumItalic.ttf", robotoMediumItalic)
)

type m3 struct {
	light, dark map[fyne.ThemeColorName]color.Color
}

// New returns an M3 "tonal spot" theme with light and dark schemes generated
// from seed, following the user's light/dark preference. Only the seed's hue
// is used, so any seed gives the same calm M3 look. A nil seed gives the M3
// baseline (#6750A4).
func New(seed color.Color) fyne.Theme { return newTheme(seed, false) }

// NewFidelity returns an M3 "fidelity" theme: it keeps the seed's chroma and
// uses the seed itself as the primary container, so a saturated brand color
// stays saturated. Tertiary is the seed's temperature complement.
func NewFidelity(seed color.Color) fyne.Theme { return newTheme(seed, true) }

func newTheme(seed color.Color, fidelity bool) fyne.Theme {
	if seed == nil {
		seed = baseline
	}
	s := hctOf(seed)
	p := palettes(s, fidelity)
	return &m3{light: scheme(s, p, fidelity, false), dark: scheme(s, p, fidelity, true)}
}

func (t *m3) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	s := t.light
	if v == theme.VariantDark {
		s = t.dark
	}
	if c, ok := s[n]; ok {
		return c
	}
	return theme.DefaultTheme().Color(n, v)
}

// Font maps bold to Roboto Medium, M3's emphasis weight (titles, labels, buttons).
func (*m3) Font(s fyne.TextStyle) fyne.Resource {
	switch {
	case s.Monospace, s.Symbol:
		return theme.DefaultTheme().Font(s)
	case s.Bold && s.Italic:
		return mediumItalic
	case s.Bold:
		return medium
	case s.Italic:
		return italic
	}
	return regular
}

func (*m3) Size(n fyne.ThemeSizeName) float32 {
	if s, ok := sizes[n]; ok {
		return s
	}
	return theme.DefaultTheme().Size(n)
}
