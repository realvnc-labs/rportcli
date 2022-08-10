package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/sirupsen/logrus"
)

type Filters map[string]string

// NewFilters constructs filters from key value pairs. Keys with empty values are ignored.
func NewFilters(keyValues ...string) Filters {
	f := make(map[string]string)
	for i := 0; 2*i+1 < len(keyValues); i++ {
		f[keyValues[2*i]] = keyValues[2*i+1]
	}
	return f
}

// NewFilterFromCombinedSearchString takes a list of key value filters separated by ampersand,
// e.g. name=elmo&os_kernel=linux, implements --combined-search
func NewFilterFromCombinedSearchString(search string) (Filters, error) {
	logrus.Debugf("Got combined search string '%s'", search)
	return NewFilterFromKVStrings(strings.Split(search, "&"))
}

// NewFilterFromKVStrings takes a slice of strings with key value separated by '='
func NewFilterFromKVStrings(searchFlags []string) (Filters, error) {
	f := make(map[string]string)
	for _, flag := range searchFlags {
		if !strings.Contains(flag, "=") {
			continue
		}
		parts := strings.Split(strings.Trim(flag, " "), "=")
		key := strings.Trim(parts[0], " ")
		if f[key] != "" {
			return f, fmt.Errorf("you cannot specify '--%s %s=' twice", config.ClientSearchFlag, key)
		}
		f[key] = strings.Trim(parts[1], " ")
	}
	return f, nil
}

func (f Filters) Apply(q url.Values) {
	for k, v := range f {
		if v != "" {
			q.Set(fmt.Sprintf("filter[%s]", k), v)
		}
	}
}
