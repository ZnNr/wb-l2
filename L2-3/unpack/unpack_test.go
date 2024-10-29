package unpack

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{`qwe\4\5`, "qwe45", false},
		{`qwe\45`, "qwe44444", false},
		{`qwe\\5`, "qwe\\\\\\\\\\", false},
		{"\\", "", true},
		{"a\\", "", true},
		{"a\\2", "a2", false},
		{"a0", "", false},
	}

	for _, test := range tests {
		result, err := Unpack(test.input)
		if test.hasError {
			if err == nil {
				t.Errorf("Ожидалась ошибка для входа %q, но ошибки не было", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Неожиданная ошибка для входа %q: %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("Для входа %q ожидается %q, получено %q", test.input, test.expected, result)
			}
		}
	}
}
