package api

type Stream interface {
	Export(client *XtreamClient, dir string) (int, int, int, error)
	GetCategoryIds() []int
}

func (l LiveStream) GetCategoryIds() []int {
	return l.CategoryIds
}

func (m Movie) GetCategoryIds() []int {
	return m.CategoryIds
}

func (s Series) GetCategoryIds() []int {
	return s.CategoryIds
}
