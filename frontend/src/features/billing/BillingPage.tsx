import { useEffect, useState } from 'react';
import {
  getInvoices,
  getReliabilityMetrics,
  triggerBillingRun,
  previewExplainableBilling,
} from '@/lib/api/modules';
import { Invoice, ReliabilityMetrics, ExplainableBreakdown } from '@/types';
import { formatCurrency } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert } from '@/components/ui/alert';
import {
  Receipt,
  Play,
  ShieldCheck,
  CheckCircle2,
  RefreshCw,
  Calculator,
  Printer,
  Calendar,
} from 'lucide-react';

export function BillingPage() {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [metrics, setMetrics] = useState<ReliabilityMetrics | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isRunningBilling, setIsRunningBilling] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  // Explainable Billing Modal
  const [showExplainModal, setShowExplainModal] = useState(false);
  const [breakdown, setBreakdown] = useState<ExplainableBreakdown | null>(null);
  const [calcBasePrice, setCalcBasePrice] = useState('4999');
  const [calcVariantAdj, setCalcVariantAdj] = useState('560'); // Odoo variant +₹560
  const [calcQty, setCalcQty] = useState('2');
  const [calcDiscount, setCalcDiscount] = useState('500'); // ₹500 discount
  const [calcTaxBps, setCalcTaxBps] = useState('1800'); // 18% GST

  // Invoice Print Modal
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);

  const loadData = async () => {
    setIsLoading(true);
    try {
      const [invRes, relRes] = await Promise.all([
        getInvoices(),
        getReliabilityMetrics(),
      ]);
      setInvoices(invRes.invoices || []);
      setMetrics(relRes.metrics || null);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleRunBilling = async () => {
    setIsRunningBilling(true);
    setMessage(null);
    try {
      const res = await triggerBillingRun();
      setMessage(
        `Automated billing cycle complete: Evaluated ${res.run.subscriptionsEvaluated} subscription(s), Generated ${res.run.invoicesGenerated} invoice(s), ${res.run.failuresCount} failures.`
      );
      loadData();
    } catch (err: any) {
      alert(err.message || 'Billing run failed');
    } finally {
      setIsRunningBilling(false);
    }
  };

  const handlePreviewCalculation = async (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    try {
      const res = await previewExplainableBilling({
        basePriceCents: Math.round(parseFloat(calcBasePrice) * 100),
        variantAdjustmentCents: Math.round(parseFloat(calcVariantAdj) * 100),
        quantity: parseInt(calcQty, 10),
        discountCents: Math.round(parseFloat(calcDiscount) * 100),
        taxRateBps: parseInt(calcTaxBps, 10),
      });
      setBreakdown(res.breakdown);
    } catch (err: any) {
      alert(err.message || 'Calculation preview failed');
    }
  };

  useEffect(() => {
    if (showExplainModal) {
      handlePreviewCalculation();
    }
  }, [showExplainModal]);

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-slate-900">Billing Engine & Reliability Center</h1>
            <Badge variant="indigo">Differentiator</Badge>
          </div>
          <p className="text-sm text-slate-500">
            Idempotent batch recurring invoice generation, duplicate prevention enforcement, and explainable pricing breakdowns.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button variant="outline" size="sm" onClick={loadData} isLoading={isLoading}>
            <RefreshCw className="h-4 w-4 mr-1 text-slate-500" /> Refresh
          </Button>
          <Button variant="secondary" size="sm" onClick={() => setShowExplainModal(true)}>
            <Calculator className="h-4 w-4 mr-1.5 text-indigo-600" /> Explainable Calculator
          </Button>
          <Button variant="primary" size="sm" onClick={handleRunBilling} isLoading={isRunningBilling}>
            <Play className="h-4 w-4 mr-1.5" /> Execute Billing Batch
          </Button>
        </div>
      </div>

      {message && (
        <Alert variant="success" className="cursor-pointer" onClick={() => setMessage(null)}>
          {message}
        </Alert>
      )}

      {/* Billing Reliability Center Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="border-indigo-100 bg-indigo-50/40">
          <CardHeader className="pb-1">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-600 uppercase tracking-wider">Due for Billing</span>
              <Calendar className="h-4 w-4 text-indigo-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-indigo-950">
              {metrics ? metrics.subscriptionsDueCount : 0}
            </div>
            <p className="text-[11px] text-slate-500 mt-1">Subscriptions ready to invoice</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-1">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Idempotency Guard</span>
              <ShieldCheck className="h-4 w-4 text-emerald-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-emerald-700">Active (DB Unique)</div>
            <p className="text-[11px] text-slate-500 mt-1">Zero duplicate period billing</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-1">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Duplicates Blocked</span>
              <CheckCircle2 className="h-4 w-4 text-emerald-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">
              {metrics ? metrics.duplicatePreventedCount : 0}
            </div>
            <p className="text-[11px] text-slate-500 mt-1">Duplicate cycles prevented</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-1">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Recent Executions</span>
              <Receipt className="h-4 w-4 text-slate-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">
              {metrics?.recentRuns?.length ?? 0}
            </div>
            <p className="text-[11px] text-slate-500 mt-1">Batch runs recorded</p>
          </CardContent>
        </Card>
      </div>

      {/* Generated Invoices Table */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle className="text-base">Generated Invoices & Statements</CardTitle>
          <span className="text-xs text-slate-500">{(invoices || []).length} statement(s) issued</span>
        </CardHeader>
        <CardContent>
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-50 text-slate-500 text-xs uppercase border-b">
              <tr>
                <th className="py-2.5 px-4 font-semibold">Invoice #</th>
                <th className="py-2.5 px-4 font-semibold">Subscription #</th>
                <th className="py-2.5 px-4 font-semibold">Customer</th>
                <th className="py-2.5 px-4 font-semibold">Period</th>
                <th className="py-2.5 px-4 font-semibold">Total Amount</th>
                <th className="py-2.5 px-4 font-semibold">Balance</th>
                <th className="py-2.5 px-4 font-semibold">Status</th>
                <th className="py-2.5 px-4 font-semibold text-right">View / Print</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {(invoices || []).length > 0 ? (
                invoices.map((inv) => (
                  <tr key={inv.id}>
                    <td className="py-3 px-4 font-mono font-bold text-indigo-600">{inv.invoiceNumber}</td>
                    <td className="py-3 px-4 font-mono text-slate-600">{inv.subscriptionNumber || 'SUB'}</td>
                    <td className="py-3 px-4 text-slate-700">{inv.customerEmail}</td>
                    <td className="py-3 px-4 text-xs text-slate-500">
                      {new Date(inv.billingPeriodStart).toLocaleDateString()} &rarr; {new Date(inv.billingPeriodEnd).toLocaleDateString()}
                    </td>
                    <td className="py-3 px-4 font-bold text-slate-900">{formatCurrency(inv.totalCents)}</td>
                    <td className="py-3 px-4 font-semibold text-slate-700">
                      {inv.outstandingBalanceCents > 0 ? (
                        <span className="text-amber-600">{formatCurrency(inv.outstandingBalanceCents)}</span>
                      ) : (
                        <span className="text-emerald-600">Paid in Full</span>
                      )}
                    </td>
                    <td className="py-3 px-4">
                      <Badge
                        variant={
                          inv.status === 'PAID'
                            ? 'success'
                            : inv.status === 'CONFIRMED'
                            ? 'indigo'
                            : 'default'
                        }
                      >
                        {inv.status}
                      </Badge>
                    </td>
                    <td className="py-3 px-4 text-right">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setSelectedInvoice(inv)}
                        className="text-xs"
                      >
                        <Printer className="h-3.5 w-3.5 mr-1" /> View Statement
                      </Button>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={8} className="py-8 text-center text-xs text-slate-400">
                    No invoices generated yet. Click 'Execute Billing Batch' to invoice all subscriptions due today.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </CardContent>
      </Card>

      {/* Modal: Explainable Billing Calculation */}
      {showExplainModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-2xl w-full shadow-2xl">
            <CardHeader className="border-b bg-slate-50/50">
              <div className="flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-2">
                    <CardTitle className="text-lg">Authoritative Pricing Breakdown</CardTitle>
                    <Badge variant="indigo">Explainable Billing</Badge>
                  </div>
                  <p className="text-xs text-slate-500 mt-1">
                    Transparent formula: Base Price &rarr; Variant Adjustment &rarr; Quantity &rarr; Discount &rarr; Tax &rarr; Final Total
                  </p>
                </div>
                <Button variant="ghost" size="sm" onClick={() => setShowExplainModal(false)}>
                  Close
                </Button>
              </div>
            </CardHeader>
            <CardContent className="p-6 space-y-6">
              {/* Interactive Tester Inputs */}
              <form onSubmit={handlePreviewCalculation} className="grid grid-cols-2 sm:grid-cols-5 gap-3 bg-slate-50 p-4 rounded-md border">
                <div>
                  <Label>Base Price (₹)</Label>
                  <Input value={calcBasePrice} onChange={(e) => setCalcBasePrice(e.target.value)} />
                </div>
                <div>
                  <Label>Variant (+₹)</Label>
                  <Input value={calcVariantAdj} onChange={(e) => setCalcVariantAdj(e.target.value)} />
                </div>
                <div>
                  <Label>Quantity</Label>
                  <Input type="number" min="1" value={calcQty} onChange={(e) => setCalcQty(e.target.value)} />
                </div>
                <div>
                  <Label>Discount (₹)</Label>
                  <Input value={calcDiscount} onChange={(e) => setCalcDiscount(e.target.value)} />
                </div>
                <div>
                  <Label>Tax (bps - 1800)</Label>
                  <Input value={calcTaxBps} onChange={(e) => setCalcTaxBps(e.target.value)} />
                </div>
              </form>

              {/* Visual Breakdown Diagram */}
              {breakdown && (
                <div className="space-y-4">
                  <div className="text-xs font-semibold uppercase tracking-wider text-slate-500">
                    Step-by-Step Computational Trail
                  </div>
                  <div className="grid grid-cols-1 sm:grid-cols-6 gap-2 text-center text-xs">
                    <div className="p-3 bg-white rounded border border-slate-200">
                      <div className="text-slate-500 mb-1">Base Price</div>
                      <div className="font-bold text-slate-900">{formatCurrency(breakdown.basePriceCents)}</div>
                    </div>
                    <div className="p-3 bg-white rounded border border-slate-200">
                      <div className="text-slate-500 mb-1">+ Variant</div>
                      <div className="font-bold text-indigo-600">+{formatCurrency(breakdown.variantAdjustmentCents)}</div>
                    </div>
                    <div className="p-3 bg-white rounded border border-slate-200">
                      <div className="text-slate-500 mb-1">&times; Quantity</div>
                      <div className="font-bold text-slate-900">{breakdown.quantity} unit(s)</div>
                    </div>
                    <div className="p-3 bg-white rounded border border-slate-200">
                      <div className="text-slate-500 mb-1">- Discount</div>
                      <div className="font-bold text-emerald-600">-{formatCurrency(breakdown.discountAmountCents)}</div>
                    </div>
                    <div className="p-3 bg-white rounded border border-slate-200">
                      <div className="text-slate-500 mb-1">+ 18% Tax</div>
                      <div className="font-bold text-slate-900">+{formatCurrency(breakdown.taxAmountCents)}</div>
                    </div>
                    <div className="p-3 bg-indigo-600 text-white rounded font-bold shadow-sm">
                      <div className="text-indigo-200 mb-1">Final Amount</div>
                      <div className="text-sm">{formatCurrency(breakdown.finalTotalCents)}</div>
                    </div>
                  </div>

                  <div className="p-4 bg-emerald-50/70 border border-emerald-200 rounded-md text-xs text-emerald-950 flex items-start gap-2">
                    <CheckCircle2 className="h-4 w-4 text-emerald-600 mt-0.5 flex-shrink-0" />
                    <div>
                      This authoritative backend calculation service is shared identically between invoice preview and batch invoice generation, ensuring zero pricing drift.
                    </div>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {/* Modal: View Printable Statement */}
      {selectedInvoice && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-2xl w-full shadow-2xl bg-white max-h-[90vh] overflow-y-auto">
            <CardHeader className="border-b flex flex-row items-center justify-between">
              <div>
                <CardTitle className="text-lg">Billing Statement</CardTitle>
                <p className="text-xs text-slate-500 font-mono">{selectedInvoice.invoiceNumber}</p>
              </div>
              <div className="flex items-center gap-2">
                <Button variant="outline" size="sm" onClick={() => window.print()} className="text-xs">
                  <Printer className="h-3.5 w-3.5 mr-1" /> Print / PDF
                </Button>
                <Button variant="ghost" size="sm" onClick={() => setSelectedInvoice(null)}>
                  Close
                </Button>
              </div>
            </CardHeader>
            <CardContent className="p-8 space-y-6">
              <div className="flex justify-between border-b pb-6">
                <div>
                  <h2 className="font-bold text-xl text-slate-900">Brightloop Software Inc.</h2>
                  <p className="text-xs text-slate-500">Veyra Subscription Billing</p>
                  <p className="text-xs text-slate-500">GSTIN: 29AABCB1234F1Z0</p>
                </div>
                <div className="text-right text-xs text-slate-600">
                  <div className="font-bold text-sm text-slate-900">{selectedInvoice.invoiceNumber}</div>
                  <div>Issued: {new Date(selectedInvoice.createdAt).toLocaleDateString()}</div>
                  <div>Due: {new Date(selectedInvoice.dueDate).toLocaleDateString()}</div>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 text-xs">
                <div>
                  <span className="text-slate-400 uppercase font-semibold">Billed To:</span>
                  <div className="font-bold text-slate-900 mt-1">{selectedInvoice.customerEmail}</div>
                </div>
                <div className="text-right">
                  <span className="text-slate-400 uppercase font-semibold">Billing Period:</span>
                  <div className="text-slate-700 mt-1">
                    {new Date(selectedInvoice.billingPeriodStart).toLocaleDateString()} &rarr; {new Date(selectedInvoice.billingPeriodEnd).toLocaleDateString()}
                  </div>
                </div>
              </div>

              {/* Line items snapshot */}
              <table className="w-full text-left text-xs border-t border-b">
                <thead className="bg-slate-50 text-slate-600">
                  <tr>
                    <th className="py-2 px-3">Description</th>
                    <th className="py-2 px-3 text-center">Qty</th>
                    <th className="py-2 px-3 text-right">Unit Price</th>
                    <th className="py-2 px-3 text-right">Tax</th>
                    <th className="py-2 px-3 text-right">Total</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {selectedInvoice.lines && selectedInvoice.lines.length > 0 ? (
                    selectedInvoice.lines.map((l) => (
                      <tr key={l.id}>
                        <td className="py-2.5 px-3 font-medium text-slate-900">{l.description}</td>
                        <td className="py-2.5 px-3 text-center text-slate-600">{l.quantity}</td>
                        <td className="py-2.5 px-3 text-right text-slate-600">{formatCurrency(l.unitPriceCents)}</td>
                        <td className="py-2.5 px-3 text-right text-slate-600">{formatCurrency(l.taxCents)}</td>
                        <td className="py-2.5 px-3 text-right font-bold text-slate-900">{formatCurrency(l.totalCents)}</td>
                      </tr>
                    ))
                  ) : (
                    <tr>
                      <td className="py-2.5 px-3 font-medium">Subscription recurring cycle statement</td>
                      <td className="py-2.5 px-3 text-center">1</td>
                      <td className="py-2.5 px-3 text-right">{formatCurrency(selectedInvoice.subtotalCents)}</td>
                      <td className="py-2.5 px-3 text-right">{formatCurrency(selectedInvoice.taxCents)}</td>
                      <td className="py-2.5 px-3 text-right font-bold">{formatCurrency(selectedInvoice.totalCents)}</td>
                    </tr>
                  )}
                </tbody>
              </table>

              <div className="flex justify-end text-xs">
                <div className="w-64 space-y-1.5">
                  <div className="flex justify-between text-slate-500">
                    <span>Subtotal:</span>
                    <span>{formatCurrency(selectedInvoice.subtotalCents)}</span>
                  </div>
                  <div className="flex justify-between text-slate-500">
                    <span>Taxes (18% GST):</span>
                    <span>{formatCurrency(selectedInvoice.taxCents)}</span>
                  </div>
                  <div className="flex justify-between font-bold text-sm text-slate-900 border-t pt-2">
                    <span>Total Statement:</span>
                    <span>{formatCurrency(selectedInvoice.totalCents)}</span>
                  </div>
                  <div className="flex justify-between text-emerald-700 font-semibold">
                    <span>Paid to Date:</span>
                    <span>-{formatCurrency(selectedInvoice.paidAmountCents)}</span>
                  </div>
                  <div className="flex justify-between text-slate-900 font-bold border-t pt-1">
                    <span>Outstanding Balance:</span>
                    <span>{formatCurrency(selectedInvoice.outstandingBalanceCents)}</span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}
