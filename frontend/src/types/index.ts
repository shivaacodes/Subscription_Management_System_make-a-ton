export type Role = 'ADMIN' | 'STAFF' | 'CUSTOMER';

export interface User {
  id: string;
  email: string;
  role: Role;
  createdAt: string;
  updatedAt: string;
}

export interface CustomerProfile {
  id: string;
  userId: string;
  companyName?: string;
  phone?: string;
  billingAddress?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ApiErrorResponse {
  error: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}

export interface AuthResponse {
  user: User;
  message?: string;
}

export interface HealthResponse {
  status: string;
  timestamp: string;
  version: string;
  database: string;
}

// Catalog
export interface ProductVariant {
  id: string;
  productId: string;
  name: string;
  priceAdjustmentCents: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Product {
  id: string;
  name: string;
  productType: string;
  salesPriceCents: number;
  costPriceCents: number;
  notes?: string;
  isActive: boolean;
  variants?: ProductVariant[];
  createdAt: string;
  updatedAt: string;
}

export interface TaxRule {
  id: string;
  name: string;
  rateBps: number;
  country: string;
  isActive: boolean;
  createdAt: string;
}

export interface DiscountRule {
  id: string;
  code: string;
  discountType: 'FIXED_AMOUNT' | 'PERCENTAGE';
  valueCents: number;
  minPurchaseCents: number;
  minQuantity: number;
  maxUsageCount?: number;
  currentUsageCount: number;
  appliesTo: 'PRODUCT' | 'SUBSCRIPTION';
  productId?: string;
  validFrom: string;
  validTo?: string;
  isActive: boolean;
  createdAt: string;
}

// Plans
export type PlanFrequency = 'DAILY' | 'WEEKLY' | 'MONTHLY' | 'YEARLY';

export interface Plan {
  id: string;
  name: string;
  frequency: PlanFrequency;
  billingPriceCents: number;
  minQuantity: number;
  availableFrom: string;
  availableTo?: string;
  autoClose: boolean;
  allowEarlyClose: boolean;
  allowPause: boolean;
  allowRenewal: boolean;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

// Quotations
export type QuotationStatus = 'DRAFT' | 'SENT' | 'ACCEPTED' | 'EXPIRED' | 'CONVERTED';

export interface QuotationLine {
  id?: string;
  quotationId?: string;
  productId: string;
  productName?: string;
  variantId?: string;
  variantName?: string;
  quantity: number;
  unitPriceCents: number;
  taxRateBps: number;
  totalCents: number;
}

export interface Quotation {
  id: string;
  quotationNumber: string;
  customerId: string;
  customerEmail?: string;
  planId: string;
  planName?: string;
  status: QuotationStatus;
  validUntil: string;
  subtotalCents: number;
  taxCents: number;
  discountCents: number;
  totalCents: number;
  lines?: QuotationLine[];
  createdAt: string;
  updatedAt: string;
}

// Subscriptions
export type SubscriptionStatus =
  | 'DRAFT'
  | 'QUOTATION'
  | 'CONFIRMED'
  | 'ACTIVE'
  | 'PAUSED'
  | 'CANCELLED'
  | 'CLOSED';

export interface SubscriptionLine {
  id?: string;
  subscriptionId?: string;
  productId: string;
  productName?: string;
  variantId?: string;
  variantName?: string;
  isAddon: boolean;
  quantity: number;
  unitPriceCents: number;
  taxRateBps: number;
  lineTotalCents: number;
  createdAt?: string;
}

export interface SubscriptionEvent {
  id: string;
  subscriptionId: string;
  eventType: string;
  actorId?: string;
  actorEmail?: string;
  payload: Record<string, any>;
  createdAt: string;
}

export interface Subscription {
  id: string;
  subscriptionNumber: string;
  customerId: string;
  customerEmail?: string;
  planId: string;
  planName?: string;
  quotationId?: string;
  status: SubscriptionStatus;
  startDate: string;
  expiryDate?: string;
  nextBillingDate: string;
  paymentTerms: string;
  currency: string;
  lines?: SubscriptionLine[];
  totalAmountCents?: number;
  createdAt: string;
  updatedAt: string;
}

// Invoices & Billing
export type InvoiceStatus = 'DRAFT' | 'CONFIRMED' | 'PAID' | 'CANCELLED';

export interface InvoiceLine {
  id: string;
  invoiceId: string;
  productId: string;
  productName?: string;
  variantId?: string;
  variantName?: string;
  description: string;
  quantity: number;
  unitPriceCents: number;
  taxCents: number;
  discountCents: number;
  totalCents: number;
  createdAt: string;
}

export interface Invoice {
  id: string;
  invoiceNumber: string;
  subscriptionId: string;
  subscriptionNumber?: string;
  customerId: string;
  customerEmail?: string;
  billingRunId?: string;
  status: InvoiceStatus;
  billingPeriodStart: string;
  billingPeriodEnd: string;
  currency: string;
  subtotalCents: number;
  taxCents: number;
  discountCents: number;
  totalCents: number;
  paidAmountCents: number;
  outstandingBalanceCents: number;
  dueDate: string;
  paidAt?: string;
  lines?: InvoiceLine[];
  createdAt: string;
  updatedAt: string;
}

export interface BillingRun {
  id: string;
  runDate: string;
  status: string;
  subscriptionsEvaluated: number;
  invoicesGenerated: number;
  failuresCount: number;
  errorSummary?: string;
  createdAt: string;
  completedAt?: string;
}

export interface BillingAttempt {
  id: string;
  billingRunId: string;
  subscriptionId: string;
  subscriptionNumber?: string;
  status: string;
  errorMessage?: string;
  attemptNumber: number;
  createdAt: string;
}

export interface ReliabilityMetrics {
  subscriptionsDueCount: number;
  recentRuns: BillingRun[];
  recentAttempts: BillingAttempt[];
  duplicatePreventedCount: number;
}

export interface ExplainableBreakdown {
  basePriceCents: number;
  variantAdjustmentCents: number;
  effectiveUnitPriceCents: number;
  quantity: number;
  grossAmountCents: number;
  discountAmountCents: number;
  taxAmountCents: number;
  finalTotalCents: number;
}

// Payments
export interface Payment {
  id: string;
  invoiceId: string;
  invoiceNumber?: string;
  customerEmail?: string;
  paymentNumber: string;
  amountCents: number;
  paymentMethod: string;
  paymentDate: string;
  notes?: string;
  createdAt: string;
}

// Reporting
export interface MonthlyTrendPoint {
  month: string;
  invoicedRevenueCents: number;
  collectedCents: number;
}

export interface DashboardMetrics {
  activeSubscriptionsCount: number;
  mrrCents: number;
  invoicedRevenueCents: number;
  collectedPaymentsCents: number;
  outstandingReceivablesCents: number;
  overdueInvoicesCount: number;
  trend?: MonthlyTrendPoint[];
}
