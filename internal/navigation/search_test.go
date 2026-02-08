package navigation

import (
	"strings"
	"testing"
)

type mockModel struct {
	slides [][]string // Each slide has multiple parts
	page   int
	part   int
}

func (m *mockModel) CurrentPage() int {
	return m.page
}

func (m *mockModel) CurrentPart() int {
	return m.part
}

func (m *mockModel) SetPage(page int) {
	m.page = page
	m.part = 0
}

func (m *mockModel) SetPageAndPart(page int, part int) {
	m.page = page
	m.part = part
}

func (m *mockModel) Pages() []string {
	pages := make([]string, len(m.slides))
	for i, parts := range m.slides {
		pages[i] = strings.Join(parts, "\n")
	}
	return pages
}

func (m *mockModel) SlideParts(page int) []string {
	if page < 0 || page >= len(m.slides) {
		return nil
	}
	return m.slides[page]
}

func TestSearch(t *testing.T) {
	data := []string{
		"hi",
		"first",
		"second",
		"third",
		"AbCdEfG",
		"abcdefg",
		"seconds",
	}

	type query struct {
		desc     string
		query    string
		expected int
	}

	// query -> expected page
	queries := []query{
		{"basic 'first'", "first", 1},
		{"basic 'abc'", "abc", 5},
		{"basic 'abc' next occurrence", "abc", 5},
		{"'abc' ignore case", "abc/i", 4},
		{"'abc' ignore case", "abc/i", 5},
		{"'abc' ignore case", "abc/i", 4},
		{"next occurrence 1/2", "sec", 6},
		{"next occurrence 2/2", "sec", 2},
		{"regex", "a.c", 5},
		{"regex next occurrence", "a.c", 5},
		{"regex ignore case", "a.c/i", 4},
		{"regex ignore case next occurrence", "a.c/i", 5},
	}

	slides := make([][]string, len(data))
	for i, slide := range data {
		slides[i] = []string{slide}
	}

	m := &mockModel{
		slides: slides,
		page:   0,
		part:   0,
	}

	s := &Search{}
	for _, query := range queries {
		s.SetQuery(query.query)
		s.Execute(m)
		if m.CurrentPage() != query.expected {
			t.Errorf("[%s] expected page %d, got %d", query.desc, query.expected, m.CurrentPage())
		}
	}
}

func TestSearchWithPausedSlides(t *testing.T) {
	// Create slides with multiple parts
	slides := [][]string{
		{"Part A1", "Part A2", "Part A3"},
		{"Part B1", "Part B2"},
		{"Part C1"},
	}

	tests := []struct {
		name         string
		startPage    int
		startPart    int
		query        string
		expectedPage int
		expectedPart int
	}{
		{
			name:         "Find in next part of current slide",
			startPage:    0,
			startPart:    0,
			query:        "A2",
			expectedPage: 0,
			expectedPart: 1,
		},
		{
			name:         "Find in last part of current slide",
			startPage:    0,
			startPart:    0,
			query:        "A3",
			expectedPage: 0,
			expectedPart: 2,
		},
		{
			name:         "Find in next slide",
			startPage:    0,
			startPart:    2,
			query:        "B1",
			expectedPage: 1,
			expectedPart: 0,
		},
		{
			name:         "Find in next slide, second part",
			startPage:    0,
			startPart:    2,
			query:        "B2",
			expectedPage: 1,
			expectedPart: 1,
		},
		{
			name:         "Wrap around to find in first slide",
			startPage:    2,
			startPart:    0,
			query:        "A1",
			expectedPage: 0,
			expectedPart: 0,
		},
		{
			name:         "Don't find current part (should find next occurrence)",
			startPage:    0,
			startPart:    1,
			query:        "A2",
			expectedPage: 0,
			expectedPart: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockModel{
				slides: slides,
				page:   tt.startPage,
				part:   tt.startPart,
			}

			s := &Search{}
			s.SetQuery(tt.query)
			s.Execute(m)

			if m.CurrentPage() != tt.expectedPage {
				t.Errorf("expected page %d, got %d", tt.expectedPage, m.CurrentPage())
			}
			if m.CurrentPart() != tt.expectedPart {
				t.Errorf("expected part %d, got %d", tt.expectedPart, m.CurrentPart())
			}
		})
	}
}

