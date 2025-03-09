package types

import (
	"sort"
)

func (c Capital) IsPositive() bool {
	return !c.NonSlashableCapital.IsZero() || len(c.SlashableBalance) > 0
}

func (c Capital) IsZero() bool {
	return c.NonSlashableCapital.IsZero() && len(c.SlashableBalance) == 0
}

// Assume that non-slashable capital tokens matches validator capital
func (c Capital) Add(addCapital Capital) Capital {
	return Capital{
		SlashableBalance:    sumTokenBalances(c.GetSlashableBalance(), addCapital.SlashableBalance),
		NonSlashableCapital: c.NonSlashableCapital.Add(addCapital.NonSlashableCapital),
	}
}

func (c Capital) Sub(addCapital Capital) Capital {
	return Capital{
		SlashableBalance:    subtractTokenBalances(c.GetSlashableBalance(), addCapital.SlashableBalance),
		NonSlashableCapital: c.NonSlashableCapital.Sub(addCapital.NonSlashableCapital),
	}
}

func (c Capital) HasEnoughFunds(required Capital) bool {
	if c.NonSlashableCapital.LT(required.NonSlashableCapital) {
		return false
	}

	for _, requiredtoken := range required.GetSlashableBalance() {
		found := false
		for _, balance := range c.GetSlashableBalance() {
			if requiredtoken.Address == balance.Address {
				found = true
				if requiredtoken.Amount.GT(balance.Amount) {
					return false
				}
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (c Capital) Equal(other Capital) bool {
	if !c.NonSlashableCapital.Equal(other.NonSlashableCapital) {
		return false
	}

	if len(c.SlashableBalance) != len(other.SlashableBalance) {
		return false
	}

	cBalances := make([]TokenBalance, len(c.SlashableBalance))
	copy(cBalances, c.SlashableBalance)

	otherBalances := make([]TokenBalance, len(other.SlashableBalance))
	copy(otherBalances, other.SlashableBalance)

	sort.Slice(cBalances, func(i, j int) bool {
		return cBalances[i].Address < cBalances[j].Address
	})
	sort.Slice(otherBalances, func(i, j int) bool {
		return otherBalances[i].Address < otherBalances[j].Address
	})

	for i := range cBalances {
		if cBalances[i].Address != otherBalances[i].Address {
			return false
		}

		if !cBalances[i].Amount.Equal(otherBalances[i].Amount) {
			return false
		}
	}

	return true
}

func sumTokenBalances(tb1, tb2 []TokenBalance) []TokenBalance {
	balanceMap := make(map[string]TokenBalance)

	for _, balance := range tb1 {
		balanceMap[balance.Address] = TokenBalance{
			Address: balance.Address,
			Amount:  balance.Amount,
		}
	}

	// Process second slice, adding to existing balances or creating new entries
	for _, balance := range tb2 {
		if existingBalance, found := balanceMap[balance.Address]; found {
			balanceMap[balance.Address] = TokenBalance{
				Address: balance.Address,
				Amount:  existingBalance.Amount.Add(balance.Amount),
			}
		} else {
			balanceMap[balance.Address] = TokenBalance{
				Address: balance.Address,
				Amount:  balance.Amount,
			}
		}
	}

	result := make([]TokenBalance, 0, len(balanceMap))
	for _, balance := range balanceMap {
		result = append(result, balance)
	}

	return result
}

func subtractTokenBalances(tb1, tb2 []TokenBalance) []TokenBalance {
	result := make([]TokenBalance, 0, len(tb1))

	for _, ac := range tb1 {
		newAmount := ac.Amount
		for _, sc := range tb2 {
			if ac.Address == sc.Address {
				newAmount = newAmount.Sub(sc.Amount)
				break
			}
		}

		result = append(result, TokenBalance{
			Address: ac.Address,
			Amount:  newAmount,
		})
	}

	return result
}
