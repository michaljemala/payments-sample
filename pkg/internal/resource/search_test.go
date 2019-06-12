package resource

import "testing"

func TestNewSearchPagination(t *testing.T) {
	testCases := []struct {
		name          string
		page, size    uint
		offset, limit uint
	}{
		{name: "For zero size use default", page: 1, size: 0, offset: 0, limit: defaultPageSize},
		{name: "For zero page use default", page: 0, size: 10, offset: 0, limit: 10},
		{name: "Correct paging 1", page: 1, size: 10, offset: 0, limit: 10},
		{name: "Correct paging 2", page: 5, size: 10, offset: 40, limit: 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewSearchPagination(tc.page, tc.size)

			limit, offset := p.Limit(), p.Offset()
			if want, have := tc.limit, limit; want != have {
				t.Errorf("unexpected limit: want %d, have %d", want, have)
			}
			if want, have := tc.offset, offset; want != have {
				t.Errorf("unexpected offset: want %d, have %d", want, have)
			}
		})
	}
}
