import { useEffect, useState } from 'react';
import {
  getSubscriptions,
  createSubscription,
  activateSubscription,
  pauseSubscription,
  resumeSubscription,
  cancelSubscription,
  toggleAddon,
  getSubscriptionTimeline,
  getPlans,
  getProducts,
} from '@/lib/api/modules';
import { Subscription, SubscriptionEvent, Plan, Product } from '@/types';
import { formatCurrency } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert } from '@/components/ui/alert';
import {
  Plus,
  RefreshCw,
  Play,
  Pause,
  XCircle,
  Clock,
  Layers,
  Calendar,
} from 'lucide-react';

export function SubscriptionsPage() {
  const [subscriptions, setSubscriptions] = useState<Subscription[]>([]);
  const [plans, setPlans] = useState<Plan[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  // Timeline State
  const [timelineSubId, setTimelineSubId] = useState<string | null>(null);
  const [events, setEvents] = useState<SubscriptionEvent[]>([]);
  const [loadingTimeline, setLoadingTimeline] = useState(false);

  // Addon Modal State
  const [addonSub, setAddonSub] = useState<Subscription | null>(null);

  // Form State
  const [planId, setPlanId] = useState('');
  const [productId, setProductId] = useState('');
  const [quantity, setQuantity] = useState('1');
  const [customerEmail, setCustomerEmail] = useState('customer@brightloop.local');

  const loadData = async () => {
    setIsLoading(true);
    try {
      const [subsRes, plansRes, prodsRes] = await Promise.all([
        getSubscriptions(),
        getPlans(),
        getProducts(true),
      ]);
      setSubscriptions(subsRes.subscriptions || []);
      setPlans(plansRes.plans || []);
      setProducts(prodsRes.products || []);
      if (plansRes.plans && plansRes.plans.length > 0 && !planId) {
        setPlanId(plansRes.plans[0].id);
      }
      if (prodsRes.products && prodsRes.products.length > 0 && !productId) {
        setProductId(prodsRes.products[0].id);
      }
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const chosenProduct = products.find((p) => p.id === productId);
      const unitPrice = chosenProduct ? chosenProduct.salesPriceCents : 499900;

      await createSubscription({
        planId,
        paymentTerms: 'NET_30',
        lines: [
          {
            productId,
            isAddon: false,
            quantity: parseInt(quantity, 10),
            unitPriceCents: unitPrice,
            taxRateBps: 1800,
          },
        ],
      });
      setShowCreateModal(false);
      setMessage('Subscription provisioned successfully');
      loadData();
    } catch (err: any) {
      alert(err.message || 'Failed provisioning subscription');
    }
  };

  const handleActivate = async (id: string) => {
    try {
      await activateSubscription(id);
      setMessage('Subscription activated');
      loadData();
    } catch (err: any) {
      alert(err.message || 'Failed activation');
    }
  };

  const handlePause = async (id: string) => {
    try {
      await pauseSubscription(id);
      setMessage('Subscription paused');
      loadData();
    } catch (err: any) {
      alert(err.message || 'Failed pause');
    }
  };

  const handleResume = async (id: string) => {
    try {
      await resumeSubscription(id);
      setMessage('Subscription resumed');
      loadData();
    } catch (err: any) {
      alert(err.message || 'Failed resume');
    }
  };

  const handleCancel = async (id: string) => {
    if (!confirm('Are you sure you want to cancel this subscription?')) return;
    try {
      await cancelSubscription(id);
      setMessage('Subscription cancelled');
      loadData();
    } catch (err: any) {
      alert(err.message || 'Failed cancellation');
    }
  };

  const openTimeline = async (sub: Subscription) => {
    setTimelineSubId(sub.subscriptionNumber);
    setLoadingTimeline(true);
    try {
      const res = await getSubscriptionTimeline(sub.id);
      setEvents(res.events || []);
    } catch {
      setEvents([]);
    } finally {
      setLoadingTimeline(false);
    }
  };

  const handleToggleAddon = async (subId: string, pId: string, currentStatus: boolean) => {
    try {
      await toggleAddon(subId, pId, !currentStatus);
      setMessage(`Addon modified for next billing cycle`);
      loadData();
      setAddonSub(null);
    } catch (err: any) {
      alert(err.message || 'Failed toggling addon');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">Subscription Lifecycles</h1>
          <p className="text-sm text-slate-500 mt-1">
            Manage contract state machine transitions (DRAFT &rarr; ACTIVE &rarr; CLOSED), add-ons, and audit timelines.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button variant="outline" size="sm" onClick={loadData} isLoading={isLoading}>
            <RefreshCw className="h-4 w-4 mr-1 text-slate-500" /> Refresh
          </Button>
          <Button variant="primary" size="sm" onClick={() => setShowCreateModal(true)}>
            <Plus className="h-4 w-4 mr-1.5" /> Provision Subscription
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
          <CardTitle className="text-base">Active Contracts & Accounts</CardTitle>
        </CardHeader>
        <CardContent>
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-50 text-slate-500 text-xs uppercase border-b">
              <tr>
                <th className="py-2.5 px-4 font-semibold">Sub #</th>
                <th className="py-2.5 px-4 font-semibold">Customer</th>
                <th className="py-2.5 px-4 font-semibold">Plan</th>
                <th className="py-2.5 px-4 font-semibold">Next Billing</th>
                <th className="py-2.5 px-4 font-semibold">Total Amount</th>
                <th className="py-2.5 px-4 font-semibold">Status</th>
                <th className="py-2.5 px-4 font-semibold text-right">Lifecycle Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {subscriptions.length > 0 ? (
                subscriptions.map((s) => (
                  <tr key={s.id}>
                    <td className="py-3 px-4 font-mono font-bold text-slate-900">{s.subscriptionNumber}</td>
                    <td className="py-3 px-4 text-slate-700">{s.customerEmail || 'Customer'}</td>
                    <td className="py-3 px-4 text-slate-700">{s.planName || 'Plan'}</td>
                    <td className="py-3 px-4 text-slate-600">
                      <div className="flex items-center gap-1.5">
                        <Calendar className="h-3.5 w-3.5 text-slate-400" />
                        <span>{s.nextBillingDate ? new Date(s.nextBillingDate).toLocaleDateString() : '—'}</span>
                      </div>
                    </td>
                    <td className="py-3 px-4 font-bold text-slate-900">
                      {formatCurrency(s.totalAmountCents || 0)}
                    </td>
                    <td className="py-3 px-4">
                      <Badge
                        variant={
                          s.status === 'ACTIVE'
                            ? 'success'
                            : s.status === 'PAUSED'
                            ? 'warning'
                            : s.status === 'CONFIRMED'
                            ? 'indigo'
                            : 'default'
                        }
                      >
                        {s.status}
                      </Badge>
                    </td>
                    <td className="py-3 px-4 text-right">
                      <div className="flex items-center justify-end gap-1.5">
                        {/* Status Transition buttons */}
                        {s.status === 'CONFIRMED' && (
                          <Button variant="primary" size="sm" onClick={() => handleActivate(s.id)} className="text-xs">
                            <Play className="h-3 w-3 mr-1" /> Activate
                          </Button>
                        )}
                        {s.status === 'ACTIVE' && (
                          <Button variant="secondary" size="sm" onClick={() => handlePause(s.id)} className="text-xs">
                            <Pause className="h-3 w-3 mr-1" /> Pause
                          </Button>
                        )}
                        {s.status === 'PAUSED' && (
                          <Button variant="primary" size="sm" onClick={() => handleResume(s.id)} className="text-xs">
                            <Play className="h-3 w-3 mr-1" /> Resume
                          </Button>
                        )}
                        {s.status !== 'CANCELLED' && s.status !== 'CLOSED' && (
                          <Button variant="ghost" size="sm" onClick={() => handleCancel(s.id)} className="text-xs text-rose-600">
                            <XCircle className="h-3 w-3" />
                          </Button>
                        )}

                        <Button variant="outline" size="sm" onClick={() => setAddonSub(s)} className="text-xs">
                          <Layers className="h-3 w-3 mr-1" /> Add-ons
                        </Button>

                        <Button variant="ghost" size="sm" onClick={() => openTimeline(s)} className="text-xs text-indigo-600">
                          <Clock className="h-3.5 w-3.5" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={7} className="py-8 text-center text-xs text-slate-400">
                    No subscriptions found. Click 'Provision Subscription' or convert an accepted quotation.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </CardContent>
      </Card>

      {/* Modal: Provision Subscription */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-md w-full shadow-lg">
            <CardHeader>
              <CardTitle className="text-lg">Provision Subscription</CardTitle>
            </CardHeader>
            <form onSubmit={handleCreate}>
              <CardContent className="space-y-4">
                <div>
                  <Label>Customer Email</Label>
                  <Input required value={customerEmail} onChange={(e) => setCustomerEmail(e.target.value)} />
                </div>
                <div>
                  <Label>Plan</Label>
                  <select
                    className="w-full border rounded px-3 py-2 text-sm bg-white"
                    value={planId}
                    onChange={(e) => setPlanId(e.target.value)}
                  >
                    {plans.map((p) => (
                      <option key={p.id} value={p.id}>
                        {p.name} ({p.frequency} &bull; {formatCurrency(p.billingPriceCents)})
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <Label>Main Service Product</Label>
                  <select
                    className="w-full border rounded px-3 py-2 text-sm bg-white"
                    value={productId}
                    onChange={(e) => setProductId(e.target.value)}
                  >
                    {products.map((p) => (
                      <option key={p.id} value={p.id}>
                        {p.name} ({formatCurrency(p.salesPriceCents)})
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <Label>Quantity</Label>
                  <Input type="number" min="1" required value={quantity} onChange={(e) => setQuantity(e.target.value)} />
                </div>
              </CardContent>
              <div className="p-4 border-t flex justify-end gap-2 bg-slate-50 rounded-b-lg">
                <Button type="button" variant="outline" onClick={() => setShowCreateModal(false)}>
                  Cancel
                </Button>
                <Button type="submit" variant="primary">
                  Provision Now
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}

      {/* Modal: Timeline */}
      {timelineSubId && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-lg w-full max-h-[85vh] flex flex-col shadow-xl">
            <CardHeader className="border-b">
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="text-base">Subscription Timeline</CardTitle>
                  <p className="text-xs text-slate-500 font-mono mt-0.5">{timelineSubId}</p>
                </div>
                <Button variant="ghost" size="sm" onClick={() => setTimelineSubId(null)}>
                  Close
                </Button>
              </div>
            </CardHeader>
            <CardContent className="overflow-y-auto p-6 space-y-4">
              {loadingTimeline ? (
                <div className="text-center py-6 text-sm text-slate-500">Loading timeline events...</div>
              ) : events.length > 0 ? (
                <div className="relative border-l border-indigo-200 ml-4 space-y-6">
                  {events.map((ev) => (
                    <div key={ev.id} className="relative pl-6">
                      <div className="absolute -left-2 top-0 h-4 w-4 rounded-full bg-indigo-600 ring-4 ring-white" />
                      <div className="flex items-center justify-between">
                        <span className="text-xs font-bold text-slate-900 uppercase tracking-wider">{ev.eventType}</span>
                        <span className="text-[11px] text-slate-400">
                          {new Date(ev.createdAt).toLocaleString()}
                        </span>
                      </div>
                      <p className="text-xs text-slate-500 mt-0.5">Actor: {ev.actorEmail || 'System Worker'}</p>
                      {ev.payload && (
                        <pre className="bg-slate-50 p-2 rounded text-[11px] font-mono text-slate-700 mt-1 border overflow-x-auto">
                          {JSON.stringify(ev.payload, null, 2)}
                        </pre>
                      )}
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-6 text-xs text-slate-400">No events logged for this subscription yet.</div>
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {/* Modal: Addons */}
      {addonSub && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-md w-full shadow-lg">
            <CardHeader className="border-b">
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="text-base">Permitted Add-on Services</CardTitle>
                  <p className="text-xs text-slate-500">{addonSub.subscriptionNumber}</p>
                </div>
                <Button variant="ghost" size="sm" onClick={() => setAddonSub(null)}>
                  Close
                </Button>
              </div>
            </CardHeader>
            <CardContent className="p-6 space-y-4">
              <p className="text-xs text-slate-500">
                Add-on adjustments take effect on the <strong>next billing cycle</strong> without immediate proration.
              </p>
              <div className="space-y-3">
                {products.map((prod) => {
                  const isCurrent = addonSub.lines?.some((l) => l.productId === prod.id && l.isAddon);
                  return (
                    <div
                      key={prod.id}
                      className="flex items-center justify-between p-3 border rounded-md bg-slate-50"
                    >
                      <div>
                        <div className="text-xs font-semibold text-slate-900">{prod.name}</div>
                        <div className="text-[11px] text-slate-500">{formatCurrency(prod.salesPriceCents)}</div>
                      </div>
                      <Button
                        variant={isCurrent ? 'danger' : 'primary'}
                        size="sm"
                        onClick={() => handleToggleAddon(addonSub.id, prod.id, !!isCurrent)}
                        className="text-xs"
                      >
                        {isCurrent ? 'Deactivate' : 'Enable Add-on'}
                      </Button>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}
