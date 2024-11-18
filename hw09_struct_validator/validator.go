package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	errBelowMin      = errors.New("integer is below minimum")
	errAboveMax      = errors.New("integer is above maximum")
	errIntNotInSet   = errors.New("integer is not in the set")
	errStrNotInSet   = errors.New("string is not in the set")
	errRegexMismatch = errors.New("regexp doesn't match the value")
	errLenMismatch   = errors.New("string length doesn't match")
	errSliceEmpty    = errors.New("the slice is empty")
	errUnsuppRule    = errors.New("unsupported rule for the type")
	errUnsuppType    = errors.New("unsupported type")
)

func (v ValidationErrors) Error() string {
	strBuilder := strings.Builder{}
	for _, err := range v {
		strBuilder.WriteString(fmt.Sprintf("%s: %s\n", err.Field, err.Err.Error()))
	}
	return strings.TrimRight(strBuilder.String(), "\n")
}

// method for collecting all validation errors
func (v *ValidationErrors) logErr(name string, err error) {
	*v = append(*v, ValidationError{
		Field: name,
		Err:   err,
	})
}

func Validate(v interface{}) error {
	if t := reflect.TypeOf(v).Kind().String(); t != "struct" {
		return errors.New("incorrect type, need struct")
	}

	var err error

	valErrors := make(ValidationErrors, 0, 20)
	fields := reflect.VisibleFields(reflect.TypeOf(v))

	for i, f := range fields {
		// a field is ignored if there is no "validate" tage
		valStr, ok := f.Tag.Lookup("validate")
		if !ok {
			continue
		}
		rules := strings.Split(valStr, "|")

		// a loop to validate each individuall field against each rule
		// in case there are multiple rules
		for _, r := range rules {
			rule := strings.Split(r, ":")

			// call a specific function based on a field type
			// nested structs will recursively call Validate func
			switch f.Type.Kind().String() {
			case "struct":
				if valStr == "nested" {
					validation := Validate(reflect.ValueOf(v).Field(i).Interface())
					if validation != nil {
						valErrors.logErr(f.Name, validation)
					}
				}
			case "slice":
				// throw error if the slice is empty, an attempt to continue will cause panic
				if reflect.ValueOf(v).Field(i).Len() == 0 {
					valErrors.logErr(f.Name, errSliceEmpty)
					continue
				}
				err = valSlice(rule, &valErrors, reflect.ValueOf(v).Field(i).Interface(), f.Name)
			case "string":
				err = valString(rule, &valErrors, reflect.ValueOf(v).Field(i).String(), f.Name)
			case "int", "int8", "int16", "int32", "int64":
				err = valInt(rule, &valErrors, reflect.ValueOf(v).Field(i).Int(), f.Name)
			default:
				valErrors.logErr(f.Name, errUnsuppType)
			}
			if err != nil {
				return err
			}
		}
	}
	if len(valErrors) != 0 {
		return valErrors
	}
	return nil
}

// function responsible for calling a corresponding validation func
// for each of its element
func valSlice(rule []string, valErrors *ValidationErrors, val interface{}, name string) error {
	switch reflect.ValueOf(val).Index(0).Kind().String() {
	case "string":
		for j := range reflect.ValueOf(val).Len() {
			err := valString(rule, valErrors, reflect.ValueOf(val).Index(j).String(), name)
			if err != nil {
				return err
			}
		}
	case "int":
		for j := range reflect.ValueOf(val).Len() {
			err := valInt(rule, valErrors, reflect.ValueOf(val).Index(j).Int(), name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// function responsible for validating strings.
// each case represents a validation rule.
func valString(rule []string, valErrors *ValidationErrors, val string, name string) error {
	switch rule[0] {
	case "len":
		l, err := strconv.Atoi(rule[1])
		if err != nil {
			return err
		}
		if len(val) != l {
			valErrors.logErr(name, errLenMismatch)
		}
	case "regexp":
		reg, err := regexp.Compile(rule[1])
		if err != nil {
			return err
		}
		if !reg.MatchString(val) {
			valErrors.logErr(name, errRegexMismatch)
		}
	case "in":
		set := strings.Split(rule[1], ",")
		if !slices.Contains(set, val) {
			valErrors.logErr(name, errStrNotInSet)
		}
	default:
		valErrors.logErr(name, errUnsuppRule)
	}
	return nil
}

// function responsible for validating integers.
// each case represents a validation rule.
func valInt(rule []string, valErrors *ValidationErrors, val int64, name string) error {
	switch rule[0] {
	case "min":
		lim, err := strconv.Atoi(rule[1])
		if err != nil {
			return err
		}
		if val < int64(lim) {
			valErrors.logErr(name, errBelowMin)
		}
	case "max":
		lim, err := strconv.Atoi(rule[1])
		if err != nil {
			return err
		}
		if val > int64(lim) {
			valErrors.logErr(name, errAboveMax)
		}
	case "in":
		set := strings.Split(rule[1], ",")
		if !slices.Contains(set, strconv.FormatInt(val, 10)) {
			valErrors.logErr(name, errIntNotInSet)
		}
	default:
		valErrors.logErr(name, errUnsuppRule)
	}
	return nil
}
