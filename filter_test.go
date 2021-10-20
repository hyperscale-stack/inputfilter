// Copyright 2021 Hyperscale. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package inputfilter

import (
	"net/url"
	"testing"

	"github.com/hyperscale-stack/filter"
	"github.com/hyperscale-stack/validator"
	"github.com/stretchr/testify/assert"
)

func TestInputFilter(t *testing.T) {
	i := New(map[string]InputDefinition{
		"url": {
			Filters: []filter.Filter{
				filter.NewStringToLowerFilter(),
				filter.NewURLFilter(),
			},
		},
		"id": {
			Validators: []validator.Validator{
				validator.NewUUIDValidator(),
			},
		},
		"size": {
			Filters: []filter.Filter{
				filter.NewIntFilter(),
			},
		},
	})

	{
		data, errs := i.FilterMap(map[string]interface{}{
			"id":  "9D2C8507-5F9D-4CB0-A098-2E307B39DC91",
			"url": "HTTPS://google.COM",
		})
		assert.Equal(t, 0, len(errs))
		assert.Equal(t, "9D2C8507-5F9D-4CB0-A098-2E307B39DC91", data["id"])
		assert.Equal(t, "https://google.com", data["url"])
	}

	{
		values := url.Values{}
		values.Set("id", "9D2C8507-5F9D-4CB0-A098-2E307B39DC91")
		values.Set("url", "HTTPS://google.COM")
		values.Set("size", "123")

		data, errs := i.FilterValues(values)
		assert.Equal(t, 0, len(errs))
		assert.Equal(t, "9D2C8507-5F9D-4CB0-A098-2E307B39DC91", data.Get("id"))
		assert.Equal(t, "https://google.com", data.Get("url"))
		assert.Equal(t, "123", data.Get("size"))
	}
}

func TestInputFilterWithGlobalFilter(t *testing.T) {
	i := New(map[string]InputDefinition{
		"*": {
			Filters: []filter.Filter{
				filter.NewStringToLowerFilter(),
			},
		},
		"url": {
			Filters: []filter.Filter{
				filter.NewURLFilter(),
			},
		},
		"id": {
			Validators: []validator.Validator{
				validator.NewUUIDValidator(),
			},
		},
	})

	{
		data, errs := i.FilterMap(map[string]interface{}{
			"id":  "9D2C8507-5F9D-4CB0-A098-2E307B39DC91",
			"url": "HTTPS://google.COM",
		})
		assert.Equal(t, 0, len(errs))
		assert.Equal(t, "9d2c8507-5f9d-4cb0-a098-2e307b39dc91", data["id"])
		assert.Equal(t, "https://google.com", data["url"])
	}

	{
		values := url.Values{}
		values.Set("id", "9D2C8507-5F9D-4CB0-A098-2E307B39DC91")
		values.Set("url", "HTTPS://google.COM")

		data, errs := i.FilterValues(values)
		assert.Equal(t, 0, len(errs))
		assert.Equal(t, "9d2c8507-5f9d-4cb0-a098-2e307b39dc91", data.Get("id"))
		assert.Equal(t, "https://google.com", data.Get("url"))
	}
}

func TestInputFilterWithBadInput(t *testing.T) {
	i := New(map[string]InputDefinition{
		"*": {
			Filters: []filter.Filter{
				filter.NewStringToLowerFilter(),
			},
		},
		"url": {
			Filters: []filter.Filter{
				filter.NewURLFilter(),
			},
		},
		"id": {
			Validators: []validator.Validator{
				validator.NewUUIDValidator(),
			},
		},
	})

	_, errs := i.FilterMap(map[string]interface{}{
		"id":  "9D2C8507-5F9D-4CB0-A098-2E307",
		"url": 1244,
	})
	assert.Equal(t, 2, len(errs))

	assert.Contains(t, errs, "id")
	assert.Error(t, errs["id"][0])

	assert.Contains(t, errs, "url")
	assert.Error(t, errs["url"][0])

	values := url.Values{}
	values.Set("id", "9D2C8507-5F9D-4CB0-A098-2E307")
	values.Set("url", "134://foo")

	_, errs = i.FilterValues(values)
	assert.Equal(t, 2, len(errs))

	assert.Contains(t, errs, "id")
	assert.Error(t, errs["id"][0])

	assert.Contains(t, errs, "url")
	assert.Error(t, errs["url"][0])
}
