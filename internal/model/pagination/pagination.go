package pagination

import "errors"

type Pagination struct {
	page  uint64
	limit uint64
}

func New(page, limit uint64) *Pagination {
	return &Pagination{
		page:  page,
		limit: limit,
	}
}

func (pagination *Pagination) Page() uint64 {
	return pagination.page
}

func (pagination *Pagination) Limit() uint64 {
	return pagination.limit
}

func (pagination *Pagination) Offset() uint64 {
	return (pagination.page - 1) * pagination.limit
}

func (pagination *Pagination) Validate() []error {
	var result []error

	if pagination.page < 1 {
		result = append(result, errors.New("page must be greater than 0"))
	}

	if pagination.limit < 1 {
		result = append(result, errors.New("limit must be greater than 0"))
	}

	return result
}
