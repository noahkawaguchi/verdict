package api

import "testing"

func TestGetFirstSegment(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"/health", "/health", false},
		{"/polls/abc123", "/polls", false},
		{"/results/abc123/anything-else-here", "/results", false},
		{"noslash", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		got, err := getFirstSegment(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf(
				"getFirstSegment(%q): wantErr=%t, got err=%v",
				tt.input,
				tt.wantErr,
				err,
			)
		}
		if got != tt.want {
			t.Errorf("getFirstSegment(%q): want %q, got %q", tt.input, tt.want, got)
		}
	}
}
