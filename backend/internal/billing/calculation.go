package billing

// CalculateLineBreakdown implements the authoritative Explainable Billing calculation logic:
// Base Price -> Variant Adjustment -> Quantity -> Discount -> Tax -> Final Total
func CalculateLineBreakdown(
	basePriceCents int64,
	variantAdjustmentCents int64,
	quantity int,
	discountCents int64,
	taxRateBps int,
) ExplainableBreakdown {
	if quantity < 1 {
		quantity = 1
	}

	effectiveUnitPrice := basePriceCents + variantAdjustmentCents
	if effectiveUnitPrice < 0 {
		effectiveUnitPrice = 0
	}

	grossAmount := effectiveUnitPrice * int64(quantity)

	appliedDiscount := discountCents
	if appliedDiscount > grossAmount {
		appliedDiscount = grossAmount
	}
	if appliedDiscount < 0 {
		appliedDiscount = 0
	}

	taxableAmount := grossAmount - appliedDiscount
	taxAmount := (taxableAmount * int64(taxRateBps)) / 10000

	finalTotal := taxableAmount + taxAmount

	return ExplainableBreakdown{
		BasePriceCents:      basePriceCents,
		VariantAdjustment:   variantAdjustmentCents,
		EffectiveUnitPrice:  effectiveUnitPrice,
		Quantity:            quantity,
		GrossAmountCents:    grossAmount,
		DiscountAmountCents: appliedDiscount,
		TaxAmountCents:      taxAmount,
		FinalTotalCents:     finalTotal,
	}
}
