-- Phase 0: Initial Normalized Schema for Veyra (Subflow)
-- Brightloop Software Subscription Operations Platform

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. Users and Authentication
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(32) NOT NULL CHECK (role IN ('ADMIN', 'STAFF', 'CUSTOMER')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Unique index on normalized lowercase email
CREATE UNIQUE INDEX idx_users_email_lower ON users (LOWER(email));

-- 2. Customer Profiles
CREATE TABLE customer_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_name VARCHAR(255),
    phone VARCHAR(50),
    billing_address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_customer_profiles_user UNIQUE (user_id)
);

-- 3. Sessions
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_agent TEXT,
    ip_address VARCHAR(45),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);

-- 4. Password Reset Tokens
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_password_reset_user_id ON password_reset_tokens (user_id);

-- 5. Products & Variants (Member 2 Ownership)
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    product_type VARCHAR(64) NOT NULL DEFAULT 'SUBSCRIPTION_SERVICE',
    sales_price_cents BIGINT NOT NULL CHECK (sales_price_cents >= 0),
    cost_price_cents BIGINT NOT NULL DEFAULT 0 CHECK (cost_price_cents >= 0),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    price_adjustment_cents BIGINT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_variants_product ON product_variants (product_id);

-- 6. Subscription Plans (Member 2 Ownership)
CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    frequency VARCHAR(32) NOT NULL CHECK (frequency IN ('DAILY', 'WEEKLY', 'MONTHLY', 'YEARLY')),
    billing_price_cents BIGINT NOT NULL CHECK (billing_price_cents >= 0),
    min_quantity INT NOT NULL DEFAULT 1 CHECK (min_quantity >= 1),
    available_from DATE NOT NULL,
    available_to DATE,
    auto_close BOOLEAN NOT NULL DEFAULT false,
    allow_early_close BOOLEAN NOT NULL DEFAULT true,
    allow_pause BOOLEAN NOT NULL DEFAULT true,
    allow_renewal BOOLEAN NOT NULL DEFAULT true,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE plan_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    is_addon BOOLEAN NOT NULL DEFAULT false,
    is_required BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_plan_products UNIQUE (plan_id, product_id)
);

-- 7. Quotation Templates & Quotations (Member 4 Ownership)
CREATE TABLE quotation_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    validity_days INT NOT NULL DEFAULT 30 CHECK (validity_days > 0),
    plan_id UUID REFERENCES plans(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE quotation_template_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES quotation_templates(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0)
);

CREATE TABLE quotations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quotation_number VARCHAR(64) NOT NULL UNIQUE,
    customer_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE RESTRICT,
    status VARCHAR(32) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'SENT', 'ACCEPTED', 'EXPIRED', 'CONVERTED')),
    valid_until DATE NOT NULL,
    subtotal_cents BIGINT NOT NULL DEFAULT 0 CHECK (subtotal_cents >= 0),
    tax_cents BIGINT NOT NULL DEFAULT 0 CHECK (tax_cents >= 0),
    discount_cents BIGINT NOT NULL DEFAULT 0 CHECK (discount_cents >= 0),
    total_cents BIGINT NOT NULL DEFAULT 0 CHECK (total_cents >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE quotation_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quotation_id UUID NOT NULL REFERENCES quotations(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0),
    tax_rate_bps INT NOT NULL DEFAULT 1800 CHECK (tax_rate_bps >= 0),
    total_cents BIGINT NOT NULL CHECK (total_cents >= 0)
);

-- 8. Subscriptions (Member 1 / Team Lead Ownership)
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_number VARCHAR(64) NOT NULL UNIQUE,
    customer_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE RESTRICT,
    quotation_id UUID REFERENCES quotations(id) ON DELETE SET NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'QUOTATION', 'CONFIRMED', 'ACTIVE', 'PAUSED', 'CANCELLED', 'CLOSED')),
    start_date DATE NOT NULL,
    expiry_date DATE,
    next_billing_date DATE NOT NULL,
    payment_terms VARCHAR(64) NOT NULL DEFAULT 'NET_30',
    currency VARCHAR(3) NOT NULL DEFAULT 'INR',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_customer ON subscriptions (customer_id);
CREATE INDEX idx_subscriptions_status_billing ON subscriptions (status, next_billing_date);

CREATE TABLE subscription_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,
    is_addon BOOLEAN NOT NULL DEFAULT false,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0),
    tax_rate_bps INT NOT NULL DEFAULT 1800 CHECK (tax_rate_bps >= 0),
    line_total_cents BIGINT NOT NULL CHECK (line_total_cents >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscription_lines_sub ON subscription_lines (subscription_id);

-- Subscription Audit Timeline
CREATE TABLE subscription_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    event_type VARCHAR(64) NOT NULL,
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscription_events_sub_created ON subscription_events (subscription_id, created_at DESC);

-- 9. Tax & Discount Rules (Member 2 Ownership)
CREATE TABLE tax_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    rate_bps INT NOT NULL CHECK (rate_bps >= 0), -- e.g. 1800 = 18.00%
    country VARCHAR(64) NOT NULL DEFAULT 'IN',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE discount_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(64) NOT NULL UNIQUE,
    discount_type VARCHAR(32) NOT NULL CHECK (discount_type IN ('FIXED_AMOUNT', 'PERCENTAGE')),
    value_cents BIGINT NOT NULL CHECK (value_cents > 0),
    min_purchase_cents BIGINT NOT NULL DEFAULT 0 CHECK (min_purchase_cents >= 0),
    min_quantity INT NOT NULL DEFAULT 1 CHECK (min_quantity >= 1),
    max_usage_count INT CHECK (max_usage_count > 0),
    current_usage_count INT NOT NULL DEFAULT 0 CHECK (current_usage_count >= 0),
    applies_to VARCHAR(32) NOT NULL CHECK (applies_to IN ('PRODUCT', 'SUBSCRIPTION')),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    valid_from TIMESTAMPTZ NOT NULL,
    valid_to TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_discount_usage_limit CHECK (max_usage_count IS NULL OR current_usage_count <= max_usage_count)
);

-- 10. Billing Engine & Invoicing (Member 3 Ownership)
CREATE TABLE billing_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_date DATE NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'IN_PROGRESS' CHECK (status IN ('IN_PROGRESS', 'COMPLETED', 'FAILED')),
    subscriptions_evaluated INT NOT NULL DEFAULT 0,
    invoices_generated INT NOT NULL DEFAULT 0,
    failures_count INT NOT NULL DEFAULT 0,
    error_summary TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE TABLE billing_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    billing_run_id UUID NOT NULL REFERENCES billing_runs(id) ON DELETE CASCADE,
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL CHECK (status IN ('SUCCESS', 'FAILED', 'RETRYABLE')),
    error_message TEXT,
    attempt_number INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(64) NOT NULL UNIQUE,
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE RESTRICT,
    customer_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    billing_run_id UUID REFERENCES billing_runs(id) ON DELETE SET NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'CONFIRMED', 'PAID', 'CANCELLED')),
    billing_period_start DATE NOT NULL,
    billing_period_end DATE NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'INR',
    subtotal_cents BIGINT NOT NULL CHECK (subtotal_cents >= 0),
    tax_cents BIGINT NOT NULL DEFAULT 0 CHECK (tax_cents >= 0),
    discount_cents BIGINT NOT NULL DEFAULT 0 CHECK (discount_cents >= 0),
    total_cents BIGINT NOT NULL CHECK (total_cents >= 0),
    due_date DATE NOT NULL,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_sub_billing_period UNIQUE (subscription_id, billing_period_start, billing_period_end)
);

CREATE INDEX idx_invoices_customer ON invoices (customer_id);
CREATE INDEX idx_invoices_status_due ON invoices (status, due_date);

CREATE TABLE invoice_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,
    description VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0),
    tax_cents BIGINT NOT NULL DEFAULT 0 CHECK (tax_cents >= 0),
    discount_cents BIGINT NOT NULL DEFAULT 0 CHECK (discount_cents >= 0),
    total_cents BIGINT NOT NULL CHECK (total_cents >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE discount_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    discount_rule_id UUID NOT NULL REFERENCES discount_rules(id) ON DELETE RESTRICT,
    applied_amount_cents BIGINT NOT NULL CHECK (applied_amount_cents > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 11. Payments (Member 3 Ownership)
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
    payment_number VARCHAR(64) NOT NULL UNIQUE,
    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    payment_method VARCHAR(64) NOT NULL DEFAULT 'BANK_TRANSFER',
    payment_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_invoice ON payments (invoice_id);
