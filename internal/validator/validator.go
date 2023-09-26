package validator

import (
	"strings"
	"unicode/utf8"
)

// Validator contains a map of validation errors for the form fields.
type Validator struct {
	FieldErrors map[string]string
}

// Valid returns true if the FieldErrors map is empty.
func (v Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldsError adds an error message to the FieldErrors map so long as no entry already exists for a given key.
func (v Validator) AddFieldsError(key, message string) {
	// Need to init the map first if it isn't already initialized
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField adds an error message to the FieldErrors map only if a validation check is not ok.
func (v Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldsError(key, message)
	}
}

// NotBlank returns true if a value us not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if the value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt returns true if a value is in a list of permitted integers.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}