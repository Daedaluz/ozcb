package ozcb

import "strings"

func fixHexString(str string, pad int) string {
	fields := strings.FieldsFunc(str, func(r rune) bool {
		switch r {
		case '_', '.', ',', ';', ':', ' ', '\t', '\r', '\n':
			return true
		}
		return false
	})
	for i, field := range fields {
		fields[i] = strings.TrimPrefix(field, "0x")
	}
	str = strings.Join(fields, "")
	if x := pad - len(str); x > 0 {
		str = strings.Repeat("0", x) + str
	}
	return strings.ToLower(str)
}
