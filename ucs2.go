package main

import (
	"encoding/hex"

	"golang.org/x/text/encoding/unicode"
)

// UCS-2 is a strict subset of UTF-16
func decodeUCS2Hex(ucs2Hex string) (string, error) {
	bytes, err := hex.DecodeString(ucs2Hex)
	if err != nil {
		return "", err
	}

	e := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	dec := e.NewDecoder()
	utf8, err := dec.Bytes(bytes)
	if err != nil {
		return "", err
	}
	return string(utf8), nil
}

// Ok, this is technically returning UTF-16 instead of UCS-2, but I
// don't care.  Just don't stick emoji in your strings.
func encodeUCS2Hex(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	e := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	enc := e.NewEncoder()
	utf16, err := enc.Bytes([]byte(s))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(utf16), nil
}
