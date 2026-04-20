package api

type Stream interface {
	Export(client *XtreamClient, dir string) (int, int, int, error)
	Equals(t Stream) bool

	GetCover() string
	GetCategoryIds() []int
}

// Function: Equals()
func (l LiveStream) Equals(t Stream) bool {
	switch any(t).(type) {
	case Movie:
		return false
	case Series:
		return false
	}

	livestream := t.(LiveStream)

	// Compare
	if l.Id != livestream.Id {
		return false
	}

	if l.Name != livestream.Name {
		return false
	}

	if l.Added != livestream.Added {
		return false
	}

	if l.Number != livestream.Number {
		return false
	}

	if l.CategoryId != livestream.CategoryId {
		return false
	}

	if l.EpgId != livestream.EpgId {
		return false
	}

	if l.HasCatchup != livestream.HasCatchup {
		return false
	}

	return true
}

func (m Movie) Equals(t Stream) bool {
	switch any(t).(type) {
	case LiveStream:
		return false
	case Series:
		return false
	}

	movie := t.(Movie)

	// Compare
	if m.Id != movie.Id {
		return false
	}

	if m.Name != movie.Name {
		return false
	}

	if m.Added != movie.Added {
		return false
	}

	if m.Extension != movie.Extension {
		return false
	}

	if m.Number != movie.Number {
		return false
	}

	if m.CategoryId != movie.CategoryId {
		return false
	}

	return true
}

func (s Series) Equals(t Stream) bool {
	switch any(t).(type) {
	case LiveStream:
		return false
	case Movie:
		return false
	}

	show := t.(Series)

	// Compare
	if s.LastModified != show.LastModified {
		return false
	}

	return true
}

// Function: GetCover()
func (l LiveStream) GetCover() string {
	return l.Icon
}

func (m Movie) GetCover() string {
	return m.Icon
}

func (s Series) GetCover() string {
	return s.Cover
}

// Function: GetCategoryIds()
func (l LiveStream) GetCategoryIds() []int {
	return l.CategoryIds
}

func (m Movie) GetCategoryIds() []int {
	return m.CategoryIds
}

func (s Series) GetCategoryIds() []int {
	return s.CategoryIds
}
