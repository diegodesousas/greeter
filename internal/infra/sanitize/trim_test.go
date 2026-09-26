package sanitize_test

import (
	"testing"

	"github.com/diegodesousas/greeter/internal/infra/sanitize"
	"github.com/stretchr/testify/assert"
)

type status string

type address struct {
	Street string
}

type payload struct {
	Name       string
	Password   string `trim:"-"`
	Status     status
	Age        int
	Nickname   *string
	Tags       []string
	Codes      [2]string
	Labels     map[string]string
	Extra      map[string]any
	Address    address
	AddressPtr *address
	Addresses  []address
	Skipped    address `trim:"-"`
	address
	unexported string
}

func strPtr(s string) *string { return &s }

func TestTrimStrings(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  any
	}{
		{
			name:  "trims a string pointer",
			input: strPtr("  Diego \t\n"),
			want:  strPtr("Diego"),
		},
		{
			name: "trims every string reachable from a struct",
			input: &payload{
				Name:       "  Diego  ",
				Password:   "  secret  ",
				Status:     " active ",
				Age:        30,
				Nickname:   strPtr(" Di "),
				Tags:       []string{" a ", "b "},
				Codes:      [2]string{" x", "y "},
				Labels:     map[string]string{" key ": " value "},
				Extra:      map[string]any{"k": " v ", "n": 1},
				Address:    address{Street: " Rua A "},
				AddressPtr: &address{Street: " Rua B "},
				Addresses:  []address{{Street: " Rua C "}},
				Skipped:    address{Street: " Rua D "},
				address:    address{Street: " Rua E "},
				unexported: " hidden ",
			},
			want: &payload{
				Name:       "Diego",
				Password:   "  secret  ",
				Status:     "active",
				Age:        30,
				Nickname:   strPtr("Di"),
				Tags:       []string{"a", "b"},
				Codes:      [2]string{"x", "y"},
				Labels:     map[string]string{" key ": "value"},
				Extra:      map[string]any{"k": "v", "n": 1},
				Address:    address{Street: "Rua A"},
				AddressPtr: &address{Street: "Rua B"},
				Addresses:  []address{{Street: "Rua C"}},
				Skipped:    address{Street: " Rua D "},
				address:    address{Street: " Rua E "},
				unexported: " hidden ",
			},
		},
		{
			name:  "nil pointers are left untouched",
			input: &payload{},
			want:  &payload{},
		},
		{
			name:  "nil input does not panic",
			input: nil,
			want:  nil,
		},
		{
			name:  "non-pointer input is ignored",
			input: payload{Name: " Diego "},
			want:  payload{Name: " Diego "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sanitize.TrimStrings(tt.input)

			assert.Equal(t, tt.want, tt.input)
		})
	}
}
