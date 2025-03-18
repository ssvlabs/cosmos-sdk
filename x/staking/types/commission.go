package types

import (
	"time"

	"cosmossdk.io/math"
)

// NewCommissionRates returns an initialized validator commission rates.
func NewCommissionRates(rate math.LegacyDec) CommissionRates {
	return CommissionRates{
		Rate: rate,
	}
}

// NewCommission returns an initialized validator commission.
func NewCommission(rate, maxRate, maxChangeRate math.LegacyDec) Commission {
	return Commission{
		CommissionRates: NewCommissionRates(rate),
		UpdateTime:      time.Unix(0, 0).UTC(),
	}
}

// NewCommissionWithTime returns an initialized validator commission with a specified
// update time which should be the current block BFT time.
func NewCommissionWithTime(rate math.LegacyDec, updatedAt time.Time) Commission {
	return Commission{
		CommissionRates: NewCommissionRates(rate),
		UpdateTime:      updatedAt,
	}
}

// Validate performs basic sanity validation checks of initial commission
// parameters. If validation fails, an SDK error is returned.
func (cr CommissionRates) Validate() error {
	switch {
	case cr.Rate.IsNegative():
		// rate cannot be negative
		return ErrCommissionNegative

	}
	return nil
}

// ValidateNewRate performs basic sanity validation checks of a new commission
// rate. If validation fails, an SDK error is returned.
func (c Commission) ValidateNewRate(newRate math.LegacyDec, blockTime time.Time) error {
	switch {
	case blockTime.Sub(c.UpdateTime).Hours() < 24:
		// new rate cannot be changed more than once within 24 hours
		return ErrCommissionUpdateTime

	case newRate.IsNegative():
		// new rate cannot be negative
		return ErrCommissionNegative
	}

	return nil
}
