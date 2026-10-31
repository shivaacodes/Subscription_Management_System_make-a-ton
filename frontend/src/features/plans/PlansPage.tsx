import { useEffect, useState } from 'react';
import { getPlans, createPlan } from '@/lib/api/modules';
import { Plan, PlanFrequency } from '@/types';
import { formatCurrency } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert } from '@/components/ui/alert';
import { Plus, RefreshCw, CheckCircle2, XCircle } from 'lucide-react';

export function PlansPage() {
  const [plans, setPlans] = useState<Plan[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  // Form State
  const [name, setName] = useState('');
  const [frequency, setFrequency] = useState<PlanFrequency>('MONTHLY');
  const [price, setPrice] = useState('4999');
  const [minQty, setMinQty] = useState('1');
  const [allowPause, setAllowPause] = useState(true);
  const [allowEarlyClose, setAllowEarlyClose] = useState(true);
  const [allowRenewal, setAllowRenewal] = useState(true);

  const loadPlans = async () => {
    setIsLoading(true);
    try {
      const res = await getPlans();
      setPlans(res.plans || []);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadPlans();
  }, []);

  const handleCreatePlan = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createPlan({
        name,
        frequency,
        billingPriceCents: Math.round(parseFloat(price) * 100),
        minQuantity: parseInt(minQty, 10),
        allowPause,
        allowEarlyClose,
        allowRenewal,
        autoClose: false,
        isActive: true,
      });
      setShowModal(false);
      setName('');
      setMessage('Subscription plan created successfully');
      loadPlans();
    } catch (err: any) {
      alert(err.message || 'Failed creating plan');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">Subscription Plans</h1>
          <p className="text-sm text-slate-500 mt-1">
            Configure recurring billing cadences, availability windows, and subscriber pause/renewal permissions.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button variant="outline" size="sm" onClick={loadPlans} isLoading={isLoading}>
            <RefreshCw className="h-4 w-4 mr-1 text-slate-500" /> Refresh
          </Button>
          <Button variant="primary" size="sm" onClick={() => setShowModal(true)}>
            <Plus className="h-4 w-4 mr-1.5" /> New Subscription Plan
          </Button>
        </div>
      </div>

      {message && (
        <Alert variant="success" className="cursor-pointer" onClick={() => setMessage(null)}>
          {message}
        </Alert>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {plans.map((p) => (
          <Card key={p.id} className="border-slate-200 flex flex-col justify-between">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between mb-2">
                <Badge variant="indigo" className="font-mono">{p.frequency}</Badge>
                <Badge variant={p.isActive ? 'success' : 'default'}>
                  {p.isActive ? 'Active' : 'Archived'}
                </Badge>
              </div>
              <CardTitle className="text-lg text-slate-900">{p.name}</CardTitle>
              <div className="text-2xl font-bold text-slate-900 mt-2">
                {formatCurrency(p.billingPriceCents)}
                <span className="text-xs font-normal text-slate-500 ml-1">/ {p.frequency.toLowerCase()}</span>
              </div>
            </CardHeader>
            <CardContent className="pt-0 space-y-3">
              <div className="border-t border-slate-100 pt-3 space-y-2 text-xs text-slate-600">
                <div className="flex items-center justify-between">
                  <span>Minimum Quantity:</span>
                  <span className="font-semibold">{p.minQuantity} unit(s)</span>
                </div>
                <div className="flex items-center justify-between">
                  <span>Subscriber Pause:</span>
                  <span className="flex items-center gap-1 font-medium">
                    {p.allowPause ? (
                      <CheckCircle2 className="h-3.5 w-3.5 text-emerald-600" />
                    ) : (
                      <XCircle className="h-3.5 w-3.5 text-slate-400" />
                    )}
                    {p.allowPause ? 'Permitted' : 'Disabled'}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span>Early Cancellation:</span>
                  <span className="flex items-center gap-1 font-medium">
                    {p.allowEarlyClose ? (
                      <CheckCircle2 className="h-3.5 w-3.5 text-emerald-600" />
                    ) : (
                      <XCircle className="h-3.5 w-3.5 text-slate-400" />
                    )}
                    {p.allowEarlyClose ? 'Allowed' : 'Locked'}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span>Automatic Renewal:</span>
                  <span className="flex items-center gap-1 font-medium">
                    {p.allowRenewal ? (
                      <CheckCircle2 className="h-3.5 w-3.5 text-emerald-600" />
                    ) : (
                      <XCircle className="h-3.5 w-3.5 text-slate-400" />
                    )}
                    {p.allowRenewal ? 'Supported' : 'No Renewal'}
                  </span>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Modal: New Plan */}
      {showModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-md w-full shadow-lg">
            <CardHeader>
              <CardTitle className="text-lg">Create Subscription Plan</CardTitle>
            </CardHeader>
            <form onSubmit={handleCreatePlan}>
              <CardContent className="space-y-4">
                <div>
                  <Label>Plan Name</Label>
                  <Input
                    required
                    placeholder="e.g. Standard Annual Plan"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                  />
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <Label>Billing Frequency</Label>
                    <select
                      className="w-full border rounded px-3 py-2 text-sm bg-white"
                      value={frequency}
                      onChange={(e: any) => setFrequency(e.target.value)}
                    >
                      <option value="DAILY">Daily</option>
                      <option value="WEEKLY">Weekly</option>
                      <option value="MONTHLY">Monthly</option>
                      <option value="YEARLY">Yearly</option>
                    </select>
                  </div>
                  <div>
                    <Label>Billing Price (₹)</Label>
                    <Input
                      type="number"
                      required
                      value={price}
                      onChange={(e) => setPrice(e.target.value)}
                    />
                  </div>
                </div>

                <div>
                  <Label>Minimum Quantity</Label>
                  <Input
                    type="number"
                    min="1"
                    required
                    value={minQty}
                    onChange={(e) => setMinQty(e.target.value)}
                  />
                </div>

                <div className="border-t pt-3 space-y-2 text-xs">
                  <div className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      id="allowPause"
                      checked={allowPause}
                      onChange={(e) => setAllowPause(e.target.checked)}
                      className="rounded text-indigo-600"
                    />
                    <label htmlFor="allowPause" className="text-slate-700 font-medium">
                      Allow customers to pause subscription
                    </label>
                  </div>

                  <div className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      id="allowEarlyClose"
                      checked={allowEarlyClose}
                      onChange={(e) => setAllowEarlyClose(e.target.checked)}
                      className="rounded text-indigo-600"
                    />
                    <label htmlFor="allowEarlyClose" className="text-slate-700 font-medium">
                      Allow early cancellation before expiration
                    </label>
                  </div>

                  <div className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      id="allowRenewal"
                      checked={allowRenewal}
                      onChange={(e) => setAllowRenewal(e.target.checked)}
                      className="rounded text-indigo-600"
                    />
                    <label htmlFor="allowRenewal" className="text-slate-700 font-medium">
                      Permit automated recurring renewals
                    </label>
                  </div>
                </div>
              </CardContent>
              <div className="p-4 border-t flex justify-end gap-2 bg-slate-50 rounded-b-lg">
                <Button type="button" variant="outline" onClick={() => setShowModal(false)}>
                  Cancel
                </Button>
                <Button type="submit" variant="primary">
                  Save Plan
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  );
}
