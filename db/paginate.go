package db

type Paginate[T any] struct {
	Last       bool   `json:"last"`
	Lists      []T    `json:"lists"`
	Total      uint64 `json:"total"`
	Page       uint   `json:"page"`
	PageSize   uint   `json:"pageSize"`
	TotalPages uint   `json:"totalPages"`
}

func (p Paginate[T]) Map[R any](fc func(T) R) Paginate[R] {
	r := make([]R, len(p.Lists), cap(p.Lists))
	for i := 0; i < len(p.Lists); i++ {
		r[i] = fc(p.Lists[i])
	}
	return Paginate[R]{
		Last:       p.Last,
		Lists:      r,
		Total:      p.Total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: p.TotalPages,
	}
}
