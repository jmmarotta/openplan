package plan

import (
	"bytes"
	"errors"
	"testing"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{name: "valid", raw: "OPN-12_7K4M9XQ2"},
		{name: "missing suffix", raw: "OPN-12", wantErr: true},
		{name: "lowercase prefix", raw: "opn-12_7K4M9XQ2", wantErr: true},
		{name: "zero number", raw: "OPN-0_7K4M9XQ2", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseID(tc.raw)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q", tc.raw)
			}
			if tc.wantErr && !errors.Is(err, errInvalidID) {
				t.Fatalf("expected wrapped invalid ID sentinel for %q, got %v", tc.raw, err)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.raw, err)
			}
		})
	}
}

func TestNewID(t *testing.T) {
	id, err := NewID("OPN", []string{"OPN-1_ABCDEFGH", "OPN-3_12345678"}, bytes.NewReader([]byte{0, 0, 0, 0, 0}))
	if err != nil {
		t.Fatalf("NewID returned error: %v", err)
	}

	if got, want := FormatID(id), "OPN-4_00000000"; got != want {
		t.Fatalf("FormatID() = %q, want %q", got, want)
	}
}
