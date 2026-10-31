import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Receipt } from 'lucide-react';

export function CustomerInvoicesPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-900">Invoices & Receipts</h1>
        <p className="text-sm text-slate-500 mt-1">
          View past billing statements, payment receipts, and download PDF copies.
        </p>
      </div>

      <Card className="border-dashed border-2 border-slate-200 bg-white">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center">
              <Receipt className="h-5 w-5" />
            </div>
            <div>
              <CardTitle className="text-base">No Invoices Issued</CardTitle>
              <CardDescription>
                When recurring billing cycles execute, your itemized statements will appear here with explainable price breakdowns.
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <p className="text-xs text-slate-500">
            PDF invoices and payment receipt downloads will be enabled once Member 3's billing engine branch is integrated.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
