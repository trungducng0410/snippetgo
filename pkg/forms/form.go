package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		if strings.TrimSpace(val) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	val := f.Get(field)
	if val == "" {
		return
	}
	if utf8.RuneCountInString(val) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maxixmum is %d)", d))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	val := f.Get(field)
	if val == "" {
		return
	}
	for _, opt := range opts {
		if val == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d)", d))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
