import { useAuth } from '@/features/auth/AuthContext';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Layers, Database, ShieldCheck, Users, ArrowUpRight } from 'lucide-react';
import { Link } from 'react-router-dom';

export function DashboardPage() {
  const { user } = useAuth();

  return (
    <div className="space-y-8">
      {/* Welcome Banner */}
      <div className="bg-white border border-slate-200 rounded-lg p-6 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-slate-900">Welcome to Veyra</h1>
            <Badge variant="indigo">{user?.role}</Badge>
          </div>
          <p className="text-sm text-slate-600">
            Unified Subscription Management and Recurring Billing for Brightloop Software.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Link
            to="/admin/system"
            className="inline-flex items-center gap-1.5 px-3 py-2 bg-slate-900 text-white rounded-md text-xs font-medium hover:bg-slate-800 transition-colors shadow-sm"
          >
            <ShieldCheck className="h-4 w-4" />
            Verify System & Auth
          </Link>
        </div>
      </div>

      {/* Foundation Status Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card>
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <CardTitle className="text-sm font-medium text-slate-500 uppercase tracking-wider">
                Auth & Identity
              </CardTitle>
              <Users className="h-5 w-5 text-indigo-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold text-slate-900">Active & Enforced</div>
            <p className="text-xs text-slate-500 mt-1">
              Role-based session access for ADMIN, STAFF, and CUSTOMER with bcrypt hashing.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <CardTitle className="text-sm font-medium text-slate-500 uppercase tracking-wider">
                Database Engine
              </CardTitle>
              <Database className="h-5 w-5 text-indigo-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold text-slate-900">PostgreSQL Schema</div>
            <p className="text-xs text-slate-500 mt-1">
              Normalized tables for catalog, plans, quotations, subscriptions, billing, and audit.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <CardTitle className="text-sm font-medium text-slate-500 uppercase tracking-wider">
                Architecture
              </CardTitle>
              <Layers className="h-5 w-5 text-indigo-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold text-slate-900">Modular Monolith</div>
            <p className="text-xs text-slate-500 mt-1">
              Handler &rarr; Service &rarr; Repository with 0 circular dependencies.
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Feature Tracks */}
      <div className="space-y-4">
        <h2 className="text-base font-semibold text-slate-900">Feature Branch Tracks (4 Members)</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Link to="/admin/subscriptions" className="group">
            <Card className="h-full hover:border-indigo-300 transition-colors">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <Badge variant="indigo">Member 1 &bull; Team Lead</Badge>
                  <ArrowUpRight className="h-4 w-4 text-slate-400 group-hover:text-indigo-600 transition-colors" />
                </div>
                <CardTitle className="text-base mt-2">feat/subscription-lifecycle</CardTitle>
                <CardDescription>
                  Subscription CRUD, state machine (DRAFT &rarr; ACTIVE &rarr; CLOSED), renewal policies, add-on toggling, and audit timeline.
                </CardDescription>
              </CardHeader>
            </Card>
          </Link>

          <Link to="/admin/catalog" className="group">
            <Card className="h-full hover:border-indigo-300 transition-colors">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <Badge variant="default">Member 2</Badge>
                  <ArrowUpRight className="h-4 w-4 text-slate-400 group-hover:text-indigo-600 transition-colors" />
                </div>
                <CardTitle className="text-base mt-2">feat/catalog-plans</CardTitle>
                <CardDescription>
                  Products, variants (e.g. Odoo brand +₹560), frequency-based plans, tax rules, and administrative discounts.
                </CardDescription>
              </CardHeader>
            </Card>
          </Link>

          <Link to="/admin/billing" className="group">
            <Card className="h-full hover:border-indigo-300 transition-colors">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <Badge variant="default">Member 3</Badge>
                  <ArrowUpRight className="h-4 w-4 text-slate-400 group-hover:text-indigo-600 transition-colors" />
                </div>
                <CardTitle className="text-base mt-2">feat/billing-engine</CardTitle>
                <CardDescription>
                  Recurring billing worker, idempotent invoice generation, payment recording, and Billing Reliability Center.
                </CardDescription>
              </CardHeader>
            </Card>
          </Link>

          <Link to="/admin/quotations" className="group">
            <Card className="h-full hover:border-indigo-300 transition-colors">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <Badge variant="default">Member 4</Badge>
                  <ArrowUpRight className="h-4 w-4 text-slate-400 group-hover:text-indigo-600 transition-colors" />
                </div>
                <CardTitle className="text-base mt-2">feat/quotations-portal</CardTitle>
                <CardDescription>
                  Quotation templates, quotation conversion, Customer self-service portal, and executive revenue reporting.
                </CardDescription>
              </CardHeader>
            </Card>
          </Link>
        </div>
      </div>
    </div>
  );
}
