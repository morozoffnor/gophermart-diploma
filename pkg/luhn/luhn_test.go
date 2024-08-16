package luhn

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "[Positive] Check if number is valid",
			value: "12345678903",
			want:  true,
		},
		{
			name:  "[Negative] Check if number is invalid",
			value: "243432433334",
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, Valid(test.value))
		})
	}
}
