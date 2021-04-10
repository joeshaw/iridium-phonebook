package main

import "testing"

func TestUCS2(t *testing.T) {
	testcases := []struct {
		utf8    string
		ucs2hex string
	}{
		{"", ""},
		{"abc123", "006100620063003100320033"},
		{"Joe Shaw", "004a006f006500200053006800610077"},
	}

	for _, tc := range testcases {
		t.Run(tc.utf8, func(t *testing.T) {
			e, err := encodeUCS2Hex(tc.utf8)
			if err != nil {
				t.Errorf("encodeUCS2Hex: got unexpected error %v", err)
			} else if got, want := e, tc.ucs2hex; got != want {
				t.Errorf("encodeUCS2Hex: got %q, want %q", got, want)
			}

			d, err := decodeUCS2Hex(tc.ucs2hex)
			if err != nil {
				t.Errorf("decodeUCS2Hex: got unexpected error %v", err)
			} else if got, want := d, tc.utf8; got != want {
				t.Errorf("decodeUCS2Hex: got %q, want %q", got, want)
			}
		})
	}
}
