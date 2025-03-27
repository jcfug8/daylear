package model

// NewPageToken creates a new page token with the given limit, offset, and
// tail.
func NewPageToken[T any](pageSize, skip int32, checksum uint32, tail *T) *PageToken[T] {
	return &PageToken[T]{
		PageSize:        pageSize,
		RequestChecksum: checksum,
		Skip:            skip,
		Tail:            tail,
	}
}

// PageToken represents a page token for paginated requests.
type PageToken[T any] struct {
	PageSize        int32
	RequestChecksum uint32
	Skip            int32
	Tail            *T
}

// GetPageSize returns the page size of the page token.
func (p *PageToken[T]) GetPageSize() int32 {
	if p != nil {
		return p.PageSize
	}
	return 0
}

// GetRequestChecksum returns the request checksum of the page token.
func (p *PageToken[T]) GetRequestChecksum() uint32 {
	if p != nil {
		return p.RequestChecksum
	}
	return 0
}

// GetSkip returns the skip value of the page token.
func (p *PageToken[T]) GetSkip() int32 {
	if p != nil {
		return p.Skip
	}
	return 0
}

// GetTail returns the tail value of the page token.
func (p *PageToken[T]) GetTail() *T {
	if p != nil {
		return p.Tail
	}
	return nil
}

// ----------------------------------------------------------------------------

// Next returns a new page token for the next page.
func (t *PageToken[T]) Next(page []T) *PageToken[T] {
	if t == nil || len(page) == 0 || int32(len(page)) < t.PageSize {
		return nil
	}

	return &PageToken[T]{
		PageSize:        t.PageSize,
		RequestChecksum: t.RequestChecksum,
		Tail:            &page[len(page)-1],
	}
}
