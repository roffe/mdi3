// Demo of the mdi3 theme. Pick a seed color and flip dark mode to see the
// whole scheme regenerate.
package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/roffe/mdi3"
)

var seeds = map[string]color.Color{
	"Baseline": nil,
	"Teal":     color.NRGBA{R: 0x00, G: 0x6a, B: 0x6a, A: 0xff},
	"Green":    color.NRGBA{R: 0x38, G: 0x6a, B: 0x20, A: 0xff},
	"Orange":   color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0xff},
	"Red":      color.NRGBA{R: 0xb3, G: 0x26, B: 0x1e, A: 0xff},
	"Blue":     color.NRGBA{R: 0x1e, G: 0x88, B: 0xe5, A: 0xff},
}

type state struct {
	seed string
	dark bool
}

func (s state) theme() fyne.Theme {
	v := theme.VariantLight
	if s.dark {
		v = theme.VariantDark
	}
	return forced{mdi3.New(seeds[s.seed]), v}
}

// forced pins the variant so the demo can flip light/dark whatever the OS says.
type forced struct {
	fyne.Theme
	v fyne.ThemeVariant
}

func (f forced) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(n, f.v)
}

func main() {
	a := app.NewWithID("com.github.roffe.mdi3.demo")
	w := a.NewWindow("mdi3 demo")
	var apply func(state)
	apply = func(s state) {
		a.Settings().SetTheme(s.theme())
		w.SetContent(demo(s, apply)) // rebuilt so the swatches pick up the new colors
	}
	apply(state{"Baseline", a.Settings().ThemeVariant() == theme.VariantDark})
	w.ShowAndRun()
}

func demo(s state, apply func(state)) fyne.CanvasObject {
	seed := widget.NewSelect([]string{"Baseline", "Teal", "Green", "Orange", "Red", "Blue"}, nil)
	seed.Selected = s.seed
	seed.OnChanged = func(name string) { apply(state{name, s.dark}) }
	dark := widget.NewCheck("Dark mode", nil)
	dark.Checked = s.dark
	dark.OnChanged = func(on bool) { apply(state{s.seed, on}) }

	placeholder := widget.NewEntry()
	placeholder.SetPlaceHolder("Placeholder")
	entry := widget.NewEntry()
	entry.SetText("Some text")
	pick := widget.NewSelect([]string{"Alpha", "Beta", "Gamma"}, nil)
	pick.Selected = "Alpha"

	high := widget.NewButtonWithIcon("Filled", theme.ConfirmIcon(), nil)
	high.Importance = widget.HighImportance
	low := widget.NewButton("Text", nil)
	low.Importance = widget.LowImportance
	disabled := widget.NewButton("Disabled", nil)
	disabled.Disable()
	danger := widget.NewButton("Danger", nil)
	danger.Importance = widget.DangerImportance
	success := widget.NewButton("Success", nil)
	success.Importance = widget.SuccessImportance
	warning := widget.NewButton("Warning", nil)
	warning.Importance = widget.WarningImportance

	checked := widget.NewCheck("Checked", nil)
	checked.Checked = true
	radio := widget.NewRadioGroup([]string{"One", "Two"}, nil)
	radio.Horizontal = true
	radio.Selected = "One"
	slider := widget.NewSlider(0, 100)
	slider.Value = 40
	progress := widget.NewProgressBar()
	progress.Value = 0.6

	icons := container.NewHBox()
	for _, n := range []fyne.ThemeIconName{
		theme.IconNameHome, theme.IconNameSettings, theme.IconNameSearch, theme.IconNameDelete,
		theme.IconNameFolder, theme.IconNameMailSend, theme.IconNameMediaPlay, theme.IconNameInfo,
		theme.IconNameWarning, theme.IconNameError, theme.IconNameAccount, theme.IconNameCalendar,
	} {
		icons.Add(widget.NewIcon(theme.Icon(n)))
	}

	left := container.NewVBox(
		container.NewHBox(seed, dark),
		widget.NewRichTextFromMarkdown("# Material Design 3\n## Fyne theme\nRoboto body text with *italic* and **medium** weight."),
		placeholder, entry, pick,
		container.NewHBox(widget.NewButton("Tonal", nil), high, low, disabled),
		container.NewHBox(danger, success, warning),
		container.NewHBox(checked, widget.NewCheck("Unchecked", nil)),
		radio, slider, progress, icons,
	)

	swatches := container.NewGridWithColumns(3,
		swatch("Primary", theme.ColorNamePrimary, mdi3.ColorNameOnPrimary),
		swatch("Secondary", mdi3.ColorNameSecondary, mdi3.ColorNameOnSecondary),
		swatch("Tertiary", mdi3.ColorNameTertiary, mdi3.ColorNameOnTertiary),
		swatch("Primary container", mdi3.ColorNamePrimaryContainer, mdi3.ColorNameOnPrimaryContainer),
		swatch("Secondary container", mdi3.ColorNameSecondaryContainer, mdi3.ColorNameOnSecondaryContainer),
		swatch("Tertiary container", mdi3.ColorNameTertiaryContainer, mdi3.ColorNameOnTertiaryContainer),
		swatch("Error", theme.ColorNameError, mdi3.ColorNameOnError),
		swatch("Success", theme.ColorNameSuccess, theme.ColorNameForegroundOnSuccess),
		swatch("Warning", theme.ColorNameWarning, theme.ColorNameForegroundOnWarning),
	)
	list := widget.NewList(
		func() int { return 4 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) { o.(*widget.Label).SetText("List item") },
	)
	list.Select(1)
	listHeight := canvas.NewRectangle(color.Transparent)
	listHeight.SetMinSize(fyne.NewSize(0, 150))

	right := container.NewVBox(
		container.NewAppTabs(
			container.NewTabItemWithIcon("Home", theme.HomeIcon(),
				widget.NewCard("Card title", "Subtitle", widget.NewLabel("Cards get the medium (12px) corners."))),
			container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), widget.NewLabel("")),
		),
		swatches,
		container.NewStack(listHeight, list),
	)
	return container.NewPadded(container.NewHBox(left, right))
}

func swatch(label string, bg, fg fyne.ThemeColorName) fyne.CanvasObject {
	r := canvas.NewRectangle(theme.Color(bg))
	r.CornerRadius = theme.Size(mdi3.SizeNameCornerSmall)
	t := canvas.NewText(label, theme.Color(fg))
	t.TextSize = theme.Size(mdi3.SizeNameLabelMedium)
	t.TextStyle.Bold = true
	return container.NewStack(r, container.NewPadded(t))
}
