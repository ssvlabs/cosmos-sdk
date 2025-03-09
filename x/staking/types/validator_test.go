package types

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestPotentialConsensusPower(t *testing.T) {
	powerReduction := math.NewIntFromUint64(1000000)

	// Test with empty slashable balance
	t.Run("empty slashable balance", func(t *testing.T) {
		validator := Validator{
			Capital: Capital{
				SlashableBalance:    []TokenBalance{},
				NonSlashableCapital: math.NewInt(1000),
			},
		}

		require.EqualValues(t, 0, validator.PotentialConsensusPower(powerReduction))
	})

	// Test with one slashable balance
	t.Run("single slashable balance", func(t *testing.T) {
		validator := Validator{
			Capital: Capital{
				SlashableBalance: []TokenBalance{
					{Amount: math.NewInt(1000000)},
				},
				NonSlashableCapital: math.NewInt(1000),
			},
		}

		require.EqualValues(t, 1, validator.PotentialConsensusPower(powerReduction))
	})

	// Test with multiple equal slashable balances
	t.Run("multiple equal slashable balances", func(t *testing.T) {
		validator := Validator{
			Capital: Capital{
				SlashableBalance: []TokenBalance{
					{Amount: math.NewInt(1000000000)},
					{Amount: math.NewInt(1000000000)},
					{Amount: math.NewInt(1000000000)},
				},
				NonSlashableCapital: math.NewInt(1000),
			},
		}

		require.EqualValues(t, 1000, validator.PotentialConsensusPower(powerReduction))
	})

	t.Run("multiple different slashable balances", func(t *testing.T) {
		validator := Validator{
			Capital: Capital{
				SlashableBalance: []TokenBalance{
					{Amount: math.NewInt(50000000000000)},
					{Amount: MustCreateNewInt("12500000000000000000")},
				},
				NonSlashableCapital: math.NewInt(1000),
			},
		}

		require.EqualValues(t, 99999600, validator.PotentialConsensusPower(powerReduction))
	})

	t.Run("multiple slashable balances but one is zero", func(t *testing.T) {
		validator := Validator{
			Capital: Capital{
				SlashableBalance: []TokenBalance{
					{Amount: math.NewInt(0)},
					{Amount: MustCreateNewInt("75000000000000000")},
				},
				NonSlashableCapital: math.NewInt(1000),
			},
		}

		require.EqualValues(t, 75000000000, validator.PotentialConsensusPower(powerReduction))
	})

	t.Run("ignore zero amounts", func(t *testing.T) {
		validator := Validator{
			Capital: Capital{
				SlashableBalance: []TokenBalance{
					{Amount: math.NewInt(0)},
					{Amount: math.NewInt(2000000)},
					{Amount: math.NewInt(0)},
				},
				NonSlashableCapital: math.NewInt(1000),
			},
		}

		require.EqualValues(t, 2, validator.PotentialConsensusPower(powerReduction))
	})

	// Test with some negative amounts (should be ignored)
	t.Run("ignore negative amounts", func(t *testing.T) {
		validator := Validator{
			Capital: Capital{
				SlashableBalance: []TokenBalance{
					{Amount: math.NewInt(1000000)},
					{Amount: math.NewInt(-2000000)}, // Should be ignored
				},
				NonSlashableCapital: math.NewInt(1000),
			},
		}

		power := validator.PotentialConsensusPower(powerReduction)
		require.EqualValues(t, 1, power)
	})

	// Test with all zero/negative amounts
	t.Run("all zero or negative amounts", func(t *testing.T) {
		validator := Validator{
			Capital: Capital{
				SlashableBalance: []TokenBalance{
					{Amount: math.NewInt(0)},
					{Amount: math.NewInt(-1000)},
				},
				NonSlashableCapital: math.NewInt(1000),
			},
		}

		require.EqualValues(t, 0, validator.PotentialConsensusPower(powerReduction))
	})

	// Test with large numbers
	t.Run("large numbers", func(t *testing.T) {
		largeNumber := new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil) // 10^20
		validator := Validator{
			Capital: Capital{
				SlashableBalance: []TokenBalance{
					{Amount: math.NewIntFromBigInt(largeNumber)},
				},
				NonSlashableCapital: math.NewInt(1000),
			},
		}

		power := validator.PotentialConsensusPower(powerReduction)
		expected := new(big.Int).Div(largeNumber, big.NewInt(1000000)).Int64()
		require.EqualValues(t, expected, power)
	})

	// Test with typical staking scenario
	t.Run("typical staking scenario", func(t *testing.T) {
		validator := Validator{
			Capital: Capital{
				SlashableBalance: []TokenBalance{
					{Amount: math.NewInt(10000000)}, // 10 tokens with powerReduction of 1000000
					{Amount: math.NewInt(20000000)}, // 20 tokens
					{Amount: math.NewInt(15000000)}, // 15 tokens
				},
				NonSlashableCapital: math.NewInt(5000000), // Not used in calculation
			},
		}

		// Harmonic mean calculation:
		// count = 3
		// sum of inverses = 1/10000000 + 1/20000000 + 1/15000000 = 0.0000001 + 0.00000005 + 0.00000006667 = 0.00000021667
		// harmonic mean = 3 / 0.00000021667 = 13847343.43
		// divide by powerReduction (1000000) = 13.85
		// floor to int64 = 13

		require.EqualValues(t, 13, validator.PotentialConsensusPower(powerReduction))
	})
}

func MustCreateNewInt(bigIntStr string) math.Int {
	v, _ := math.NewIntFromString(bigIntStr)
	return v
}
