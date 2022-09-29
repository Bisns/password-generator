package gui

import (
	theme2 "password-generator/theme"

	"fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type NewScreen struct {
	Current        fyne.Window
	APP            fyne.App
	Password       *widget.Entry
	pwdEntropy     binding.String
	pwdOptionsBind binding.BoolList
	lengthBind     binding.Float
}

// Start .
func Start(s *NewScreen) {
	w := s.Current

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("生成密码", theme.DocumentIcon(), container.NewPadded(mainWindow(s))),
		container.NewTabItemWithIcon("主题设置", theme.SettingsIcon(), container.NewPadded(settingsWindow())),
	)

	tabs.OnSelected = func(t *container.TabItem) {
		t.Content.Refresh()
	}

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(w.Canvas().Size().Width, w.Canvas().Size().Height))
	w.CenterOnScreen()
	w.SetMaster()
	w.ShowAndRun()
}

// InitNewScreen .
func InitNewScreen() *NewScreen {
	a := app.NewWithID("password-generator")

	t := a.Preferences().StringWithFallback("Theme", "Light")
	a.Settings().SetTheme(&theme2.MyTheme{Theme: t})
	a.SetIcon(theme2.Ico)
	w := a.NewWindow("随机密码生成器")

	return &NewScreen{
		Current: w,
		APP:     a,
	}
}
