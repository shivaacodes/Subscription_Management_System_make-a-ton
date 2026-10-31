import { useEffect, useState } from 'react';
import { getDashboardMetrics } from '@/lib/api/modules';
import { DashboardMetrics } from '@/types';
import { formatCurrency } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { RefreshCw, ShieldCheck } from 'lucide-react';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  CartesianGrid,
} from 'recharts';

export function ReportingPage() {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const fetchMetrics = async () => {
    setIsLoading(true);
    try {
      const res = await getDashboardMetrics();
      setMetrics(res.metrics);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchMetrics();
  }, []);

  const chartData = (metrics?.trend || []).map((t) => ({
    name: t.month,
    Invoiced: t.invoicedRevenueCents / 100,
    Collected: t.collectedCents / 100,
  }));

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">Revenue & Operations Reporting</h1>
          <p className="text-sm text-slate-500 mt-1">
            Authoritative accounting metrics with strict operational distinction between MRR, Invoiced Revenue, and Cash Collected.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button variant="outline" size="sm" onClick={fetchMetrics} isLoading={isLoading}>
            <RefreshCw className="h-4 w-4 mr-1 text-slate-500" /> Refresh Metrics
          </Button>
        </div>
      </div>

      {/* Accounting Clarification Alert */}
      <div className="p-4 bg-indigo-50/70 border border-indigo-200 rounded-md text-xs text-indigo-950 flex items-start gap-3">
        <ShieldCheck className="h-5 w-5 text-indigo-600 flex-shrink-0 mt-0.5" />
        <div>
          <span className="font-bold">Accounting Separation Guarantee:</span> Veyra calculates Monthly Recurring Revenue (MRR) based exclusively on active contracted normalized run rates. Invoices issued represent legally binding accounting debits, while Collected Payments represent real cash settled in bank accounts.
        </div>
      </div>

      {/* Metric Breakdown Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card className="border-emerald-200 bg-emerald-50/20">
          <CardHeader className="pb-2">
            <span className="text-xs font-semibold text-emerald-800 uppercase tracking-wider">Monthly Recurring Revenue</span>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-extrabold text-emerald-800">
              {metrics ? formatCurrency(metrics.mrrCents) : '—'}
            </div>
            <p className="text-xs text-emerald-700 mt-1">
              From {metrics?.activeSubscriptionsCount || 0} active subscriptions
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Cumulative Invoiced</span>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-extrabold text-slate-900">
              {metrics ? formatCurrency(metrics.invoicedRevenueCents) : '—'}
            </div>
            <p className="text-xs text-slate-500 mt-1">Total recognized billing statements</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Cash Collected</span>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-extrabold text-slate-900">
              {metrics ? formatCurrency(metrics.collectedPaymentsCents) : '—'}
            </div>
            <p className="text-xs text-slate-500 mt-1">Confirmed payments settled</p>
          </CardContent>
        </Card>

        <Card className="border-amber-200 bg-amber-50/20">
          <CardHeader className="pb-2">
            <span className="text-xs font-semibold text-amber-800 uppercase tracking-wider">Accounts Receivable</span>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-extrabold text-amber-800">
              {metrics ? formatCurrency(metrics.outstandingReceivablesCents) : '—'}
            </div>
            <p className="text-xs text-amber-700 mt-1">
              {metrics?.overdueInvoicesCount || 0} invoice(s) past due
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Visual Analytics Charts */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Invoicing vs Cash Inflow Comparison</CardTitle>
          <p className="text-xs text-slate-500">Historical monthly billing velocity</p>
        </CardHeader>
        <CardContent>
          <div className="h-80 w-full pt-4">
            {chartData.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={chartData} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f1f5f9" />
                  <XAxis dataKey="name" stroke="#94a3b8" fontSize={12} tickLine={false} />
                  <YAxis stroke="#94a3b8" fontSize={12} tickLine={false} tickFormatter={(val) => `₹${val}`} />
                  <Tooltip
                    formatter={(value: any) => [`₹${Number(value).toLocaleString('en-IN')}`, '']}
                    contentStyle={{ backgroundColor: '#ffffff', borderRadius: '6px', border: '1px solid #e2e8f0' }}
                  />
                  <Legend />
                  <Bar dataKey="Invoiced" fill="#6366f1" radius={[4, 4, 0, 0]} name="Invoiced Revenue" />
                  <Bar dataKey="Collected" fill="#10b981" radius={[4, 4, 0, 0]} name="Cash Collected" />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="h-full flex items-center justify-center text-xs text-slate-400">
                Awaiting historical billing runs...
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
