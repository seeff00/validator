package validator

import (
	"reflect"
	"regexp"
	"strings"
)

type Validator func(s string) bool

type Struct struct {
	Key     string
	Value   string
	Tag     string
	IsValid bool
}

type Error struct {
	Key     string
	Value   string
	Tag     string
	IsValid bool
}

var TagMap = map[string]Validator{
	"comma_separated_numbers": IsCommaSeparatedNumber,
	"number":                  IsNumber,
	"alphabet":                IsAlphabet,
	"sort_format":             IsSortFormat,
}

func RegisterNewValidation(tag string, fn func(string) bool) {
	TagMap[tag] = fn
}

func ValidateStruct(obj interface{}) (bool, Error) {
	structElements := GetStructElements(obj, "validation")
	for i, structElement := range structElements {
		fn, isExists := TagMap[structElement.Tag]
		if isExists {
			structElements[i].IsValid = fn(structElement.Value)
			if !structElements[i].IsValid {
				return false, Error(structElements[i])
			}
		}
	}

	return true, Error{}
}

func GetStructElements(obj interface{}, tag string) []Struct {
	var result []Struct

	keyValues := reflect.ValueOf(obj)
	tags := reflect.TypeOf(obj)
	for i := 0; i < tags.NumField(); i++ {
		var s Struct
		if keyValues.Field(i).String() != "" && !keyValues.Field(i).IsZero() {
			s.Key = keyValues.Type().Field(i).Name
			s.Value = keyValues.Field(i).String()
		}

		if strings.TrimSpace(tags.Field(i).Tag.Get(tag)) != "" {
			s.Tag = tags.Field(i).Tag.Get(tag)
		}

		if strings.TrimSpace(s.Key) != "" && strings.TrimSpace(s.Value) != "" && strings.TrimSpace(s.Tag) != "" {
			result = append(result, s)
		}
	}

	return result
}

func IsCommaSeparatedNumber(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}

	for _, numAsStr := range strings.Split(s, ",") {
		if !IsNumber(numAsStr) {
			return false
		}
	}

	return true
}

func IsNumber(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}

	return regexp.MustCompile(`^[+\-]?(?:(?:0|[1-9]\d*)(?:\.\d*)?|\.\d+)$`).MatchString(s)
}

func IsAlphabet(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}

	return regexp.MustCompile(`^[A-Za-z]+$`).MatchString(s)
}

func IsSortFormat(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}

	var hasValidCharacters = regexp.MustCompile(`^[a-zA-Z\s_]+$`).MatchString(s)
	if !hasValidCharacters {
		return false
	}

	sortElements := strings.Split(s, " ")
	if len(sortElements) == 2 {
		if strings.ToLower(sortElements[1]) != "asc" && strings.ToLower(sortElements[1]) != "desc" {
			return false
		}
	}

	if len(sortElements) > 2 {
		return false
	}

	return true
}
