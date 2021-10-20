// Copyright 2021 Hyperscale. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package inputfilter

import (
	"fmt"
	"net/url"

	"github.com/hyperscale-stack/filter"
	"github.com/hyperscale-stack/validator"
)

// InputDefinition struct.
type InputDefinition struct {
	Filters    []filter.Filter
	Validators []validator.Validator
}

// InputFilter interface.
type InputFilter interface {
	FilterMap(input map[string]interface{}) (map[string]interface{}, map[string][]error)
	FilterValues(values url.Values) (url.Values, map[string][]error)
}

type inputFilter struct {
	filters    map[string][]filter.Filter
	validators map[string][]validator.Validator
}

// New input filter.
func New(definitions map[string]InputDefinition) InputFilter {
	f := &inputFilter{
		filters:    make(map[string][]filter.Filter),
		validators: make(map[string][]validator.Validator),
	}

	globalFilters := []filter.Filter{}

	for key, def := range definitions {
		// extract filters for wildcard field
		if key == "*" {
			globalFilters = append(globalFilters, def.Filters...)

			continue
		}

		if len(def.Filters) > 0 {
			f.filters[key] = append(f.filters[key], def.Filters...)
		}

		if len(def.Validators) > 0 {
			f.validators[key] = append(f.validators[key], def.Validators...)
		}
	}

	// apply global filter to all input fields
	if len(globalFilters) > 0 {
		for key := range definitions {
			f.filters[key] = append(f.filters[key], globalFilters...)
		}
	}

	return f
}

func (f inputFilter) filterField(key string, value interface{}) (interface{}, error) {
	filters, ok := f.filters[key]
	if !ok {
		return value, nil
	}

	val := value

	var err error

	for _, filter := range filters {
		val, err = filter.Filter(val)
		if err != nil {
			return val, err
		}
	}

	return val, nil
}

func (f inputFilter) validateField(key string, value interface{}) []error {
	errs := []error{}

	validators, ok := f.validators[key]
	if !ok {
		return errs
	}

	for _, validator := range validators {
		if err := validator.Validate(value); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (f inputFilter) FilterMap(input map[string]interface{}) (map[string]interface{}, map[string][]error) {
	errs := make(map[string][]error, len(input))

	for key, val := range input {
		val, err := f.filterField(key, val)
		if err != nil {
			errs[key] = append(errs[key], err)

			continue
		}

		input[key] = val
	}

	for key, val := range input {
		if err := f.validateField(key, val); len(err) > 0 {
			errs[key] = append(errs[key], err...)
		}
	}

	return input, errs
}

func (f inputFilter) filterFieldValues(key string, values []interface{}) ([]interface{}, error) {
	filters, ok := f.filters[key]
	if !ok {
		return values, nil
	}

	retvals := make([]interface{}, len(values))

	for i, value := range values {
		val := value

		var err error

		for _, filter := range filters {
			val, err = filter.Filter(val)
			if err != nil {
				return retvals, err
			}

			retvals[i] = val
		}
	}

	return retvals, nil
}

func (f inputFilter) castSliceStringToSliceValue(values []string) []interface{} {
	retvals := make([]interface{}, len(values))

	for i, val := range values {
		retvals[i] = val
	}

	return retvals
}

func (f inputFilter) castSliceValueToSliceString(values []interface{}) []string {
	retvals := make([]string, len(values))

	for i, val := range values {
		switch v := val.(type) {
		case string:
			retvals[i] = v
		default:
			retvals[i] = fmt.Sprintf("%v", val)
		}
	}

	return retvals
}

func (f inputFilter) FilterValues(values url.Values) (url.Values, map[string][]error) {
	errs := make(map[string][]error, len(values))

	for key, vals := range values {
		vals, err := f.filterFieldValues(key, f.castSliceStringToSliceValue(vals))
		if err != nil {
			errs[key] = append(errs[key], err)

			continue
		}

		values[key] = f.castSliceValueToSliceString(vals)
	}

	for key, vals := range values {
		for _, val := range vals {
			if err := f.validateField(key, val); len(err) > 0 {
				errs[key] = append(errs[key], err...)
			}
		}
	}

	return values, errs
}
