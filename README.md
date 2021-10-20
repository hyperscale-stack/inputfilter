Hyperscale InputFilter [![Last release](https://img.shields.io/github/release/hyperscale-stack/inputfilter.svg)](https://github.com/hyperscale-stack/inputfilter/releases/latest) [![Documentation](https://godoc.org/github.com/hyperscale-stack/inputfilter?status.svg)](https://godoc.org/github.com/hyperscale-stack/inputfilter)
====================

[![Go Report Card](https://goreportcard.com/badge/github.com/hyperscale-stack/inputfilter)](https://goreportcard.com/report/github.com/hyperscale-stack/inputfilter)

| Branch  | Status | Coverage |
|---------|--------|----------|
| master  | [![Build Status](https://github.com/hyperscale-stack/inputfilter/workflows/Go/badge.svg?branch=master)](https://github.com/hyperscale-stack/inputfilter/actions?query=workflow%3AGo) | [![Coveralls](https://img.shields.io/coveralls/hyperscale-stack/inputfilter/master.svg)](https://coveralls.io/github/hyperscale-stack/inputfilter?branch=master) |

The Hyperscale InputFilter library provides a simple inputfilter chaining mechanism by which multiple filters and validator may be applied to a single datum in a user-defined order. 

## Example

Filter by `map[string]interface{}`

```go
package main

import (
    "fmt"

    "github.com/hyperscale-stack/filter"
    "github.com/hyperscale-stack/validator"
    "github.com/hyperscale-stack/inputfilter"
)

func main() {
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

	data, errs := i.FilterMap(map[string]interface{}{
        "id":  "9D2C8507-5F9D-4CB0-A098-2E307B39DC91",
        "url": "HTTPS://google.COM",
    })
    // return 
    // map[string]interface{}{
	//     "id":  "9d2c8507-5f9d-4cb0-a098-2e307b39dc91",
    //     "url": "https://google.com",
    // }
}

```


Filter by `url.Values`

```go
package main

import (
    "fmt"

    "github.com/hyperscale-stack/filter"
    "github.com/hyperscale-stack/validator"
    "github.com/hyperscale-stack/inputfilter"
)

func main() {
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

    values := url.Values{}
    values.Set("id", "9D2C8507-5F9D-4CB0-A098-2E307B39DC91")
    values.Set("url", "HTTPS://google.COM")

	data, errs := i.FilterValues(values)
    // return 
    // url.Values{
	//     "id":  []string{"9d2c8507-5f9d-4cb0-a098-2e307b39dc91"},
    //     "url": []string{"https://google.com"},
    // }
}

```


## License

Hyperscale Filter is licensed under [the MIT license](LICENSE.md).
