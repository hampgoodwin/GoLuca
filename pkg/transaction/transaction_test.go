package transaction

import (
	"testing"
)

func TestBalanced(t *testing.T) {
	testCases := []struct {
		input    Transaction
		expected bool
	}{
		{
			Transaction{
				Description: "balanced one Entry",
				Entries: []Entry{
					{
						Amount: float64(0),
					},
				},
			},
			true,
		},
		{
			Transaction{
				Description: "balanced two Entries",
				Entries: []Entry{
					{
						Amount: float64(2),
					},
					{
						Amount: float64(-2),
					},
				},
			},
			true,
		},
		{
			Transaction{
				Description: "balanced three Entries",
				Entries: []Entry{
					{
						Amount: float64(2),
					},
					{
						Amount: float64(-1),
					},
					{
						Amount: float64(-1),
					},
				},
			},
			true,
		},
		{
			Transaction{
				Description: "imbalanced one Entry",
				Entries: []Entry{
					{
						Amount: float64(1),
					},
				},
			},
			false,
		},
		{
			Transaction{
				Description: "imbalanced two Entries",
				Entries: []Entry{
					{
						Amount: float64(1),
					},
					{
						Amount: float64(-2),
					},
				},
			},
			false,
		},
		{
			Transaction{
				Description: "imbalanced three Entries",
				Entries: []Entry{
					{
						Amount: float64(1),
					},
					{
						Amount: float64(-2),
					},
					{
						Amount: float64(2),
					},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.input.Description, func(t *testing.T) {
			t.Parallel()
			output := tc.input.Balanced()
			if tc.expected != output {
				t.Errorf("expected %t\nreceived %t\n", output, tc.expected)
			}
		})
	}
}
