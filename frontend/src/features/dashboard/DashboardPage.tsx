import { useEffect, useState } from 'react';
import { useAuth } from '@/features/auth/AuthContext';
import { getDashboardMetrics } from '@/lib/api/modules';
import { DashboardMetrics } from '@/types';
import { formatCurrency } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Repeat,
  DollarSign,
  Receipt,
  CreditCard,
  AlertTriangle,
  ArrowUpRight,
  TrendingUp,
  RefreshCw,
  PlusCircle,
} from 'lucide-react';
import { Link } from 'react-router-dom';
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

export function DashboardPage() {
  const { user } = useAuth();
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const fetchMetrics = async () => {
    setIsLoading(true);
    try {
      const res = await getDashboardMetrics();
      setMetrics(res.metrics);
    } catch {
      // fallback
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
      {/* Top Banner */}
      <div className="bg-white border border-slate-200 rounded-lg p-6 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-slate-900">Brightloop Operations Center</h1>
            <Badge variant="indigo">{user?.role}</Badge>
          </div>
          <p className="text-sm text-slate-600">
            Real-time subscription lifecycles, recurring revenue health, and billing reliability.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button variant="outline" size="sm" onClick={fetchMetrics} isLoading={isLoading}>
            <RefreshCw className="h-4 w-4 mr-1.5 text-slate-500" />
            Refresh
          </Button>
          <Link to="/admin/subscriptions">
            <Button variant="primary" size="sm">
              <PlusCircle className="h-4 w-4 mr-1.5" />
              New Subscription
            </Button>
          </Link>
        </div>
      </div>

      {/* KPI Stat Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-4">
        <Card className="border-slate-200">
          <CardHeader className="pb-1">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Active Subs</span>
              <Repeat className="h-4 w-4 text-indigo-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">
              {metrics ? metrics.activeSubscriptionsCount : '—'}
            </div>
            <p className="text-[11px] text-slate-500 mt-1">Contracted recurring</p>
          </CardContent>
        </Card>

        <Card className="border-slate-200">
          <CardHeader className="pb-1">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">MRR</span>
              <DollarSign className="h-4 w-4 text-emerald-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-emerald-700">
              {metrics ? formatCurrency(metrics.mrrCents) : '—'}
            </div>
            <p className="text-[11px] text-slate-500 mt-1">Normalized monthly</p>
          </CardContent>
        </Card>

        <Card className="border-slate-200">
          <CardHeader className="pb-1">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Invoiced</span>
              <Receipt className="h-4 w-4 text-indigo-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">
              {metrics ? formatCurrency(metrics.invoicedRevenueCents) : '—'}
            </div>
            <p className="text-[11px] text-slate-500 mt-1">Total statements issued</p>
          </CardContent>
        </Card>

        <Card className="border-slate-200">
          <CardHeader className="pb-1">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Cash Collected</span>
              <CreditCard className="h-4 w-4 text-indigo-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">
              {metrics ? formatCurrency(metrics.collectedPaymentsCents) : '—'}
            </div>
            <p className="text-[11px] text-slate-500 mt-1">Recorded payments</p>
          </CardContent>
        </Card>

        <Card className="border-slate-200">
          <CardHeader className="pb-1">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Receivables</span>
              <TrendingUp className="h-4 w-4 text-amber-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-amber-700">
              {metrics ? formatCurrency(metrics.outstandingReceivablesCents) : '—'}
            </div>
            <p className="text-[11px] text-slate-500 mt-1">Outstanding balances</p>
          </CardContent>
        </Card>

        <Card className="border-slate-200">
          <CardHeader className="pb-1">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Overdue</span>
              <AlertTriangle className="h-4 w-4 text-rose-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-rose-700">
              {metrics ? metrics.overdueInvoicesCount : '—'}
            </div>
            <p className="text-[11px] text-slate-500 mt-1">Statements past due</p>
          </CardContent>
        </Card>
      </div>

      {/* Chart Section */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <div>
            <CardTitle className="text-base">Revenue Execution Trend</CardTitle>
            <p className="text-xs text-slate-500 mt-0.5">
              Comparison of Invoiced Billings vs Collected Payments across recent months
            </p>
          </div>
          <Link to="/admin/billing" className="text-xs text-indigo-600 hover:underline flex items-center gap-1">
            Billing Reliability Center <ArrowUpRight className="h-3.5 w-3.5" />
          </Link>
        </CardHeader>
        <CardContent>
          <div className="h-72 w-full pt-4">
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
                  <Bar dataKey="Invoiced" fill="#6366f1" radius={[4, 4, 0, 0]} />
                  <Bar dataKey="Collected" fill="#10b981" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="h-full flex items-center justify-center text-xs text-slate-400">
                Generating trend data...
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Quick Launchpad */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Link to="/admin/billing" className="group">
          <Card className="h-full hover:border-indigo-300 transition-colors">
            <CardHeader>
              <div className="flex items-center justify-between">
                <Badge variant="indigo">Differentiator</Badge>
                <ArrowUpRight className="h-4 w-4 text-slate-400 group-hover:text-indigo-600 transition-colors" />
              </div>
              <CardTitle className="text-base mt-2">Billing Reliability Center</CardTitle>
              <p className="text-xs text-slate-500 mt-1">
                Trigger automated batch invoicing, monitor due subscriptions, review duplicate prevention, and inspect explainable calculations.
              </p>
            </CardHeader>
          </Card>
        </Link>

        <Link to="/admin/quotations" className="group">
          <Card className="h-full hover:border-indigo-300 transition-colors">
            <CardHeader>
              <div className="flex items-center justify-between">
                <Badge variant="default">Quotations</Badge>
                <ArrowUpRight className="h-4 w-4 text-slate-400 group-hover:text-indigo-600 transition-colors" />
              </div>
              <CardTitle className="text-base mt-2">Proposals & Conversions</CardTitle>
              <p className="text-xs text-slate-500 mt-1">
                Draft client proposals with preconfigured plan and product lines. Convert accepted quotes directly into active subscriptions with 1 click.
              </p>
            </CardHeader>
          </Card>
        </Link>

        <Link to="/admin/catalog" className="group">
          <Card className="h-full hover:border-indigo-300 transition-colors">
            <CardHeader>
              <div className="flex items-center justify-between">
                <Badge variant="default">Catalog</Badge>
                <ArrowUpRight className="h-4 w-4 text-slate-400 group-hover:text-indigo-600 transition-colors" />
              </div>
              <CardTitle className="text-base mt-2">Products, Variants & Plans</CardTitle>
              <p className="text-xs text-slate-500 mt-1">
                Configure Brightloop Core software products, price variants (e.g. +₹560 Odoo connector), tax rates (18% GST), and discount rules.
              </p>
            </CardHeader>
          </Card>
        </Link>
      </div>
    </div>
  );
}
