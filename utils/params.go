package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var BoolParseError = errors.New("bad value for bool field") // placeholder not passed to user

type Parameters map[string]string
type ParameterType int

const (
	ParameterTypeInt = iota
	ParameterTypeBool
	ParameterTypeDuration
	ParameterTypeDurationFalse
	ParameterTypeString
	ParameterTypeMultiString
)

// ParameterHint contains info on how a parameter
type ParameterHint struct {
	Key          string
	Description  string
	Type         ParameterType
	DefaultValue string
}

// Todo: was this necessary?
func TryGetParamsSectionIndex(indexedSchemePart string) (string, int, bool) {
	if strings.HasSuffix(indexedSchemePart, "]") {
		openingIndex := strings.LastIndex(indexedSchemePart, "[")
		if openingIndex != -1 {
			// now get what's in between
			between := indexedSchemePart[openingIndex : len(indexedSchemePart)-1]
			value, err := strconv.Atoi(between)
			if err != nil {
				return indexedSchemePart, 0, false
			}
			return indexedSchemePart[:openingIndex], value, true
		}
	}

	// it seems that it wasn't indexed!
	return indexedSchemePart, 0, false
}

func GetIndexedParamsSection(schemePart string, paramsSectionIndex int) string {
	return fmt.Sprintf("%s[%d]", schemePart, paramsSectionIndex)
}

func CommonParametersFromURL(q url.Values) Parameters {
	p := make(map[string]string)
	for key, values := range q {
		if !strings.ContainsAny(key, ".") {
			if values != nil && len(values) > 0 {
				p[key] = values[0]
			}
		}
	}
	return p
}

func CommonParametersFromMap(q map[string]string) Parameters {
	p := make(map[string]string)
	for key, value := range q {
		if !strings.ContainsAny(key, ".") {
			if value != "" {
				p[key] = value
			}
		}
	}
	return p
}

func ParametersBySectionFromURL(q url.Values, section string) Parameters {
	p := make(map[string]string)
	prefix := strings.TrimRight(section, ".") + "."
	for key, values := range q {
		if strings.HasPrefix(key, prefix) {
			if values != nil && len(values) > 0 {
				p[key[len(prefix):]] = values[0]
			}
		}
	}
	return p
}

func ParametersBySectionFromMap(q map[string]string, section string) Parameters {
	p := make(map[string]string)
	prefix := strings.TrimRight(section, ".") + "."
	for key, value := range q {
		if strings.HasPrefix(key, prefix) {
			if value != "" {
				p[key[len(prefix):]] = value
			}
		}
	}
	return p
}

func ParametersFromURL(q url.Values) Parameters {
	p := make(map[string]string)
	for key, values := range q {
		if values != nil && len(values) > 0 && values[0] != "" {
			p[key] = values[0]
		}
	}
	return p
}

func ParametersFromMap(q map[string]string) Parameters {
	p := make(map[string]string)
	for key, value := range q {
		if value != "" {
			p[key] = value
		}
	}
	return p
}

func (p Parameters) Has(key string) bool {
	_, found := p[key]
	return found
}

func (p Parameters) Get(key string) (string, bool) {
	v, found := p[key]
	return v, found
}

func (p Parameters) Set(key string) (string, bool) {
	v, found := p[key]
	return v, found
}

func (p Parameters) Remove(key string) {
	delete(p, key)
}

func (p Parameters) Section(section string) Parameters {
	newP := make(map[string]string)
	prefix := strings.TrimRight(section, ".") + "."
	for key, value := range p {
		if strings.HasPrefix(key, prefix) {
			if value != "" {
				newP[key[len(prefix):]] = value
			}
		}
	}
	return newP
}

func (p Parameters) SectionWithCommon(section string) Parameters {
	newP := make(map[string]string)
	prefix := strings.TrimRight(section, ".") + "."
	for key, value := range p {
		if strings.HasPrefix(key, prefix) || !strings.ContainsAny(key, ".") {
			if value != "" {
				newP[key[len(prefix):]] = value
			}
		}
	}
	return newP
}

func CombineParameters(ps ...Parameters) Parameters {
	parameters := make(map[string]string)
	for _, p := range ps {
		for key, value := range p {
			parameters[key] = value
		}
	}
	return parameters
}

func DurationFromParameters(params Parameters, key string, defaultValue time.Duration) time.Duration {
	if value, found := params.Get(key); found {
		parsed, err := time.ParseDuration(value)
		if err == nil {
			return parsed
		}
	}
	return defaultValue
}

func MultiStringFromParameters(params Parameters, key string, defaultValue []string) []string {
	if value, found := params.Get(key); found {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func StringFromParameters(params Parameters, key string, defaultValue string) string {
	if value, found := params.Get(key); found {
		return value
	}
	return defaultValue
}

func IntegerFromParameters(params Parameters, key string, defaultValue int) int {
	if value, found := params.Get(key); found {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}
	return defaultValue
}

func BoolFromParameters(params Parameters, key string, defaultValue bool) bool {
	if value, found := params.Get(key); found {
		parsed, err := ParseBool(value)
		if err == nil {
			return parsed
		}
	}
	return defaultValue
}

func ParseBool(value string) (bool, error) {
	if StrIsTrue(value) {
		return true, nil
	}
	if StrIsFalse(value) {
		return false, nil
	}
	return false, BoolParseError
}
func StrIsTrue(value string) bool {
	value = strings.ToLower(value)
	if value == "1" || value == "y" || value == "yes" || value == "true" || value == "t" {
		return true
	}
	return false
}
func StrIsFalse(value string) bool {
	value = strings.ToLower(value)
	if value == "0" || value == "n" || value == "no" || value == "false" || value == "f" {
		return true
	}
	return false
}
