import { apiClient } from './client';
import {
  Product,
  ProductVariant,
  TaxRule,
  DiscountRule,
  Plan,
  Quotation,
  Subscription,
  SubscriptionEvent,
  Invoice,
  Payment,
  DashboardMetrics,
  ReliabilityMetrics,
  ExplainableBreakdown,
} from '@/types';

// Catalog
export async function getProducts(activeOnly = false): Promise<{ products: Product[] }> {
  return apiClient<{ products: Product[] }>(`/catalog/products${activeOnly ? '?active=true' : ''}`);
}

export async function createProduct(payload: Partial<Product>): Promise<{ product: Product }> {
  return apiClient<{ product: Product }>('/catalog/products', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function updateProduct(id: string, payload: Partial<Product>): Promise<{ product: Product }> {
  return apiClient<{ product: Product }>(`/catalog/products/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
}

export async function createVariant(productId: string, payload: Partial<ProductVariant>): Promise<{ variant: ProductVariant }> {
  return apiClient<{ variant: ProductVariant }>(`/catalog/products/${productId}/variants`, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function getTaxes(): Promise<{ taxes: TaxRule[] }> {
  return apiClient<{ taxes: TaxRule[] }>('/taxes');
}

export async function createTax(payload: Partial<TaxRule>): Promise<{ tax: TaxRule }> {
  return apiClient<{ tax: TaxRule }>('/taxes', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function getDiscounts(): Promise<{ discounts: DiscountRule[] }> {
  return apiClient<{ discounts: DiscountRule[] }>('/discounts');
}

export async function createDiscount(payload: Partial<DiscountRule>): Promise<{ discount: DiscountRule }> {
  return apiClient<{ discount: DiscountRule }>('/discounts', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

// Plans
export async function getPlans(): Promise<{ plans: Plan[] }> {
  return apiClient<{ plans: Plan[] }>('/plans');
}

export async function createPlan(payload: Partial<Plan>): Promise<{ plan: Plan }> {
  return apiClient<{ plan: Plan }>('/plans', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

// Quotations
export async function getQuotations(customerId?: string): Promise<{ quotations: Quotation[] }> {
  return apiClient<{ quotations: Quotation[] }>(`/quotations${customerId ? `?customerId=${customerId}` : ''}`);
}

export async function getQuotation(id: string): Promise<{ quotation: Quotation }> {
  return apiClient<{ quotation: Quotation }>(`/quotations/${id}`);
}

export async function createQuotation(payload: any): Promise<{ quotation: Quotation }> {
  return apiClient<{ quotation: Quotation }>('/quotations', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function convertQuotation(id: string): Promise<{ message: string; subscriptionId: string }> {
  return apiClient<{ message: string; subscriptionId: string }>(`/quotations/${id}/convert`, {
    method: 'POST',
  });
}

// Subscriptions
export async function getSubscriptions(customerId?: string, status?: string): Promise<{ subscriptions: Subscription[] }> {
  const params = new URLSearchParams();
  if (customerId) params.set('customerId', customerId);
  if (status) params.set('status', status);
  const q = params.toString() ? `?${params.toString()}` : '';
  return apiClient<{ subscriptions: Subscription[] }>(`/subscriptions${q}`);
}

export async function getSubscription(id: string): Promise<{ subscription: Subscription }> {
  return apiClient<{ subscription: Subscription }>(`/subscriptions/${id}`);
}

export async function createSubscription(payload: any): Promise<{ subscription: Subscription }> {
  return apiClient<{ subscription: Subscription }>('/subscriptions', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function activateSubscription(id: string): Promise<{ message: string }> {
  return apiClient<{ message: string }>(`/subscriptions/${id}/activate`, { method: 'POST' });
}

export async function pauseSubscription(id: string): Promise<{ message: string }> {
  return apiClient<{ message: string }>(`/subscriptions/${id}/pause`, { method: 'POST' });
}

export async function resumeSubscription(id: string): Promise<{ message: string }> {
  return apiClient<{ message: string }>(`/subscriptions/${id}/resume`, { method: 'POST' });
}

export async function cancelSubscription(id: string): Promise<{ message: string }> {
  return apiClient<{ message: string }>(`/subscriptions/${id}/cancel`, { method: 'POST' });
}

export async function toggleAddon(subId: string, productId: string, active: boolean): Promise<{ message: string }> {
  return apiClient<{ message: string }>(`/subscriptions/${subId}/addons`, {
    method: 'POST',
    body: JSON.stringify({ productId, active }),
  });
}

export async function getSubscriptionTimeline(id: string): Promise<{ events: SubscriptionEvent[] }> {
  return apiClient<{ events: SubscriptionEvent[] }>(`/subscriptions/${id}/timeline`);
}

// Billing & Invoices
export async function getInvoices(customerId?: string, status?: string): Promise<{ invoices: Invoice[] }> {
  const params = new URLSearchParams();
  if (customerId) params.set('customerId', customerId);
  if (status) params.set('status', status);
  const q = params.toString() ? `?${params.toString()}` : '';
  return apiClient<{ invoices: Invoice[] }>(`/invoices${q}`);
}

export async function getInvoice(id: string): Promise<{ invoice: Invoice }> {
  return apiClient<{ invoice: Invoice }>(`/invoices/${id}`);
}

export async function confirmInvoice(id: string): Promise<{ message: string }> {
  return apiClient<{ message: string }>(`/invoices/${id}/confirm`, { method: 'POST' });
}

export async function getReliabilityMetrics(): Promise<{ metrics: ReliabilityMetrics }> {
  return apiClient<{ metrics: ReliabilityMetrics }>('/billing/reliability');
}

export async function triggerBillingRun(): Promise<{ message: string; run: any }> {
  return apiClient<{ message: string; run: any }>('/billing/run', { method: 'POST' });
}

export async function previewExplainableBilling(payload: {
  basePriceCents: number;
  variantAdjustmentCents: number;
  quantity: number;
  discountCents: number;
  taxRateBps: number;
}): Promise<{ breakdown: ExplainableBreakdown }> {
  return apiClient<{ breakdown: ExplainableBreakdown }>('/billing/preview', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

// Payments
export async function getPayments(invoiceId?: string, customerId?: string): Promise<{ payments: Payment[] }> {
  const params = new URLSearchParams();
  if (invoiceId) params.set('invoiceId', invoiceId);
  if (customerId) params.set('customerId', customerId);
  const q = params.toString() ? `?${params.toString()}` : '';
  return apiClient<{ payments: Payment[] }>(`/payments${q}`);
}

export async function recordPayment(invoiceId: string, payload: {
  amountCents: number;
  method: string;
  notes?: string;
}): Promise<{ payment: Payment }> {
  return apiClient<{ payment: Payment }>(`/invoices/${invoiceId}/payments`, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

// Reporting
export async function getDashboardMetrics(): Promise<{ metrics: DashboardMetrics }> {
  return apiClient<{ metrics: DashboardMetrics }>('/reporting/dashboard');
}
