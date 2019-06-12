package resource

type SearchFilter map[string]interface{}

type SearchPagination struct{ page, size uint }

const (
	defaultPageNumber = 1
	defaultPageSize   = 100
)

func NewSearchPagination(page, size uint) *SearchPagination {
	if page == 0 {
		page = defaultPageNumber
	}
	if size == 0 {
		size = defaultPageSize
	}
	return &SearchPagination{
		page: page,
		size: size,
	}
}
func (p *SearchPagination) Size() uint {
	return p.size
}

func (p *SearchPagination) Page() uint {
	return p.page
}

func (p *SearchPagination) Limit() uint {
	return p.Size()
}

func (p *SearchPagination) Offset() uint {
	return p.Size() * (p.page - 1)
}
