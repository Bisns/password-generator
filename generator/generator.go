package generator

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
)

const (
	LengthStrong              uint = 24
	DefaultLetterSet               = "abcdefghijklmnopqrstuvwxyz"
	DefaultLetterAmbiguousSet      = "ijlo"
	DefaultNumberSet               = "0123456789"
	DefaultNumberAmbiguousSet      = "01"
	DefaultSymbolSet               = "~!@#$%^&*-=+"
	DefaultSymbolAmbiguousSet      = "`()_{}[]\\|:;\"'<>,.?/"
)

var (
	DefaultConfig = Config{
		Length:                   LengthStrong,
		IncludeNumbers:           true,
		IncludeLowercaseLetters:  true,
		IncludeUppercaseLetters:  true,
		IncludeSymbols:           true,
		IncludeAllSymbols:        true,
		ExcludeSimilarCharacters: true,
		Duplicate:                false,
	}

	ErrConfigIsEmpty = errors.New("config is empty")
)

// Generator is what generates the password
type Generator struct {
	*Config
}

// Config is the config struct to hold the settings about
// what type of password to generate
type Config struct {
	// Length is the length of password to generate
	Length uint

	// CharacterSet is the setting to manually set the
	// character set
	CharacterSet string

	// IncludeSymbols is the setting to include symbols in
	// the character set
	// i.e. !"£*
	IncludeSymbols bool

	// IncludeNumbers is the setting to include number in
	// the character set
	// i.e. 1234
	IncludeNumbers bool

	// IncludeLowercaseLetters is the setting to include
	// lowercase letters in the character set
	// i.e. abcde
	IncludeLowercaseLetters bool

	// IncludeUppercaseLetters is the setting to include
	// uppercase letters in the character set
	// i.e. ABCD
	IncludeUppercaseLetters bool

	// ExcludeSimilarCharacters is the setting to exclude
	// characters that look the same in the character set
	// i.e. i1jIo0
	ExcludeSimilarCharacters bool

	// ExcludeAmbiguousCharacters is the setting to exclude
	// characters that can be hard to remember or symbols
	// that are rarely used
	// i.e. <>{}[]()/|\`
	IncludeAllSymbols bool

	Duplicate bool
}

// New returns a new generator
func New(config *Config) (*Generator, error) {
	if config == nil {
		config = &DefaultConfig
	}

	if !config.IncludeSymbols &&
		!config.IncludeUppercaseLetters &&
		!config.IncludeLowercaseLetters &&
		!config.IncludeNumbers &&
		!config.IncludeAllSymbols &&
		config.CharacterSet == "" {
		return nil, ErrConfigIsEmpty
	}

	if config.Length == 0 {
		config.Length = LengthStrong
	}

	if config.CharacterSet == "" {
		config.CharacterSet = buildCharacterSet(config)
	}

	return &Generator{Config: config}, nil
}

func buildCharacterSet(config *Config) string {
	var characterSet string
	if config.IncludeLowercaseLetters {
		characterSet += DefaultLetterSet
		if config.ExcludeSimilarCharacters {
			characterSet = removeCharacters(characterSet, DefaultLetterAmbiguousSet)
		}
	}

	if config.IncludeUppercaseLetters {
		characterSet += strings.ToUpper(DefaultLetterSet)
		if config.ExcludeSimilarCharacters {
			characterSet = removeCharacters(characterSet, strings.ToUpper(DefaultLetterAmbiguousSet))
		}
	}

	if config.IncludeNumbers {
		characterSet += DefaultNumberSet
		if config.ExcludeSimilarCharacters {
			characterSet = removeCharacters(characterSet, DefaultNumberAmbiguousSet)
		}
	}

	if config.IncludeSymbols {
		characterSet += DefaultSymbolSet
	}

	if config.IncludeAllSymbols {
		if !config.IncludeSymbols {
			characterSet += DefaultSymbolSet
		}
		characterSet += DefaultSymbolAmbiguousSet
	}

	return characterSet
}

func removeCharacters(str, characters string) string {
	return strings.Map(func(r rune) rune {
		if !strings.ContainsRune(characters, r) {
			return r
		}
		return -1
	}, str)
}

// NewWithDefault returns a new generator with the default
// config
func NewWithDefault() (*Generator, error) {
	return New(&DefaultConfig)
}

// Generate generates one password with length set in the
// config
func (g *Generator) Generate() (string, error) {
	var generated string
	characterSet := strings.Split(g.Config.CharacterSet, "")
	max := big.NewInt(int64(len(characterSet)))

	if g.Config.Duplicate {
		var length uint
		if g.Config.Length > uint(len(characterSet)) {
			length = uint(len(characterSet))
		} else {
			length = g.Config.Length
		}

		chars := make(map[string]struct{}, 0)
		for uint(len(chars)) < length {
			val, err := rand.Int(rand.Reader, max)
			if err != nil {
				return "", err
			}
			c := characterSet[val.Int64()]
			chars[c] = struct{}{}
		}
		for k := range chars {
			generated += k
		}
	} else {
		for i := uint(0); i < g.Config.Length; i++ {
			val, err := rand.Int(rand.Reader, max)
			if err != nil {
				return "", err
			}
			generated += characterSet[val.Int64()]
		}
	}
	return generated, nil
}

// GenerateMany generates multiple passwords with length set
// in the config
func (g *Generator) GenerateMany(amount uint) ([]string, error) {
	var generated []string
	for i := uint(0); i < amount; i++ {
		str, err := g.Generate()
		if err != nil {
			return nil, err
		}

		generated = append(generated, str)
	}
	return generated, nil
}

// GenerateWithLength generate one password with set length
func (g *Generator) GenerateWithLength(length uint) (string, error) {
	var generated string
	characterSet := strings.Split(g.Config.CharacterSet, "")
	max := big.NewInt(int64(len(characterSet)))
	for i := uint(0); i < length; i++ {
		val, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		generated += characterSet[val.Int64()]
	}
	return generated, nil
}

// GenerateManyWithLength generates multiple passwords with set length
func (g *Generator) GenerateManyWithLength(amount, length uint) ([]string, error) {
	var generated []string
	for i := uint(0); i < amount; i++ {
		str, err := g.GenerateWithLength(length)
		if err != nil {
			return nil, err
		}
		generated = append(generated, str)
	}
	return generated, nil
}
