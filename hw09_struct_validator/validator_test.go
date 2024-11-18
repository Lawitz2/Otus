package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,staff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	BadRulesAndTypes struct {
		Code  int     `validate:"badRule:500"`
		State float64 `validate:"min:10"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "testname",
				Age:    25,
				Email:  "test@test.com",
				Role:   "staff",
				Phones: []string{"89991234567"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			in:          App{Version: "1.0.0"},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 404,
				Body: "something",
			},
			expectedErr: nil,
		},
		{
			in: App{Version: "1.0.0.0.0"},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Version",
				Err:   errLenMismatch,
			}},
		},
		{
			in: Response{
				Code: 123,
				Body: "something",
			},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Code",
				Err:   errIntNotInSet,
			}},
		},
		{
			in: User{
				ID:     "badID",
				Name:   "notValidated",
				Age:    16,
				Email:  "badEmail",
				Role:   "badRole",
				Phones: []string{"12", "34", "56"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{ValidationError{
				Field: "ID",
				Err:   errLenMismatch,
			}, ValidationError{
				Field: "Age",
				Err:   errBelowMin,
			}, ValidationError{
				Field: "Email",
				Err:   errRegexMismatch,
			}, ValidationError{
				Field: "Role",
				Err:   errStrNotInSet,
			}, ValidationError{ // each slice element that didn't satisfy the validation rule
				Field: "Phones", // will add another error, hence three errLenMismatch
				Err:   errLenMismatch,
			}, ValidationError{
				Field: "Phones",
				Err:   errLenMismatch,
			}, ValidationError{
				Field: "Phones",
				Err:   errLenMismatch,
			}},
		},
		{
			in: Token{
				Header:    []byte{},
				Payload:   []byte{25, 255},
				Signature: []byte{0},
			},
			expectedErr: nil,
		},
		{
			in: BadRulesAndTypes{
				Code:  256,
				State: 256,
			},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Code",
				Err:   errUnsuppRule,
			}, ValidationError{
				Field: "State",
				Err:   errUnsuppType,
			}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			require.Equal(t, Validate(tt.in), tt.expectedErr)
			_ = tt
		})
	}
}
