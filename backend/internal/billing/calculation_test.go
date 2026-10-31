package billing

import (
	"testing"
)

func TestCalculateLineBreakdown(t *testing.T) {
	// Reference test: Brightloop Core (499900 cents) + Odoo variant (56000 cents)
	// Quantity: 2
	// Discount: 50000 cents
	// Tax: 1800 bps (18%)
	basePrice := int64(499900)
	variantAdj := int64(56000)
	qty := 2
	discount := int64(50000)
	taxRateBps := 1800

	breakdown := CalculateLineBreakdown(basePrice, variantAdj, qty, discount, taxRateBps)

	expectedEffective := int64(555900)
	if breakdown.EffectiveUnitPrice != expectedEffective {
		t.Errorf("expected effective unit price %d, got %d", expectedEffective, breakdown.EffectiveUnitPrice)
	}

	expectedGross := expectedEffective * int64(qty) // 1111800
	if breakdown.GrossAmountCents != expectedGross {
		t.Errorf("expected gross %d, got %d", expectedGross, breakdown.GrossAmountCents)
	}

	expectedTaxable := expectedGross - discount     // 1061800
	expectedTax := (expectedTaxable * 1800) / 10000 // 191124
	if breakdown.TaxAmountCents != expectedTax {
		t.Errorf("expected tax %d, got %d", expectedTax, breakdown.TaxAmountCents)
	}

	expectedTotal := expectedTaxable + expectedTax // 1252924
	if breakdown.FinalTotalCents != expectedTotal {
		t.Errorf("expected final total %d, got %d", expectedTotal, breakdown.FinalTotalCents)
	}
}
