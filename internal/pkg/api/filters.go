package api

import (
	"fmt"
	"net/url"
)

type Filters map[string]string

func NewFilters(keyValues ...string) Filters {
	f := make(map[string]string)
	for i := 0; 2*i+1 < len(keyValues); i++ {
		f[keyValues[2*i]] = keyValues[2*i+1]
	}
	return f
}

func (f Filters) Apply(q url.Values) {
	for k, v := range f {
		if v != "" {
			q.Set(fmt.Sprintf("filter[%s]", k), v)
		}
	}
}
