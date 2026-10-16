package mdi3

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// M3 color roles, for theme.Color(mdi3.ColorNameSurfaceContainerHigh) etc.
// primary, error, background, success and warning use Fyne's own names.
const (
	ColorNameOnPrimary             fyne.ThemeColorName = "onPrimary"
	ColorNamePrimaryContainer      fyne.ThemeColorName = "primaryContainer"
	ColorNameOnPrimaryContainer    fyne.ThemeColorName = "onPrimaryContainer"
	ColorNameInversePrimary        fyne.ThemeColorName = "inversePrimary"
	ColorNamePrimaryFixed          fyne.ThemeColorName = "primaryFixed"
	ColorNamePrimaryFixedDim       fyne.ThemeColorName = "primaryFixedDim"
	ColorNameOnPrimaryFixed        fyne.ThemeColorName = "onPrimaryFixed"
	ColorNameOnPrimaryFixedVariant fyne.ThemeColorName = "onPrimaryFixedVariant"
	ColorNameSurfaceTint           fyne.ThemeColorName = "surfaceTint"

	ColorNameSecondary               fyne.ThemeColorName = "secondary"
	ColorNameOnSecondary             fyne.ThemeColorName = "onSecondary"
	ColorNameSecondaryContainer      fyne.ThemeColorName = "secondaryContainer"
	ColorNameOnSecondaryContainer    fyne.ThemeColorName = "onSecondaryContainer"
	ColorNameSecondaryFixed          fyne.ThemeColorName = "secondaryFixed"
	ColorNameSecondaryFixedDim       fyne.ThemeColorName = "secondaryFixedDim"
	ColorNameOnSecondaryFixed        fyne.ThemeColorName = "onSecondaryFixed"
	ColorNameOnSecondaryFixedVariant fyne.ThemeColorName = "onSecondaryFixedVariant"

	ColorNameTertiary               fyne.ThemeColorName = "tertiary"
	ColorNameOnTertiary             fyne.ThemeColorName = "onTertiary"
	ColorNameTertiaryContainer      fyne.ThemeColorName = "tertiaryContainer"
	ColorNameOnTertiaryContainer    fyne.ThemeColorName = "onTertiaryContainer"
	ColorNameTertiaryFixed          fyne.ThemeColorName = "tertiaryFixed"
	ColorNameTertiaryFixedDim       fyne.ThemeColorName = "tertiaryFixedDim"
	ColorNameOnTertiaryFixed        fyne.ThemeColorName = "onTertiaryFixed"
	ColorNameOnTertiaryFixedVariant fyne.ThemeColorName = "onTertiaryFixedVariant"

	ColorNameOnError          fyne.ThemeColorName = "onError"
	ColorNameErrorContainer   fyne.ThemeColorName = "errorContainer"
	ColorNameOnErrorContainer fyne.ThemeColorName = "onErrorContainer"

	ColorNameSurface                 fyne.ThemeColorName = "surface"
	ColorNameSurfaceDim              fyne.ThemeColorName = "surfaceDim"
	ColorNameSurfaceBright           fyne.ThemeColorName = "surfaceBright"
	ColorNameSurfaceContainerLowest  fyne.ThemeColorName = "surfaceContainerLowest"
	ColorNameSurfaceContainerLow     fyne.ThemeColorName = "surfaceContainerLow"
	ColorNameSurfaceContainer        fyne.ThemeColorName = "surfaceContainer"
	ColorNameSurfaceContainerHigh    fyne.ThemeColorName = "surfaceContainerHigh"
	ColorNameSurfaceContainerHighest fyne.ThemeColorName = "surfaceContainerHighest"
	ColorNameOnSurface               fyne.ThemeColorName = "onSurface"
	ColorNameSurfaceVariant          fyne.ThemeColorName = "surfaceVariant"
	ColorNameOnSurfaceVariant        fyne.ThemeColorName = "onSurfaceVariant"
	ColorNameInverseSurface          fyne.ThemeColorName = "inverseSurface"
	ColorNameInverseOnSurface        fyne.ThemeColorName = "inverseOnSurface"
	ColorNameOnBackground            fyne.ThemeColorName = "onBackground"
	ColorNameOutline                 fyne.ThemeColorName = "outline"
	ColorNameOutlineVariant          fyne.ThemeColorName = "outlineVariant"
	ColorNameScrim                   fyne.ThemeColorName = "scrim"

	// Custom color groups (M3 "extended colors"), harmonized with the seed.
	ColorNameSuccessContainer   fyne.ThemeColorName = "successContainer"
	ColorNameOnSuccessContainer fyne.ThemeColorName = "onSuccessContainer"
	ColorNameWarningContainer   fyne.ThemeColorName = "warningContainer"
	ColorNameOnWarningContainer fyne.ThemeColorName = "onWarningContainer"
)

const (
	palPrimary = iota
	palSecondary
	palTertiary
	palNeutral
	palNeutralVariant
	palError
	palSuccess
	palWarning
)

// roles are the M3 2021-spec standard contrast tones (with the 2024
// tone-30 on*Container update) as {palette, light tone, dark tone}.
// ponytail: fixed tones, port DynamicColor contrast curves if medium/high contrast schemes are needed.
var roles = map[fyne.ThemeColorName]struct {
	pal         int
	light, dark float64
}{
	theme.ColorNamePrimary:         {palPrimary, 40, 80},
	ColorNameOnPrimary:             {palPrimary, 100, 20},
	ColorNamePrimaryContainer:      {palPrimary, 90, 30},
	ColorNameOnPrimaryContainer:    {palPrimary, 30, 90},
	ColorNameInversePrimary:        {palPrimary, 80, 40},
	ColorNamePrimaryFixed:          {palPrimary, 90, 90},
	ColorNamePrimaryFixedDim:       {palPrimary, 80, 80},
	ColorNameOnPrimaryFixed:        {palPrimary, 10, 10},
	ColorNameOnPrimaryFixedVariant: {palPrimary, 30, 30},
	ColorNameSurfaceTint:           {palPrimary, 40, 80},

	ColorNameSecondary:               {palSecondary, 40, 80},
	ColorNameOnSecondary:             {palSecondary, 100, 20},
	ColorNameSecondaryContainer:      {palSecondary, 90, 30},
	ColorNameOnSecondaryContainer:    {palSecondary, 30, 90},
	ColorNameSecondaryFixed:          {palSecondary, 90, 90},
	ColorNameSecondaryFixedDim:       {palSecondary, 80, 80},
	ColorNameOnSecondaryFixed:        {palSecondary, 10, 10},
	ColorNameOnSecondaryFixedVariant: {palSecondary, 30, 30},

	ColorNameTertiary:               {palTertiary, 40, 80},
	ColorNameOnTertiary:             {palTertiary, 100, 20},
	ColorNameTertiaryContainer:      {palTertiary, 90, 30},
	ColorNameOnTertiaryContainer:    {palTertiary, 30, 90},
	ColorNameTertiaryFixed:          {palTertiary, 90, 90},
	ColorNameTertiaryFixedDim:       {palTertiary, 80, 80},
	ColorNameOnTertiaryFixed:        {palTertiary, 10, 10},
	ColorNameOnTertiaryFixedVariant: {palTertiary, 30, 30},

	theme.ColorNameError:      {palError, 40, 80},
	ColorNameOnError:          {palError, 100, 20},
	ColorNameErrorContainer:   {palError, 90, 30},
	ColorNameOnErrorContainer: {palError, 30, 90},

	ColorNameSurface:                 {palNeutral, 98, 6},
	ColorNameSurfaceDim:              {palNeutral, 87, 6},
	ColorNameSurfaceBright:           {palNeutral, 98, 24},
	ColorNameSurfaceContainerLowest:  {palNeutral, 100, 4},
	ColorNameSurfaceContainerLow:     {palNeutral, 96, 10},
	ColorNameSurfaceContainer:        {palNeutral, 94, 12},
	ColorNameSurfaceContainerHigh:    {palNeutral, 92, 17},
	ColorNameSurfaceContainerHighest: {palNeutral, 90, 22},
	ColorNameOnSurface:               {palNeutral, 10, 90},
	ColorNameInverseSurface:          {palNeutral, 20, 90},
	ColorNameInverseOnSurface:        {palNeutral, 95, 20},
	ColorNameScrim:                   {palNeutral, 0, 0},
	ColorNameSurfaceVariant:          {palNeutralVariant, 90, 30},
	ColorNameOnSurfaceVariant:        {palNeutralVariant, 30, 80},
	ColorNameOutline:                 {palNeutralVariant, 50, 60},
	ColorNameOutlineVariant:          {palNeutralVariant, 80, 30},

	theme.ColorNameSuccess:             {palSuccess, 40, 80},
	theme.ColorNameForegroundOnSuccess: {palSuccess, 100, 20},
	ColorNameSuccessContainer:          {palSuccess, 90, 30},
	ColorNameOnSuccessContainer:        {palSuccess, 30, 90},
	theme.ColorNameWarning:             {palWarning, 40, 80},
	theme.ColorNameForegroundOnWarning: {palWarning, 100, 20},
	ColorNameWarningContainer:          {palWarning, 90, 30},
	ColorNameOnWarningContainer:        {palWarning, 30, 90},
}

// aliases map Fyne's widget color names onto M3 roles.
var aliases = map[fyne.ThemeColorName]fyne.ThemeColorName{
	theme.ColorNameBackground:                ColorNameSurface,
	ColorNameOnBackground:                    ColorNameOnSurface,
	theme.ColorNameForeground:                ColorNameOnSurface,
	theme.ColorNameForegroundOnPrimary:       ColorNameOnPrimary,
	theme.ColorNameForegroundOnError:         ColorNameOnError,
	theme.ColorNameHyperlink:                 theme.ColorNamePrimary,
	theme.ColorNameButton:                    ColorNameSecondaryContainer, // filled tonal button
	theme.ColorNameSelection:                 ColorNameSecondaryContainer, // selected list item
	theme.ColorNameHeaderBackground:          ColorNameSurfaceContainer,
	theme.ColorNameInputBackground:           ColorNameSurfaceContainerHighest, // filled text field, slider track
	theme.ColorNameInputBorder:               ColorNameOutline,
	theme.ColorNameMenuBackground:            ColorNameSurfaceContainer,
	theme.ColorNameOverlayBackground:         ColorNameSurfaceContainerHigh, // dialogs, popups
	theme.ColorNameInnerWindowBorder:         ColorNameSurfaceContainerHighest,
	theme.ColorNameInnerWindowBorderInactive: ColorNameSurfaceContainer,
	theme.ColorNamePlaceHolder:               ColorNameOnSurfaceVariant,
	theme.ColorNameScrollBar:                 ColorNameOutline,
	theme.ColorNameScrollBarBackground:       ColorNameSurfaceContainerHigh,
	theme.ColorNameSeparator:                 ColorNameOutlineVariant,
}

// stateLayers are M3 state layer / disabled opacities over a role.
var stateLayers = map[fyne.ThemeColorName]struct {
	role  fyne.ThemeColorName
	alpha uint8
}{
	theme.ColorNameHover:          {ColorNameOnSurface, 0x14},     // 8%
	theme.ColorNamePressed:        {ColorNameOnSurface, 0x1a},     // 10%
	theme.ColorNameFocus:          {theme.ColorNamePrimary, 0x29}, // 16%, M3 says 10% + a focus ring Fyne can't draw
	theme.ColorNameDisabledButton: {ColorNameOnSurface, 0x1f},     // 12%
	theme.ColorNameDisabled:       {ColorNameOnSurface, 0x61},     // 38%
}

var (
	baseline    = color.NRGBA{0x67, 0x50, 0xa4, 0xff} // M3 baseline seed
	successBase = color.NRGBA{0x4c, 0xaf, 0x50, 0xff}
	warningBase = color.NRGBA{0xff, 0x98, 0x00, 0xff}
)

type tonal struct{ h, c float64 }

func customPalette(base color.Color, seed hct) tonal {
	h := hctOf(harmonize(hctOf(base), seed).color())
	return tonal{h.h, max(48, h.c)}
}

// scheme builds an M3 "tonal spot" scheme (the Android/Theme Builder default) from seed.
func scheme(seed color.Color, dark bool) map[fyne.ThemeColorName]color.Color {
	s := hctOf(seed)
	pals := [...]tonal{
		palPrimary:        {s.h, 36},
		palSecondary:      {s.h, 16},
		palTertiary:       {sanitizeDeg(s.h + 60), 24},
		palNeutral:        {s.h, 6},
		palNeutralVariant: {s.h, 8},
		palError:          {25, 84},
		palSuccess:        customPalette(successBase, s),
		palWarning:        customPalette(warningBase, s),
	}

	m := make(map[fyne.ThemeColorName]color.Color, len(roles)+len(aliases)+len(stateLayers)+1)
	for n, r := range roles {
		t := r.light
		if dark {
			t = r.dark
		}
		m[n] = hct{pals[r.pal].h, pals[r.pal].c, t}.color()
	}
	for n, role := range aliases {
		m[n] = m[role]
	}
	for n, l := range stateLayers {
		c := m[l.role].(color.NRGBA)
		c.A = l.alpha
		m[n] = c
	}
	m[theme.ColorNameShadow] = color.NRGBA{A: 0x44} // Fyne's shadow alphas, tuned for its shadow renderer
	if dark {
		m[theme.ColorNameShadow] = color.NRGBA{A: 0x77}
	}
	return m
}
