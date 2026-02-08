package model

import "testing"

func TestSplitSlidesWithPausedDelimiter(t *testing.T) {
	slides := []string{
		"# Title\nfirst\n+++\nsecond\n+++\nthird",
		"# Next\nonly",
	}

	parsed := splitSlides(slides)
	if len(parsed) != 2 {
		t.Fatalf("expected 2 slides, got %d", len(parsed))
	}
	if len(parsed[0].Parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parsed[0].Parts))
	}
	if parsed[0].Parts[1] != "second" {
		t.Fatalf("expected second part to be %q, got %q", "second", parsed[0].Parts[1])
	}
	if len(parsed[1].Parts) != 1 {
		t.Fatalf("expected 1 part, got %d", len(parsed[1].Parts))
	}
}

func TestAdvanceRetreatSteps(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"a", "b"}},
			{Parts: []string{"c", "d", "e"}},
		},
	}

	m.SetPageAndPart(0, 0)
	m.advanceSteps(1)
	if m.Page != 0 || m.Part != 1 {
		t.Fatalf("expected page 0 part 1, got page %d part %d", m.Page, m.Part)
	}

	m.advanceSteps(1)
	if m.Page != 1 || m.Part != 0 {
		t.Fatalf("expected page 1 part 0, got page %d part %d", m.Page, m.Part)
	}

	m.retreatSteps(1)
	if m.Page != 0 || m.Part != 1 {
		t.Fatalf("expected page 0 part 1 after retreat, got page %d part %d", m.Page, m.Part)
	}
}

func TestCurrentSlideContentCumulative(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"first", "second", "third"}},
		},
	}
	m.SetPageAndPart(0, 1)

	content := m.currentSlideContent()
	if content != "first\nsecond" {
		t.Fatalf("expected cumulative content, got %q", content)
	}
}

func TestPagingFormats(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"a", "b", "c"}},
			{Parts: []string{"d"}},
		},
	}
	m.SetPageAndPart(0, 1)

	tests := []struct {
		format   string
		expected string
	}{
		// No placeholder
		{"Static", "Static"},
		// 1 placeholder: slide
		{"Slide %d", "Slide 1"},
		// 2 placeholders: slide, totalSlides
		{"Slide %d / %d", "Slide 1 / 2"},
		// 3 placeholders: slide, totalSlides, part
		{"Slide %d / %d | Part %d", "Slide 1 / 2 | Part 2"},
		// 4+ placeholders: slide, totalSlides, part, totalParts
		{"Slide %d / %d (Part %d / %d)", "Slide 1 / 2 (Part 2 / 3)"},
	}

	for _, tt := range tests {
		m.Paging = tt.format
		result := m.paging()
		if result != tt.expected {
			t.Errorf("format %q: expected %q, got %q", tt.format, tt.expected, result)
		}
	}
}

func TestPagingWithLastPart(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"a", "b"}},
			{Parts: []string{"c", "d", "e"}},
		},
	}
	m.SetPageAndPart(1, 2) // Last slide, last part
	m.Paging = "Slide %d / %d (Part %d / %d)"

	result := m.paging()
	expected := "Slide 2 / 2 (Part 3 / 3)"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSplitSlidesEmptySlide(t *testing.T) {
	slides := []string{"", "content", ""}
	parsed := splitSlides(slides)
	
	if len(parsed) != 3 {
		t.Fatalf("expected 3 slides, got %d", len(parsed))
	}
	
	if len(parsed[0].Parts) != 1 || parsed[0].Parts[0] != "" {
		t.Errorf("expected empty slide to have 1 empty part")
	}
	
	if len(parsed[1].Parts) != 1 || parsed[1].Parts[0] != "content" {
		t.Errorf("expected content slide to have 1 part with 'content'")
	}
}

func TestAdvanceStepsMultiple(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"a", "b"}},
			{Parts: []string{"c"}},
			{Parts: []string{"d", "e", "f"}},
		},
	}
	
	m.SetPageAndPart(0, 0)
	m.advanceSteps(5)
	// Should be: (0,0) -> (0,1) -> (1,0) -> (2,0) -> (2,1) -> (2,2)
	if m.Page != 2 || m.Part != 2 {
		t.Errorf("after 5 steps from (0,0): expected page 2 part 2, got page %d part %d", m.Page, m.Part)
	}
}

func TestRetreatStepsMultiple(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"a", "b"}},
			{Parts: []string{"c"}},
			{Parts: []string{"d", "e", "f"}},
		},
	}
	
	m.SetPageAndPart(2, 2)
	m.retreatSteps(5)
	// Should be: (2,2) -> (2,1) -> (2,0) -> (1,0) -> (0,1) -> (0,0)
	if m.Page != 0 || m.Part != 0 {
		t.Errorf("after 5 retreat steps from (2,2): expected page 0 part 0, got page %d part %d", m.Page, m.Part)
	}
}

func TestAdvanceStepsAtEnd(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"a"}},
			{Parts: []string{"b", "c"}},
		},
	}
	
	m.SetPageAndPart(1, 1) // At the end
	m.advanceSteps(10) // Try to advance way past the end
	
	// Should stay at the end
	if m.Page != 1 || m.Part != 1 {
		t.Errorf("advancing past end: expected page 1 part 1, got page %d part %d", m.Page, m.Part)
	}
}

func TestRetreatStepsAtStart(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"a", "b"}},
			{Parts: []string{"c"}},
		},
	}
	
	m.SetPageAndPart(0, 0) // At the start
	m.retreatSteps(10) // Try to retreat way before the start
	
	// Should stay at the start
	if m.Page != 0 || m.Part != 0 {
		t.Errorf("retreating past start: expected page 0 part 0, got page %d part %d", m.Page, m.Part)
	}
}

func TestCurrentSlideContentAtEnd(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"first", "second", "third"}},
		},
	}
	m.SetPageAndPart(0, 2)
	
	content := m.currentSlideContent()
	expected := "first\nsecond\nthird"
	if content != expected {
		t.Errorf("expected %q, got %q", expected, content)
	}
}

func TestClampPositionInvalidPage(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"a"}},
			{Parts: []string{"b"}},
		},
		Page: 10,
		Part: 5,
	}
	
	m.clampPosition()
	
	if m.Page != 1 || m.Part != 0 {
		t.Errorf("clamping invalid position: expected page 1 part 0, got page %d part %d", m.Page, m.Part)
	}
}

func TestClampPositionEmptySlides(t *testing.T) {
	m := Model{
		Slides: []Slide{},
		Page:   5,
		Part:   3,
	}
	
	m.clampPosition()
	
	if m.Page != 0 || m.Part != 0 {
		t.Errorf("clamping with empty slides: expected page 0 part 0, got page %d part %d", m.Page, m.Part)
	}
}
