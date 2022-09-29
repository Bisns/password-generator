package gui

import (
	"password-generator/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func settingsWindow() fyne.CanvasObject {
	themeText := canvas.NewText("主题切换", nil)
	dropdown := widget.NewSelect([]string{"Light", "Dark"}, parseTheme())
	t := fyne.CurrentApp().Preferences().StringWithFallback("Theme", "Light")
	switch t {
	case "Light":
		dropdown.PlaceHolder = "Light"
	case "Dark":
		dropdown.PlaceHolder = "Dark"
	}

	dropdown.Refresh()

	settings := container.NewVBox(themeText, dropdown)
	return settings
}

func parseTheme() func(string) {
	return func(t string) {
		switch t {
		case "Light":
			fyne.CurrentApp().Preferences().SetString("Theme", "Light")
			fyne.CurrentApp().Settings().SetTheme(&theme.MyTheme{Theme: "Light"})
		case "Dark":
			fyne.CurrentApp().Preferences().SetString("Theme", "Dark")
			fyne.CurrentApp().Settings().SetTheme(&theme.MyTheme{Theme: "Dark"})
		}
	}
}
