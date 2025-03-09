package types

import (
	"cosmossdk.io/math"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

// Template for testing improved sumTokenBalances function
func TestSumTokenBalances_Template(t *testing.T) {
	testCases := []struct {
		name     string
		tb1      []TokenBalance
		tb2      []TokenBalance
		expected []TokenBalance
	}{
		{
			name:     "Both slices empty",
			tb1:      []TokenBalance{},
			tb2:      []TokenBalance{},
			expected: []TokenBalance{},
		},
		{
			name: "First slice empty",
			tb1:  []TokenBalance{},
			tb2: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(50),
				},
			},
			expected: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(50),
				},
			},
		},
		{
			name: "Second slice empty",
			tb1: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(100),
				},
			},
			tb2: []TokenBalance{},
			expected: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(100),
				},
			},
		},
		{
			name: "No matching addresses",
			tb1: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(100),
				},
			},
			tb2: []TokenBalance{
				{
					Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
					Amount:  math.NewInt(50),
				},
			},
			expected: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(100),
				},
				{
					Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
					Amount:  math.NewInt(50),
				},
			},
		},
		{
			name: "Single address match",
			tb1: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(100),
				},
			},
			tb2: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(50),
				},
				{
					Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
					Amount:  math.NewInt(50),
				},
			},
			expected: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(150),
				},
				{
					Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
					Amount:  math.NewInt(50),
				},
			},
		},
		{
			name: "Multiple matches",
			tb1: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(100),
				},
				{
					Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
					Amount:  math.NewInt(200),
				},
			},
			tb2: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(50),
				},
				{
					Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
					Amount:  math.NewInt(50),
				},
				{
					Address: "0x1111111111111111111111111111111111111111",
					Amount:  math.NewInt(75),
				},
			},
			expected: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(150),
				},
				{
					Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
					Amount:  math.NewInt(250),
				},
				{
					Address: "0x1111111111111111111111111111111111111111",
					Amount:  math.NewInt(75),
				},
			},
		},
		{
			name: "Zero amounts",
			tb1: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(0),
				},
			},
			tb2: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(0),
				},
			},
			expected: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(0),
				},
			},
		},
		{
			name: "Negative amounts",
			tb1: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(100),
				},
			},
			tb2: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(-50),
				},
			},
			expected: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(50),
				},
			},
		},
		{
			name: "Mix of matched and unmatched addresses",
			tb1: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(100),
				},
				{
					Address: "0x2222222222222222222222222222222222222222",
					Amount:  math.NewInt(300),
				},
			},
			tb2: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(50),
				},
				{
					Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
					Amount:  math.NewInt(150),
				},
				{
					Address: "0x3333333333333333333333333333333333333333",
					Amount:  math.NewInt(200),
				},
			},
			expected: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(150),
				},
				{
					Address: "0x2222222222222222222222222222222222222222",
					Amount:  math.NewInt(300),
				},
				{
					Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
					Amount:  math.NewInt(150),
				},
				{
					Address: "0x3333333333333333333333333333333333333333",
					Amount:  math.NewInt(200),
				},
			},
		},
		{
			name: "Multiple tokens from both lists",
			tb1: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(100),
				},
				{
					Address: "0x2222222222222222222222222222222222222222",
					Amount:  math.NewInt(300),
				},
				{
					Address: "0x4444444444444444444444444444444444444444",
					Amount:  math.NewInt(400),
				},
			},
			tb2: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(50),
				},
				{
					Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
					Amount:  math.NewInt(150),
				},
				{
					Address: "0x4444444444444444444444444444444444444444",
					Amount:  math.NewInt(75),
				},
			},
			expected: []TokenBalance{
				{
					Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
					Amount:  math.NewInt(150),
				},
				{
					Address: "0x2222222222222222222222222222222222222222",
					Amount:  math.NewInt(300),
				},
				{
					Address: "0x4444444444444444444444444444444444444444",
					Amount:  math.NewInt(475),
				},
				{
					Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
					Amount:  math.NewInt(150),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := sumTokenBalances(tc.tb1, tc.tb2)

			require.Len(t, got, len(tc.expected), "Result length doesn't match expected for case: "+tc.name)

			// For each expected item, find a matching result item
			for _, expectedItem := range tc.expected {
				found := false
				for _, gotItem := range got {
					if gotItem.Address == expectedItem.Address && gotItem.Amount.Equal(expectedItem.Amount) {
						found = true
						break
					}
				}
				require.True(t, found, "Expected item not found in result: %v for case: %s", expectedItem, tc.name)
			}
		})
	}
}

func TestSumTokenBalances_SameAddress(t *testing.T) {
	tb1 := []TokenBalance{
		{
			Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
			Amount:  math.NewInt(100),
		},
	}
	tb2 := []TokenBalance{
		{
			Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
			Amount:  math.NewInt(50),
		},
		{
			Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Amount:  math.NewInt(50),
		},
	}

	got := sumTokenBalances(tb1, tb2)

	// Should now have 2 entries (one matched address + one unmatched from tb2)
	require.Len(t, got, 2)

	// Find the entry with the first address and check it
	found := false
	for _, item := range got {
		if item.Address == "0xad45A78180961079BFaeEe349704F411dfF947C6" {
			require.True(t, item.Amount.Equal(math.NewInt(150)))
			found = true
			break
		}
	}
	require.True(t, found, "Expected to find matched address with summed amount")

	// Find the entry with the second address and check it
	found = false
	for _, item := range got {
		if item.Address == "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE" {
			require.True(t, item.Amount.Equal(math.NewInt(50)))
			found = true
			break
		}
	}
	require.True(t, found, "Expected to find unmatched address from tb2")

	fmt.Println(got)
}

func TestSumTokenBalances_Empty(t *testing.T) {
	// Test with both slices empty
	tb1 := []TokenBalance{}
	tb2 := []TokenBalance{}

	got := sumTokenBalances(tb1, tb2)
	require.Len(t, got, 0, "Expected empty result for empty inputs")

	// Test with only first slice empty
	tb2 = []TokenBalance{
		{
			Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
			Amount:  math.NewInt(50),
		},
	}

	got = sumTokenBalances(tb1, tb2)
	require.Len(t, got, 1, "Expected one token in result when first slice is empty")
	require.Equal(t, "0xad45A78180961079BFaeEe349704F411dfF947C6", got[0].Address)
	require.True(t, got[0].Amount.Equal(math.NewInt(50)))

	// Test with only second slice empty
	tb1 = []TokenBalance{
		{
			Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
			Amount:  math.NewInt(100),
		},
	}
	tb2 = []TokenBalance{}

	got = sumTokenBalances(tb1, tb2)
	require.Len(t, got, 1, "Expected one token in result when second slice is empty")
	require.Equal(t, "0xad45A78180961079BFaeEe349704F411dfF947C6", got[0].Address)
	require.True(t, got[0].Amount.Equal(math.NewInt(100)))
}

func TestSumTokenBalances_NoMatches(t *testing.T) {
	tb1 := []TokenBalance{
		{
			Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
			Amount:  math.NewInt(100),
		},
	}
	tb2 := []TokenBalance{
		{
			Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Amount:  math.NewInt(50),
		},
	}

	got := sumTokenBalances(tb1, tb2)
	require.Len(t, got, 2, "Expected both tokens in result when no addresses match")

	// Find and verify first token
	found := false
	for _, item := range got {
		if item.Address == "0xad45A78180961079BFaeEe349704F411dfF947C6" {
			require.True(t, item.Amount.Equal(math.NewInt(100)))
			found = true
			break
		}
	}
	require.True(t, found, "Expected to find token from tb1")

	// Find and verify second token
	found = false
	for _, item := range got {
		if item.Address == "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE" {
			require.True(t, item.Amount.Equal(math.NewInt(50)))
			found = true
			break
		}
	}
	require.True(t, found, "Expected to find token from tb2")
}

func TestSumTokenBalances_MultipleMatches(t *testing.T) {
	tb1 := []TokenBalance{
		{
			Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
			Amount:  math.NewInt(100),
		},
		{
			Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Amount:  math.NewInt(200),
		},
		{
			Address: "0x1111111111111111111111111111111111111111",
			Amount:  math.NewInt(300),
		},
	}
	tb2 := []TokenBalance{
		{
			Address: "0xad45A78180961079BFaeEe349704F411dfF947C6",
			Amount:  math.NewInt(50),
		},
		{
			Address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Amount:  math.NewInt(50),
		},
		{
			Address: "0x2222222222222222222222222222222222222222",
			Amount:  math.NewInt(500),
		},
	}

	got := sumTokenBalances(tb1, tb2)
	require.Len(t, got, 4, "Expected 4 tokens in result: 2 matched + 1 unique from tb1 + 1 unique from tb2")

	// Verify each token is present with correct amounts
	expectedResults := map[string]int64{
		"0xad45A78180961079BFaeEe349704F411dfF947C6": 150,
		"0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE": 250,
		"0x1111111111111111111111111111111111111111": 300,
		"0x2222222222222222222222222222222222222222": 500,
	}

	for expectedAddr, expectedAmount := range expectedResults {
		found := false
		for _, item := range got {
			if item.Address == expectedAddr {
				require.True(t, item.Amount.Equal(math.NewInt(expectedAmount)),
					"Expected amount %d for address %s, got %s", expectedAmount, expectedAddr, item.Amount)
				found = true
				break
			}
		}
		require.True(t, found, "Expected to find address %s in result", expectedAddr)
	}
}
