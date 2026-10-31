import { useEffect, useState } from 'react';
import { getInvoices } from '@/lib/api/modules';
import { Invoice } from '@/types';
import { formatCurrency } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Printer, RefreshCw } from 'lucide-react';

export function CustomerInvoicesPage() {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);

  const loadInvoices = async () => {
    setIsLoading(true);
    try {
      const res = await getInvoices();
      setInvoices(res.invoices || []);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadInvoices();
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">Invoices & Receipts</h1>
          <p className="text-sm text-slate-500 mt-1">
            Review past billing statements, payment receipts, and download PDF copies for accounting.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={loadInvoices} isLoading={isLoading}>
          <RefreshCw className="h-4 w-4 mr-1 text-slate-500" /> Refresh
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Billing Statements</CardTitle>
        </CardHeader>
        <CardContent>
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-50 text-slate-500 text-xs uppercase border-b">
              <tr>
                <th className="py-2.5 px-4 font-semibold">Invoice #</th>
                <th className="py-2.5 px-4 font-semibold">Date</th>
                <th className="py-2.5 px-4 font-semibold">Period</th>
                <th className="py-2.5 px-4 font-semibold">Total Amount</th>
                <th className="py-2.5 px-4 font-semibold">Status</th>
                <th className="py-2.5 px-4 font-semibold text-right">Statement</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {invoices.length > 0 ? (
                invoices.map((inv) => (
                  <tr key={inv.id}>
                    <td className="py-3 px-4 font-mono font-bold text-slate-900">{inv.invoiceNumber}</td>
                    <td className="py-3 px-4 text-xs text-slate-500">
                      {new Date(inv.createdAt).toLocaleDateString()}
                    </td>
                    <td className="py-3 px-4 text-xs text-slate-600">
                      {new Date(inv.billingPeriodStart).toLocaleDateString()} &rarr; {new Date(inv.billingPeriodEnd).toLocaleDateString()}
                    </td>
                    <td className="py-3 px-4 font-bold text-slate-900">{formatCurrency(inv.totalCents)}</td>
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
                  <td colSpan={6} className="py-8 text-center text-xs text-slate-400">
                    No billing statements issued yet.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </CardContent>
      </Card>

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

              <table className="w-full text-left text-xs border-t border-b">
                <thead className="bg-slate-50 text-slate-600">
                  <tr>
                    <th className="py-2 px-3">Description</th>
                    <th className="py-2 px-3 text-center">Qty</th>
                    <th className="py-2 px-3 text-right">Tax</th>
                    <th className="py-2 px-3 text-right">Total</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  <tr>
                    <td className="py-2.5 px-3 font-medium">Software Subscription Services</td>
                    <td className="py-2.5 px-3 text-center">1</td>
                    <td className="py-2.5 px-3 text-right">{formatCurrency(selectedInvoice.taxCents)}</td>
                    <td className="py-2.5 px-3 text-right font-bold">{formatCurrency(selectedInvoice.totalCents)}</td>
                  </tr>
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
