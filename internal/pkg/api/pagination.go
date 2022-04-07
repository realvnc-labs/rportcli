package api

import (
	"net/url"
	"strconv"

	options "github.com/breathbath/go_utils/v2/pkg/config"
)

const (
	PaginationOffset = "offset"
	PaginationLimit  = "limit"
)

type Pagination struct {
	Limit  int
	Offset int
}

func NewPaginationFromParams(params *options.ParameterBag) Pagination {
	return Pagination{
		Offset: params.ReadInt(PaginationOffset, 0),
		Limit:  params.ReadInt(PaginationLimit, -1),
	}
}

func NewPaginationWithLimit(max int) Pagination {
	return Pagination{
		Offset: 0,
		Limit:  max,
	}
}

func (p Pagination) Apply(q url.Values) {
	q.Set("page[offset]", strconv.Itoa(p.Offset))
	q.Set("page[limit]", strconv.Itoa(p.Limit))
}
