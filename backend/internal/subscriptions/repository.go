package subscriptions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type pgRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{pool: pool}
}

func (r *pgRepository) CreateSubscription(ctx context.Context, sub *Subscription, lines []SubscriptionLine) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO subscriptions (
			subscription_number, customer_id, plan_id, quotation_id,
			status, start_date, expiry_date, next_billing_date, payment_terms, currency,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query,
		sub.SubscriptionNumber, sub.CustomerID, sub.PlanID, sub.QuotationID,
		sub.Status, sub.StartDate, sub.ExpiryDate, sub.NextBillingDate, sub.PaymentTerms, sub.Currency,
	).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		return err
	}

	lineQuery := `
		INSERT INTO subscription_lines (
			subscription_id, product_id, variant_id, is_addon, quantity, unit_price_cents, tax_rate_bps, line_total_cents, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id, created_at
	`
	for i := range lines {
		lines[i].SubscriptionID = sub.ID
		err = tx.QueryRow(ctx, lineQuery,
			sub.ID, lines[i].ProductID, lines[i].VariantID, lines[i].IsAddon,
			lines[i].Quantity, lines[i].UnitPriceCents, lines[i].TaxRateBps, lines[i].LineTotalCents,
		).Scan(&lines[i].ID, &lines[i].CreatedAt)
		if err != nil {
			return fmt.Errorf("failed inserting subscription line %d: %w", i, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *pgRepository) GetSubscriptionByID(ctx context.Context, id string) (*Subscription, []SubscriptionLine, error) {
	query := `
		SELECT s.id, s.subscription_number, s.customer_id, u.email, s.plan_id, p.name, s.quotation_id,
		       s.status, s.start_date, s.expiry_date, s.next_billing_date, s.payment_terms, s.currency,
		       s.created_at, s.updated_at
		FROM subscriptions s
		JOIN users u ON s.customer_id = u.id
		JOIN plans p ON s.plan_id = p.id
		WHERE s.id = $1
	`
	var sub Subscription
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.SubscriptionNumber, &sub.CustomerID, &sub.CustomerEmail, &sub.PlanID, &sub.PlanName, &sub.QuotationID,
		&sub.Status, &sub.StartDate, &sub.ExpiryDate, &sub.NextBillingDate, &sub.PaymentTerms, &sub.Currency,
		&sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ErrSubscriptionNotFound
		}
		return nil, nil, err
	}

	lineQuery := `
		SELECT sl.id, sl.subscription_id, sl.product_id, pr.name, sl.variant_id, COALESCE(pv.name, ''),
		       sl.is_addon, sl.quantity, sl.unit_price_cents, sl.tax_rate_bps, sl.line_total_cents, sl.created_at
		FROM subscription_lines sl
		JOIN products pr ON sl.product_id = pr.id
		LEFT JOIN product_variants pv ON sl.variant_id = pv.id
		WHERE sl.subscription_id = $1
	`
	rows, err := r.pool.Query(ctx, lineQuery, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	lines := make([]SubscriptionLine, 0)
	var totalCents int64 = 0
	for rows.Next() {
		var l SubscriptionLine
		var variantName string
		if err := rows.Scan(
			&l.ID, &l.SubscriptionID, &l.ProductID, &l.ProductName, &l.VariantID, &variantName,
			&l.IsAddon, &l.Quantity, &l.UnitPriceCents, &l.TaxRateBps, &l.LineTotalCents, &l.CreatedAt,
		); err != nil {
			return nil, nil, err
		}
		if variantName != "" {
			l.VariantName = &variantName
		}
		totalCents += l.LineTotalCents
		lines = append(lines, l)
	}

	sub.Lines = lines
	sub.TotalAmountCents = totalCents
	return &sub, lines, nil
}

func (r *pgRepository) ListSubscriptions(ctx context.Context, customerID *string, status *string) ([]Subscription, error) {
	query := `
		SELECT s.id, s.subscription_number, s.customer_id, u.email, s.plan_id, p.name, s.quotation_id,
		       s.status, s.start_date, s.expiry_date, s.next_billing_date, s.payment_terms, s.currency,
		       s.created_at, s.updated_at,
		       COALESCE(SUM(sl.line_total_cents), 0) as total_amount_cents
		FROM subscriptions s
		JOIN users u ON s.customer_id = u.id
		JOIN plans p ON s.plan_id = p.id
		LEFT JOIN subscription_lines sl ON s.id = sl.subscription_id
		WHERE ($1::uuid IS NULL OR s.customer_id = $1::uuid)
		  AND ($2::varchar IS NULL OR s.status = $2::varchar)
		GROUP BY s.id, u.email, p.name
		ORDER BY s.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, customerID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subs := make([]Subscription, 0)
	for rows.Next() {
		var s Subscription
		if err := rows.Scan(
			&s.ID, &s.SubscriptionNumber, &s.CustomerID, &s.CustomerEmail, &s.PlanID, &s.PlanName, &s.QuotationID,
			&s.Status, &s.StartDate, &s.ExpiryDate, &s.NextBillingDate, &s.PaymentTerms, &s.Currency,
			&s.CreatedAt, &s.UpdatedAt, &s.TotalAmountCents,
		); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}

func (r *pgRepository) GetSubscriptionsDueForBilling(ctx context.Context, cutoffDate time.Time, limit int) ([]Subscription, error) {
	query := `
		SELECT s.id, s.subscription_number, s.customer_id, u.email, s.plan_id, p.name, s.quotation_id,
		       s.status, s.start_date, s.expiry_date, s.next_billing_date, s.payment_terms, s.currency,
		       s.created_at, s.updated_at
		FROM subscriptions s
		JOIN users u ON s.customer_id = u.id
		JOIN plans p ON s.plan_id = p.id
		WHERE s.status = 'ACTIVE' AND s.next_billing_date <= $1
		ORDER BY s.next_billing_date ASC
		LIMIT $2
	`
	rows, err := r.pool.Query(ctx, query, cutoffDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subs := make([]Subscription, 0)
	for rows.Next() {
		var s Subscription
		if err := rows.Scan(
			&s.ID, &s.SubscriptionNumber, &s.CustomerID, &s.CustomerEmail, &s.PlanID, &s.PlanName, &s.QuotationID,
			&s.Status, &s.StartDate, &s.ExpiryDate, &s.NextBillingDate, &s.PaymentTerms, &s.Currency,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	for i := range subs {
		_, lines, err := r.GetSubscriptionByID(ctx, subs[i].ID)
		if err == nil {
			subs[i].Lines = lines
		}
	}

	return subs, nil
}

func (r *pgRepository) UpdateStatus(ctx context.Context, id string, status SubscriptionStatus) error {
	query := `UPDATE subscriptions SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

func (r *pgRepository) UpdateNextBillingDate(ctx context.Context, id string, nextBillingDate time.Time) error {
	query := `UPDATE subscriptions SET next_billing_date = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, nextBillingDate, id)
	return err
}

func (r *pgRepository) UpdateAddonStatus(ctx context.Context, subscriptionID, productID string, active bool) error {
	if active {
		query := `UPDATE subscription_lines SET is_addon = true WHERE subscription_id = $1 AND product_id = $2`
		_, err := r.pool.Exec(ctx, query, subscriptionID, productID)
		return err
	}
	query := `DELETE FROM subscription_lines WHERE subscription_id = $1 AND product_id = $2 AND is_addon = true`
	_, err := r.pool.Exec(ctx, query, subscriptionID, productID)
	return err
}

func (r *pgRepository) RecordEvent(ctx context.Context, event *SubscriptionEvent) error {
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		payloadBytes = []byte("{}")
	}

	query := `
		INSERT INTO subscription_events (subscription_id, event_type, actor_id, payload, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query, event.SubscriptionID, event.EventType, event.ActorID, payloadBytes).
		Scan(&event.ID, &event.CreatedAt)
}

func (r *pgRepository) ListEvents(ctx context.Context, subscriptionID string) ([]SubscriptionEvent, error) {
	query := `
		SELECT se.id, se.subscription_id, se.event_type, se.actor_id, COALESCE(u.email, 'SYSTEM'), se.payload, se.created_at
		FROM subscription_events se
		LEFT JOIN users u ON se.actor_id = u.id
		WHERE se.subscription_id = $1
		ORDER BY se.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]SubscriptionEvent, 0)
	for rows.Next() {
		var e SubscriptionEvent
		var rawPayload []byte
		if err := rows.Scan(&e.ID, &e.SubscriptionID, &e.EventType, &e.ActorID, &e.ActorEmail, &rawPayload, &e.CreatedAt); err != nil {
			return nil, err
		}
		var payloadMap map[string]any
		_ = json.Unmarshal(rawPayload, &payloadMap)
		e.Payload = payloadMap
		events = append(events, e)
	}
	return events, nil
}
