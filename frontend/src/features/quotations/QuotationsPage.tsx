import { useEffect, useState } from 'react';
import { getQuotations, createQuotation, convertQuotation, getPlans, getProducts } from '@/lib/api/modules';
import { Quotation, Plan, Product } from '@/types';
import { formatCurrency } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert } from '@/components/ui/alert';
import { Plus, RefreshCw, CheckCircle2 } from 'lucide-react';

export function QuotationsPage() {
  const [quotations, setQuotations] = useState<Quotation[]>([]);
  const [plans, setPlans] = useState<Plan[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  // Form State
  const [selectedPlanId, setSelectedPlanId] = useState('');
  const [selectedProductId, setSelectedProductId] = useState('');
  const [quantity, setQuantity] = useState('1');
  const [unitPrice, setUnitPrice] = useState('4999');
  const [validityDays, setValidityDays] = useState('30');
  const [customerEmail, setCustomerEmail] = useState('customer@brightloop.local');

  const loadData = async () => {
    setIsLoading(true);
    try {
      const [quoteRes, planRes, prodRes] = await Promise.all([
        getQuotations(),
        getPlans(),
        getProducts(true),
      ]);
      setQuotations(quoteRes.quotations || []);
      setPlans(planRes.plans || []);
      setProducts(prodRes.products || []);
      if (planRes.plans && planRes.plans.length > 0 && !selectedPlanId) {
        setSelectedPlanId(planRes.plans[0].id);
      }
      if (prodRes.products && prodRes.products.length > 0 && !selectedProductId) {
        setSelectedProductId(prodRes.products[0].id);
      }
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleCreateQuotation = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      // Find a customer ID from existing user or mock test user
      await createQuotation({
        planId: selectedPlanId,
        validityDays: parseInt(validityDays, 10),
        lines: [
          {
            productId: selectedProductId,
            quantity: parseInt(quantity, 10),
            unitPriceCents: Math.round(parseFloat(unitPrice) * 100),
            taxRateBps: 1800,
          },
        ],
      });
      setShowModal(false);
      setMessage('Quotation proposal generated successfully');
      loadData();
    } catch (err: any) {
      alert(err.message || 'Failed creating quotation');
    }
  };

  const handleConvert = async (quoteId: string) => {
    if (!confirm('Convert this approved quotation into an active subscription?')) return;
    try {
      const res = await convertQuotation(quoteId);
      setMessage(`Quotation converted! Active Subscription ID: ${res.subscriptionId}`);
      loadData();
    } catch (err: any) {
      alert(err.message || 'Failed converting quotation');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">Quotations & Proposals</h1>
          <p className="text-sm text-slate-500 mt-1">
            Build customized client quotes with preconfigured plans, and seamlessly convert accepted proposals into active subscriptions.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button variant="outline" size="sm" onClick={loadData} isLoading={isLoading}>
            <RefreshCw className="h-4 w-4 mr-1 text-slate-500" /> Refresh
          </Button>
          <Button variant="primary" size="sm" onClick={() => setShowModal(true)}>
            <Plus className="h-4 w-4 mr-1.5" /> Draft Quotation
          </Button>
        </div>
      </div>

      {message && (
        <Alert variant="success" className="cursor-pointer" onClick={() => setMessage(null)}>
          {message}
        </Alert>
      )}

      <Card>
        <CardHeader>
          <CardTitle className="text-base">All Proposals</CardTitle>
        </CardHeader>
        <CardContent>
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-50 text-slate-500 text-xs uppercase border-b">
              <tr>
                <th className="py-2.5 px-4 font-semibold">Quote #</th>
                <th className="py-2.5 px-4 font-semibold">Customer</th>
                <th className="py-2.5 px-4 font-semibold">Target Plan</th>
                <th className="py-2.5 px-4 font-semibold">Total Amount</th>
                <th className="py-2.5 px-4 font-semibold">Valid Until</th>
                <th className="py-2.5 px-4 font-semibold">Status</th>
                <th className="py-2.5 px-4 font-semibold text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {quotations.length > 0 ? (
                quotations.map((q) => (
                  <tr key={q.id}>
                    <td className="py-3 px-4 font-mono font-bold text-slate-900">{q.quotationNumber}</td>
                    <td className="py-3 px-4 text-slate-700">{q.customerEmail || 'customer'}</td>
                    <td className="py-3 px-4 text-slate-700">{q.planName || 'Standard'}</td>
                    <td className="py-3 px-4 font-bold text-slate-900">{formatCurrency(q.totalCents)}</td>
                    <td className="py-3 px-4 text-slate-500">
                      {q.validUntil ? new Date(q.validUntil).toLocaleDateString() : '—'}
                    </td>
                    <td className="py-3 px-4">
                      <Badge
                        variant={
                          q.status === 'CONVERTED'
                            ? 'success'
                            : q.status === 'ACCEPTED'
                            ? 'indigo'
                            : 'default'
                        }
                      >
                        {q.status}
                      </Badge>
                    </td>
                    <td className="py-3 px-4 text-right">
                      {q.status !== 'CONVERTED' ? (
                        <Button
                          variant="primary"
                          size="sm"
                          onClick={() => handleConvert(q.id)}
                          className="text-xs"
                        >
                          <CheckCircle2 className="h-3.5 w-3.5 mr-1" /> Convert to Sub
                        </Button>
                      ) : (
                        <span className="text-xs text-emerald-700 font-medium">Converted</span>
                      )}
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={7} className="py-8 text-center text-xs text-slate-400">
                    No quotations drafted yet. Click 'Draft Quotation' to create a proposal.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </CardContent>
      </Card>

      {/* Modal: Draft Quotation */}
      {showModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-md w-full shadow-lg">
            <CardHeader>
              <CardTitle className="text-lg">Draft New Proposal</CardTitle>
            </CardHeader>
            <form onSubmit={handleCreateQuotation}>
              <CardContent className="space-y-4">
                <div>
                  <Label>Recipient Email</Label>
                  <Input
                    required
                    value={customerEmail}
                    onChange={(e) => setCustomerEmail(e.target.value)}
                  />
                </div>

                <div>
                  <Label>Subscription Plan</Label>
                  <select
                    className="w-full border rounded px-3 py-2 text-sm bg-white"
                    value={selectedPlanId}
                    onChange={(e) => setSelectedPlanId(e.target.value)}
                  >
                    {plans.map((p) => (
                      <option key={p.id} value={p.id}>
                        {p.name} ({p.frequency} &bull; {formatCurrency(p.billingPriceCents)})
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <Label>Core Product Line</Label>
                  <select
                    className="w-full border rounded px-3 py-2 text-sm bg-white"
                    value={selectedProductId}
                    onChange={(e) => setSelectedProductId(e.target.value)}
                  >
                    {products.map((pr) => (
                      <option key={pr.id} value={pr.id}>
                        {pr.name} ({formatCurrency(pr.salesPriceCents)})
                      </option>
                    ))}
                  </select>
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <Label>Quantity</Label>
                    <Input
                      type="number"
                      min="1"
                      required
                      value={quantity}
                      onChange={(e) => setQuantity(e.target.value)}
                    />
                  </div>
                  <div>
                    <Label>Unit Price (₹)</Label>
                    <Input
                      type="number"
                      required
                      value={unitPrice}
                      onChange={(e) => setUnitPrice(e.target.value)}
                    />
                  </div>
                </div>

                <div>
                  <Label>Proposal Validity (Days)</Label>
                  <Input
                    type="number"
                    min="1"
                    required
                    value={validityDays}
                    onChange={(e) => setValidityDays(e.target.value)}
                  />
                </div>
              </CardContent>
              <div className="p-4 border-t flex justify-end gap-2 bg-slate-50 rounded-b-lg">
                <Button type="button" variant="outline" onClick={() => setShowModal(false)}>
                  Cancel
                </Button>
                <Button type="submit" variant="primary">
                  Generate Quotation
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  );
}
