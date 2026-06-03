package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{"ascii string", []byte("hello world"), "hello world"},
		{"empty", []byte{}, ""},
		{"utf8 string", []byte("héllo"), "héllo"},
		{"binary bytes", []byte{0x01, 0x02, 0x03}, "0x010203"},
		{"tab is printable", []byte("key\tval"), "key\tval"},
		{"null byte", []byte{0x00}, "0x00"},
		{"mixed printable and binary", []byte{'h', 'i', 0x00}, "0x686900"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatValue(tc.input)
			if got != tc.want {
				t.Errorf("formatValue(%v) = %q; want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseFormat(t *testing.T) {
	if _, err := ParseFormat("table"); err != nil {
		t.Errorf("expected no error for 'table', got %v", err)
	}
	if _, err := ParseFormat("json"); err != nil {
		t.Errorf("expected no error for 'json', got %v", err)
	}
	if _, err := ParseFormat("csv"); err == nil {
		t.Error("expected error for 'csv', got nil")
	}
}

func TestPrintEntries_Table(t *testing.T) {
	entries := []Entry{
		{Key: "b", Value: []byte("two")},
		{Key: "a", Value: []byte("one")},
	}
	var buf bytes.Buffer
	if err := PrintEntries(&buf, entries, FormatTable); err != nil {
		t.Fatalf("PrintEntries error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "VALUE") {
		t.Error("expected header row")
	}
	aIdx := strings.Index(out, "a")
	bIdx := strings.Index(out, "b")
	if aIdx > bIdx {
		t.Error("expected entries sorted by key")
	}
}

func TestPrintEntries_JSON(t *testing.T) {
	entries := []Entry{
		{Key: "k", Value: []byte("v")},
	}
	var buf bytes.Buffer
	if err := PrintEntries(&buf, entries, FormatJSON); err != nil {
		t.Fatalf("PrintEntries error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"key"`) || !strings.Contains(out, `"value"`) {
		t.Errorf("expected JSON keys in output, got: %s", out)
	}
}

func TestPrintEntries_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := PrintEntries(&buf, nil, FormatTable); err != nil {
		t.Fatalf("PrintEntries error: %v", err)
	}
	if !strings.Contains(buf.String(), "KEY") {
		t.Error("expected header even for empty slice")
	}
}

func TestPrintEntry_Table(t *testing.T) {
	var buf bytes.Buffer
	e := Entry{Key: "mykey", Value: []byte{0xff, 0xfe}} // invalid UTF-8
	if err := PrintEntry(&buf, e, FormatTable); err != nil {
		t.Fatalf("PrintEntry error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "mykey") || !strings.Contains(out, "0xfffe") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestPrintEntry_JSON(t *testing.T) {
	var buf bytes.Buffer
	e := Entry{Key: "k", Value: []byte("v")}
	if err := PrintEntry(&buf, e, FormatJSON); err != nil {
		t.Fatalf("PrintEntry error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"key": "k"`) {
		t.Errorf("unexpected JSON output: %s", out)
	}
}
