package types

import (
	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSumSlashableBalances_SameAddress(t *testing.T) {
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

	require.Len(t, got, 1)

	require.Equal(t, "0xad45A78180961079BFaeEe349704F411dfF947C6", got[0].Address)
	require.True(t, got[0].Amount.Equal(math.NewInt(150)))
}
