package cli

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	tassert "github.com/stretchr/testify/assert"
)

func TestExactArgsWithError(t *testing.T) {
	tests := []struct {
		name          string
		nArgs         int
		args          []string
		expectedError error
	}{
		{
			name:          "exact args without error",
			nArgs:         2,
			args:          []string{"foo", "bar"},
			expectedError: nil,
		},
		{
			name:          "too few args",
			nArgs:         2,
			args:          []string{"foo"},
			expectedError: errors.New("expected 2 args, got 1"),
		},
		{
			name:          "too many args",
			nArgs:         2,
			args:          []string{"foo", "bar", "baz"},
			expectedError: errors.New("expected 2 args, got 3"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			cobraArg := ExactArgsWithError(test.nArgs, test.expectedError)
			err := cobraArg(&cobra.Command{}, test.args)

			if test.expectedError != nil {
				assert.EqualError(err, test.expectedError.Error())
			} else {
				assert.NoError(err)
			}
		})
	}
}
