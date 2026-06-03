package output

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"unicode"
	"unicode/utf8"
)

type Entry struct {
	Key   string
	Value []byte
}

type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatTable, FormatJSON:
		return Format(s), nil
	default:
		return "", fmt.Errorf("unknown output format %q: must be table or json", s)
	}
}

// formatValue renders a byte slice as a printable string.
// Valid printable UTF-8 is returned as-is; anything else becomes "0x<hex>".
func formatValue(b []byte) string {
	if utf8.Valid(b) {
		allPrintable := true
		for _, r := range string(b) {
			if !unicode.IsPrint(r) && r != '\t' && r != '\n' {
				allPrintable = false
				break
			}
		}
		if allPrintable {
			return string(b)
		}
	}
	return "0x" + hex.EncodeToString(b)
}

// errWriter wraps an io.Writer and stops writing after the first error.
type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) writef(format string, args ...any) {
	if ew.err != nil {
		return
	}
	_, ew.err = fmt.Fprintf(ew.w, format, args...)
}

func PrintEntries(w io.Writer, entries []Entry, f Format) error {
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	switch f {
	case FormatJSON:
		type jsonEntry struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		out := make([]jsonEntry, len(sorted))
		for i, e := range sorted {
			out[i] = jsonEntry{Key: e.Key, Value: formatValue(e.Value)}
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	default:
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		ew := &errWriter{w: tw}
		ew.writef("KEY\tVALUE\n")
		for _, e := range sorted {
			ew.writef("%s\t%s\n", e.Key, formatValue(e.Value))
		}
		if ew.err != nil {
			return ew.err
		}
		return tw.Flush()
	}
}

func PrintEntry(w io.Writer, entry Entry, f Format) error {
	switch f {
	case FormatJSON:
		type jsonEntry struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(jsonEntry{Key: entry.Key, Value: formatValue(entry.Value)})
	default:
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		ew := &errWriter{w: tw}
		ew.writef("KEY\tVALUE\n")
		ew.writef("%s\t%s\n", entry.Key, formatValue(entry.Value))
		if ew.err != nil {
			return ew.err
		}
		return tw.Flush()
	}
}
