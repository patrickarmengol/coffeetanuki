package validator

import (
	"net/url"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var (
	LocationRX = regexp.MustCompile(`^([A-Za-z\s.'-]+, )+[A-Za-z\s.'-]+$`)
	EmailRX    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func New() *Validator {
	return &Validator{
		NonFieldErrors: []string{},
		FieldErrors:    map[string]string{},
	}
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	if v.NonFieldErrors == nil {
		v.NonFieldErrors = []string{}
	}
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// validation check helpers

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxBytes(value string, n int) bool {
	return len(value) <= n
}

func MinBytes(value string, n int) bool {
	return len(value) >= n
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func IsURL(value string) bool {
	// TODO: probably be more strict validation to avoid xss
	u, err := url.Parse(value)
	if err != nil {
		return false // unable to parse
	} else if u.Scheme == "" || u.Host == "" {
		return false // relative
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return false // not website url
	} else {
		return true
	}
}
