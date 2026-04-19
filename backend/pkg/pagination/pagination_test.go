package pagination

import "testing"

func TestNewParamsDefaults(t *testing.T) {
	p := NewParams(0, 0)
	if p.Page != DefaultPage {
		t.Errorf("expected Page=%d for page=0, got %d", DefaultPage, p.Page)
	}
	if p.PageSize != DefaultPageSize {
		t.Errorf("expected PageSize=%d for pageSize=0, got %d", DefaultPageSize, p.PageSize)
	}
}

func TestNewParamsNegativeValues(t *testing.T) {
	p := NewParams(-5, -10)
	if p.Page != DefaultPage {
		t.Errorf("expected Page=%d for negative, got %d", DefaultPage, p.Page)
	}
	if p.PageSize != DefaultPageSize {
		t.Errorf("expected PageSize=%d for negative, got %d", DefaultPageSize, p.PageSize)
	}
}

func TestNewParamsMaxPageSize(t *testing.T) {
	p := NewParams(1, 500)
	if p.PageSize != MaxPageSize {
		t.Errorf("expected PageSize capped at %d, got %d", MaxPageSize, p.PageSize)
	}
}

func TestNewParamsValidValues(t *testing.T) {
	p := NewParams(3, 25)
	if p.Page != 3 {
		t.Errorf("expected Page=3, got %d", p.Page)
	}
	if p.PageSize != 25 {
		t.Errorf("expected PageSize=25, got %d", p.PageSize)
	}
}

func TestOffset(t *testing.T) {
	tests := []struct {
		page, pageSize, expected int
	}{
		{1, 20, 0},
		{2, 20, 20},
		{3, 10, 20},
		{5, 25, 100},
	}
	for _, tt := range tests {
		p := NewParams(tt.page, tt.pageSize)
		if got := p.Offset(); got != tt.expected {
			t.Errorf("Offset() for page=%d, pageSize=%d: got %d, want %d", tt.page, tt.pageSize, got, tt.expected)
		}
	}
}

func TestNewResult(t *testing.T) {
	tests := []struct {
		page, pageSize, total, expectedPages int
	}{
		{1, 20, 0, 0},
		{1, 20, 1, 1},
		{1, 20, 20, 1},
		{1, 20, 21, 2},
		{1, 20, 100, 5},
		{1, 10, 55, 6},
		{2, 25, 100, 4},
	}
	for _, tt := range tests {
		r := NewResult(tt.page, tt.pageSize, tt.total)
		if r.TotalPages != tt.expectedPages {
			t.Errorf("NewResult(%d,%d,%d).TotalPages = %d, want %d",
				tt.page, tt.pageSize, tt.total, r.TotalPages, tt.expectedPages)
		}
		if r.Page != tt.page {
			t.Errorf("NewResult page = %d, want %d", r.Page, tt.page)
		}
		if r.PageSize != tt.pageSize {
			t.Errorf("NewResult pageSize = %d, want %d", r.PageSize, tt.pageSize)
		}
		if r.Total != tt.total {
			t.Errorf("NewResult total = %d, want %d", r.Total, tt.total)
		}
	}
}
