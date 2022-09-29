package gui

import (
	"password-generator/generator"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

type Password struct {
	Numbers           bool
	Lowercase         bool
	Uppercase         bool
	Symbol            bool
	AllSymbol         bool
	SimilarCharacters bool
	Duplicate         bool
	Length            uint
}

func mainWindow(s *NewScreen) fyne.CanvasObject {
	w := s.Current

	a := s.APP

	number := a.Preferences().BoolWithFallback("Number", false)
	lowercase := a.Preferences().BoolWithFallback("Lowercase", false)
	uppercase := a.Preferences().BoolWithFallback("Uppercase", false)
	symbol := a.Preferences().BoolWithFallback("Symbol", false)
	allsymbol := a.Preferences().BoolWithFallback("AllSymbol", false)
	similarcharacters := a.Preferences().BoolWithFallback("SimilarCharacters", false)
	duplicate := a.Preferences().BoolWithFallback("Duplicate", false)

	length := a.Preferences().IntWithFallback("Length", 0)
	lengthBind := binding.NewFloat()
	_ = lengthBind.Set(float64(length))

	pwdOptionsBind := binding.NewBoolList()
	_ = pwdOptionsBind.Set([]bool{number, lowercase, uppercase, symbol, allsymbol, similarcharacters, duplicate})

	slide := widget.NewSliderWithData(0, 64, lengthBind)
	slide.Step = 1
	lengthText := widget.NewLabelWithData(binding.FloatToStringWithFormat(lengthBind, "密码长度：%0.0f"))

	buttons := container.NewGridWithColumns(4,
		widget.NewButton("8", func() {
			_ = lengthBind.Set(8)
		}),
		widget.NewButton("16", func() {
			_ = lengthBind.Set(16)
		}),
		widget.NewButton("32", func() {
			_ = lengthBind.Set(32)
		}),
		widget.NewButton("64", func() {
			_ = lengthBind.Set(64)
		}))

	lengthLabel := container.NewGridWithColumns(2, container.New(layout.NewFormLayout(), lengthText, slide), buttons)

	password := widget.NewEntry()
	s.Password = password

	pwdEntropy := binding.NewString()
	pwdEntropyLabel := widget.NewLabelWithData(pwdEntropy)
	pwdEntropyText := canvas.NewText("密码强度：", nil)
	s.pwdEntropy = pwdEntropy

	slide.OnChanged = func(f float64) {
		_ = lengthBind.Set(f)
		a.Preferences().SetInt("Length", int(f))
	}

	lengthBind.AddListener(binding.NewDataListener(func() {
		pwdSetText(s)
	}))

	s.lengthBind = lengthBind
	s.pwdOptionsBind = pwdOptionsBind

	NumberCheck := widgetCheck(s, "数字", "Number", number)
	LowercaseCheck := widgetCheck(s, "小写字母", "Lowercase", lowercase)
	UppercaseCheck := widgetCheck(s, "大写字母", "Uppercase", uppercase)
	SymbolCheck := widgetCheck(s, "常见特殊字符", "Symbol", symbol)
	AllSymbolCheck := widgetCheck(s, "所有特殊字符", "AllSymbol", allsymbol)
	SimilarCharactersCheck := widgetCheck(s, "排除相似字符", "SimilarCharacters", similarcharacters)
	DuplicateCheck := widgetCheck(s, "字符不重复", "Duplicate", duplicate)

	copyBtn := widget.NewButtonWithIcon("复制", theme.ContentCopyIcon(), func() {
		w.Clipboard().SetContent(password.Text)
	})
	copyBtn.Importance = widget.HighImportance

	updateBth := widget.NewButtonWithIcon("刷新", theme.ViewRefreshIcon(), func() {
		pwdSetText(s)
	})

	resetBth := widget.NewButtonWithIcon("重置", theme.CancelIcon(), func() {
		password.SetText("")
		slide.SetValue(0)
		NumberCheck.SetChecked(false)
		LowercaseCheck.SetChecked(false)
		UppercaseCheck.SetChecked(false)
		SymbolCheck.SetChecked(false)
		AllSymbolCheck.SetChecked(false)
		SimilarCharactersCheck.SetChecked(false)
		DuplicateCheck.SetChecked(false)
	})

	opButtons := container.New(layout.NewGridLayout(3), copyBtn, updateBth, resetBth)

	checklists := container.New(layout.NewGridLayout(4), NumberCheck, LowercaseCheck, UppercaseCheck, SymbolCheck, AllSymbolCheck, SimilarCharactersCheck, DuplicateCheck)

	content := container.NewVBox(password, container.New(layout.NewFormLayout(), pwdEntropyText, pwdEntropyLabel), lengthLabel, checklists, opButtons)

	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyRight, fyne.KeyDown:
			if slide.Value < slide.Max {
				_ = lengthBind.Set(slide.Value + slide.Step)
			}
		case fyne.KeyLeft, fyne.KeyUp:
			if slide.Value > slide.Min {
				_ = lengthBind.Set(slide.Value - slide.Step)
			}
		case fyne.KeyF5:
			pwdSetText(s)
		}
	})

	return content
}

func widgetCheck(s *NewScreen, label string, key string, checked bool) *widget.Check {
	a := s.APP

	Check := widget.NewCheck(label, func(b bool) {})
	Check.SetChecked(checked)
	Check.OnChanged = func(b bool) {
		pwdOptionsBind := s.pwdOptionsBind

		a.Preferences().SetBool(key, b)
		number := a.Preferences().BoolWithFallback("Number", false)
		lowercase := a.Preferences().BoolWithFallback("Lowercase", false)
		uppercase := a.Preferences().BoolWithFallback("Uppercase", false)
		symbol := a.Preferences().BoolWithFallback("Symbol", false)
		allsymbol := a.Preferences().BoolWithFallback("AllSymbol", false)
		similarcharacters := a.Preferences().BoolWithFallback("SimilarCharacters", false)
		duplicate := a.Preferences().BoolWithFallback("Duplicate", false)
		_ = pwdOptionsBind.Set([]bool{number, lowercase, uppercase, symbol, allsymbol, similarcharacters, duplicate})

		pwdSetText(s)
	}
	return Check
}

func genPwd(p *Password) string {
	config := generator.Config{}
	config.IncludeNumbers = p.Numbers
	config.IncludeLowercaseLetters = p.Lowercase
	config.IncludeUppercaseLetters = p.Uppercase
	config.IncludeSymbols = p.Symbol
	config.IncludeAllSymbols = p.AllSymbol
	config.Duplicate = p.Duplicate

	if !p.Numbers && !p.Lowercase && !p.Uppercase && !p.Symbol && !p.AllSymbol {
		config.IncludeNumbers = true
	}

	config.ExcludeSimilarCharacters = p.SimilarCharacters

	config.Length = p.Length

	g, _ := generator.New(&config)
	pwd, _ := g.Generate()
	return pwd
}

func pwdLevel(pwd string) string {
	e := passwordvalidator.GetEntropy(pwd)
	switch {
	case e < 20:
		return "密码强度非常低"
	case e < 40:
		return "密码强度低"
	case e < 60:
		return "密码强度一般"
	case e < 80:
		return "密码强度高"
	case e >= 80:
		return "密码强度非常高"
	default:
		return ""
	}
}

func pwdSetText(s *NewScreen) {
	pwdOptionsBind := s.pwdOptionsBind
	lengthBind := s.lengthBind
	password := s.Password
	pwdEntropy := s.pwdEntropy

	pwdOptions, _ := pwdOptionsBind.Get()
	length, _ := lengthBind.Get()

	if length > 0 {
		if !pwdOptions[0] && !pwdOptions[1] && !pwdOptions[2] && !pwdOptions[3] && !pwdOptions[4] {
			password.SetText("")
			_ = pwdEntropy.Set("")
		} else {
			pwd := genPwd(&Password{pwdOptions[0], pwdOptions[1], pwdOptions[2], pwdOptions[3], pwdOptions[4], pwdOptions[5], pwdOptions[6], uint(length)})
			password.SetText(pwd)
			_ = pwdEntropy.Set(pwdLevel(pwd))
		}
	} else {
		password.SetText("")
		_ = pwdEntropy.Set("")
	}
}
