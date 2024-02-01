package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestResp_readLine(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		n        int
		err      error
	}{
		{
			name:     "Test reading a line with multiple characters",
			input:    []byte("Hello\r\n"),
			expected: []byte("Hello"),
			n:        7,
			err:      nil,
		},
		{
			name:     "Test reading an empty line",
			input:    []byte("\r\n"),
			expected: []byte(""),
			n:        2,
			err:      nil,
		},
		{
			name:     "Test reading a line that spans multiple read operations",
			input:    []byte("Hello\r\n"),
			expected: []byte("Hello"),
			n:        7,
			err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := NewResp(bytes.NewBuffer(test.input))
			line, n, err := resp.readLine()
			if !reflect.DeepEqual(line, test.expected) {
				t.Errorf("Expected line %v, but got %v", test.expected, line)
			}
			if n != test.n {
				t.Errorf("Expected number of bytes %d, but got %d", test.n, n)
			}
			if err != test.err {
				t.Errorf("Expected error %v, but got %v", test.err, err)
			}
		})
	}
}
