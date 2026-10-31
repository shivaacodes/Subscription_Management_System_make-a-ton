import { useEffect, useState } from 'react';
import { getPayments, getInvoices, recordPayment } from '@/lib/api/modules';
import { Payment, Invoice } from '@/types';
import { formatCurrency } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert } from '@/components/ui/alert';
import { CreditCard, RefreshCw } from 'lucide-react';

export function PaymentsPage() {
  const [payments, setPayments] = useState<Payment[]>([]);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [message, setMessage] = useState<string | null>(null);

  // Record Payment Modal
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);
  const [amount, setAmount] = useState('');
  const [method, setMethod] = useState('BANK_TRANSFER');
  const [notes, setNotes] = useState('');

  const loadData = async () => {
    setIsLoading(true);
    try {
      const [payRes, invRes] = await Promise.all([
        getPayments(),
        getInvoices(),
      ]);
      setPayments(payRes.payments || []);
      setInvoices(invRes.invoices || []);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const openPaymentModal = (inv: Invoice) => {
    setSelectedInvoice(inv);
    // pre-fill with outstanding balance in rupees
    setAmount((inv.outstandingBalanceCents / 100).toFixed(2));
    setNotes(`Payment for ${inv.invoiceNumber}`);
  };

  const handleRecordPayment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedInvoice) return;
    try {
      const amountCents = Math.round(parseFloat(amount) * 100);
      await recordPayment(selectedInvoice.id, {
        amountCents,
        method,
        notes: notes || undefined,
      });
      setSelectedInvoice(null);
      setMessage(`Payment recorded against ${selectedInvoice.invoiceNumber} successfully`);
      loadData();
    } catch (err: any) {
      alert(err.message || 'Payment recording failed');
    }
  };

  const unpaidInvoices = invoices.filter((i) => i.outstandingBalanceCents > 0);

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">Payments & Cash Collection</h1>
          <p className="text-sm text-slate-500 mt-1">
            Record multi-part payments against invoices, manage balances, and ensure payments never exceed outstanding liabilities.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button variant="outline" size="sm" onClick={loadData} isLoading={isLoading}>
            <RefreshCw className="h-4 w-4 mr-1 text-slate-500" /> Refresh
          </Button>
        </div>
      </div>

      {message && (
        <Alert variant="success" className="cursor-pointer" onClick={() => setMessage(null)}>
          {message}
        </Alert>
      )}

      {/* Invoices Awaiting Payment Card */}
      {unpaidInvoices.length > 0 && (
        <Card className="border-amber-200 bg-amber-50/20">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-base text-amber-950">Invoices Awaiting Payment</CardTitle>
                <p className="text-xs text-amber-800">
                  {unpaidInvoices.length} statement(s) with outstanding receivables.
                </p>
              </div>
              <Badge variant="warning">{unpaidInvoices.length} Pending</Badge>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
              {unpaidInvoices.map((inv) => (
                <div key={inv.id} className="p-4 bg-white rounded-md border border-amber-200 flex flex-col justify-between">
                  <div>
                    <div className="flex items-center justify-between mb-1">
                      <span className="font-mono font-bold text-xs text-indigo-600">{inv.invoiceNumber}</span>
                      <span className="text-xs font-semibold text-rose-600">
                        Balance: {formatCurrency(inv.outstandingBalanceCents)}
                      </span>
                    </div>
                    <p className="text-xs text-slate-600 truncate">{inv.customerEmail}</p>
                    <p className="text-[11px] text-slate-400">Total: {formatCurrency(inv.totalCents)}</p>
                  </div>
                  <div className="mt-3 pt-2 border-t flex justify-end">
                    <Button variant="primary" size="sm" onClick={() => openPaymentModal(inv)} className="text-xs w-full">
                      <CreditCard className="h-3.5 w-3.5 mr-1" /> Record Payment
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Payments History Table */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Collected Payments Ledger</CardTitle>
        </CardHeader>
        <CardContent>
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-50 text-slate-500 text-xs uppercase border-b">
              <tr>
                <th className="py-2.5 px-4 font-semibold">Payment #</th>
                <th className="py-2.5 px-4 font-semibold">Invoice #</th>
                <th className="py-2.5 px-4 font-semibold">Customer</th>
                <th className="py-2.5 px-4 font-semibold">Amount</th>
                <th className="py-2.5 px-4 font-semibold">Method</th>
                <th className="py-2.5 px-4 font-semibold">Date</th>
                <th className="py-2.5 px-4 font-semibold">Notes</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {payments.length > 0 ? (
                payments.map((p) => (
                  <tr key={p.id}>
                    <td className="py-3 px-4 font-mono font-bold text-slate-900">{p.paymentNumber}</td>
                    <td className="py-3 px-4 font-mono text-indigo-600">{p.invoiceNumber || p.invoiceId.slice(0, 8)}</td>
                    <td className="py-3 px-4 text-slate-700">{p.customerEmail || 'Customer'}</td>
                    <td className="py-3 px-4 font-bold text-emerald-700">{formatCurrency(p.amountCents)}</td>
                    <td className="py-3 px-4 text-xs">
                      <Badge variant="indigo">{p.paymentMethod}</Badge>
                    </td>
                    <td className="py-3 px-4 text-xs text-slate-500">
                      {p.paymentDate ? new Date(p.paymentDate).toLocaleDateString() : '—'}
                    </td>
                    <td className="py-3 px-4 text-xs text-slate-500 truncate max-w-xs">{p.notes || '—'}</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={7} className="py-8 text-center text-xs text-slate-400">
                    No payment transactions recorded yet. Click 'Record Payment' on an outstanding invoice.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </CardContent>
      </Card>

      {/* Modal: Record Payment */}
      {selectedInvoice && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-md w-full shadow-lg">
            <CardHeader>
              <CardTitle className="text-lg">Record Payment</CardTitle>
              <p className="text-xs text-slate-500">
                Invoice {selectedInvoice.invoiceNumber} &bull; Outstanding: {formatCurrency(selectedInvoice.outstandingBalanceCents)}
              </p>
            </CardHeader>
            <form onSubmit={handleRecordPayment}>
              <CardContent className="space-y-4">
                <div>
                  <Label>Payment Amount (₹)</Label>
                  <Input
                    type="number"
                    step="0.01"
                    max={(selectedInvoice.outstandingBalanceCents / 100).toFixed(2)}
                    required
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
                  />
                  <p className="text-[11px] text-slate-500 mt-1">
                    Cannot exceed remaining balance of {formatCurrency(selectedInvoice.outstandingBalanceCents)}.
                  </p>
                </div>

                <div>
                  <Label>Payment Method</Label>
                  <select
                    className="w-full border rounded px-3 py-2 text-sm bg-white"
                    value={method}
                    onChange={(e) => setMethod(e.target.value)}
                  >
                    <option value="BANK_TRANSFER">Bank Wire / NEFT / RTGS</option>
                    <option value="UPI">UPI / QR Code</option>
                    <option value="CREDIT_CARD">Credit Card</option>
                    <option value="CHEQUE">Corporate Cheque</option>
                  </select>
                </div>

                <div>
                  <Label>Payment Reference / Notes</Label>
                  <Input
                    placeholder="e.g. UTR / Transaction reference"
                    value={notes}
                    onChange={(e) => setNotes(e.target.value)}
                  />
                </div>
              </CardContent>
              <div className="p-4 border-t flex justify-end gap-2 bg-slate-50 rounded-b-lg">
                <Button type="button" variant="outline" onClick={() => setSelectedInvoice(null)}>
                  Cancel
                </Button>
                <Button type="submit" variant="primary">
                  Confirm Payment
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  );
}
