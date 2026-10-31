import { useEffect, useState } from 'react';
import { getProducts, createProduct, createVariant, getTaxes, createTax, getDiscounts, createDiscount } from '@/lib/api/modules';
import { Product, TaxRule, DiscountRule } from '@/types';
import { formatCurrency } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert } from '@/components/ui/alert';
import { Package, Plus, Percent, Receipt, RefreshCw, Layers } from 'lucide-react';

export function CatalogPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [taxes, setTaxes] = useState<TaxRule[]>([]);
  const [discounts, setDiscounts] = useState<DiscountRule[]>([]);
  const [tab, setTab] = useState<'products' | 'taxes' | 'discounts'>('products');
  const [isLoading, setIsLoading] = useState(true);

  // New Product Modal State
  const [showProductModal, setShowProductModal] = useState(false);
  const [newProductName, setNewProductName] = useState('');
  const [newProductSalesPrice, setNewProductSalesPrice] = useState('4999');
  const [newProductCostPrice, setNewProductCostPrice] = useState('1200');
  const [newProductNotes, setNewProductNotes] = useState('');

  // Add Variant Modal State
  const [selectedProductForVariant, setSelectedProductForVariant] = useState<Product | null>(null);
  const [variantName, setVariantName] = useState('');
  const [variantAdjustment, setVariantAdjustment] = useState('560');

  // New Tax Modal State
  const [showTaxModal, setShowTaxModal] = useState(false);
  const [taxName, setTaxName] = useState('GST Standard');
  const [taxRateBps, setTaxRateBps] = useState('1800');

  // New Discount Modal State
  const [showDiscountModal, setShowDiscountModal] = useState(false);
  const [discountCode, setDiscountCode] = useState('WELCOME10');
  const [discountType, setDiscountType] = useState<'PERCENTAGE' | 'FIXED_AMOUNT'>('PERCENTAGE');
  const [discountValue, setDiscountValue] = useState('1000'); // 10% in bps

  const [message, setMessage] = useState<string | null>(null);

  const loadData = async () => {
    setIsLoading(true);
    try {
      const [prodRes, taxRes, discRes] = await Promise.all([
        getProducts(false),
        getTaxes(),
        getDiscounts(),
      ]);
      setProducts(prodRes.products || []);
      setTaxes(taxRes.taxes || []);
      setDiscounts(discRes.discounts || []);
    } catch {
      // error handled gracefully
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleCreateProduct = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createProduct({
        name: newProductName,
        salesPriceCents: Math.round(parseFloat(newProductSalesPrice) * 100),
        costPriceCents: Math.round(parseFloat(newProductCostPrice) * 100),
        notes: newProductNotes || undefined,
        isActive: true,
      });
      setShowProductModal(false);
      setNewProductName('');
      setMessage('Product created successfully');
      loadData();
    } catch (err: any) {
      alert(err.message || 'Failed creating product');
    }
  };

  const handleCreateVariant = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProductForVariant) return;
    try {
      await createVariant(selectedProductForVariant.id, {
        name: variantName,
        priceAdjustmentCents: Math.round(parseFloat(variantAdjustment) * 100),
        isActive: true,
      });
      setSelectedProductForVariant(null);
      setVariantName('');
      setMessage(`Added variant to ${selectedProductForVariant.name}`);
      loadData();
    } catch (err: any) {
      alert(err.message || 'Failed creating variant');
    }
  };

  const handleCreateTax = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createTax({
        name: taxName,
        rateBps: parseInt(taxRateBps, 10),
        country: 'IN',
        isActive: true,
      });
      setShowTaxModal(false);
      setMessage('Tax rule created');
      loadData();
    } catch (err: any) {
      alert(err.message || 'Failed creating tax rule');
    }
  };

  const handleCreateDiscount = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createDiscount({
        code: discountCode,
        discountType: discountType,
        valueCents: parseInt(discountValue, 10),
        minPurchaseCents: 0,
        minQuantity: 1,
        appliesTo: 'SUBSCRIPTION',
        isActive: true,
      });
      setShowDiscountModal(false);
      setMessage('Discount rule created');
      loadData();
    } catch (err: any) {
      alert(err.message || 'Failed creating discount');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">Products & Catalog</h1>
          <p className="text-sm text-slate-500 mt-1">
            Manage software products, price adjustment variants (e.g., Odoo connector +₹560), and administrative billing rules.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button variant="outline" size="sm" onClick={loadData} isLoading={isLoading}>
            <RefreshCw className="h-4 w-4 mr-1 text-slate-500" /> Refresh
          </Button>
          {tab === 'products' && (
            <Button variant="primary" size="sm" onClick={() => setShowProductModal(true)}>
              <Plus className="h-4 w-4 mr-1.5" /> New Product
            </Button>
          )}
          {tab === 'taxes' && (
            <Button variant="primary" size="sm" onClick={() => setShowTaxModal(true)}>
              <Plus className="h-4 w-4 mr-1.5" /> New Tax Rule
            </Button>
          )}
          {tab === 'discounts' && (
            <Button variant="primary" size="sm" onClick={() => setShowDiscountModal(true)}>
              <Plus className="h-4 w-4 mr-1.5" /> New Discount Code
            </Button>
          )}
        </div>
      </div>

      {message && (
        <Alert variant="success" className="cursor-pointer" onClick={() => setMessage(null)}>
          {message}
        </Alert>
      )}

      {/* Navigation Tabs */}
      <div className="flex space-x-1 border-b border-slate-200">
        <button
          onClick={() => setTab('products')}
          className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
            tab === 'products'
              ? 'border-indigo-600 text-indigo-600'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          }`}
        >
          <div className="flex items-center gap-2">
            <Package className="h-4 w-4" />
            <span>Products & Variants ({products.length})</span>
          </div>
        </button>
        <button
          onClick={() => setTab('taxes')}
          className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
            tab === 'taxes'
              ? 'border-indigo-600 text-indigo-600'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          }`}
        >
          <div className="flex items-center gap-2">
            <Receipt className="h-4 w-4" />
            <span>Tax Rules ({taxes.length})</span>
          </div>
        </button>
        <button
          onClick={() => setTab('discounts')}
          className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
            tab === 'discounts'
              ? 'border-indigo-600 text-indigo-600'
              : 'border-transparent text-slate-500 hover:text-slate-800'
          }`}
        >
          <div className="flex items-center gap-2">
            <Percent className="h-4 w-4" />
            <span>Discount Rules ({discounts.length})</span>
          </div>
        </button>
      </div>

      {/* Tab: Products */}
      {tab === 'products' && (
        <div className="space-y-4">
          <div className="grid grid-cols-1 gap-4">
            {products.map((p) => (
              <Card key={p.id} className="border-slate-200">
                <CardHeader className="flex flex-row items-center justify-between py-4">
                  <div className="flex items-center gap-3">
                    <div className="h-9 w-9 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold">
                      <Package className="h-5 w-5" />
                    </div>
                    <div>
                      <CardTitle className="text-base text-slate-900">{p.name}</CardTitle>
                      <p className="text-xs text-slate-500">{p.productType} &bull; Cost: {formatCurrency(p.costPriceCents)}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <div className="text-right">
                      <div className="text-sm font-bold text-slate-900">{formatCurrency(p.salesPriceCents)}</div>
                      <Badge variant={p.isActive ? 'success' : 'default'} className="text-[10px]">
                        {p.isActive ? 'Active' : 'Inactive'}
                      </Badge>
                    </div>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setSelectedProductForVariant(p)}
                      className="text-xs"
                    >
                      <Plus className="h-3.5 w-3.5 mr-1" /> Add Variant
                    </Button>
                  </div>
                </CardHeader>
                <CardContent className="pt-0 pb-4">
                  {p.variants && p.variants.length > 0 ? (
                    <div className="mt-2 pl-12 space-y-1.5">
                      <div className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
                        Configurable Variants
                      </div>
                      <div className="flex flex-wrap gap-2">
                        {p.variants.map((v) => (
                          <div
                            key={v.id}
                            className="inline-flex items-center gap-2 px-3 py-1 bg-slate-50 rounded-md border border-slate-200 text-xs text-slate-800"
                          >
                            <Layers className="h-3.5 w-3.5 text-indigo-500" />
                            <span className="font-medium">{v.name}</span>
                            <span className="text-indigo-600 font-semibold">
                              {v.priceAdjustmentCents >= 0 ? '+' : ''}{formatCurrency(v.priceAdjustmentCents)}
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                  ) : (
                    <div className="pl-12 text-xs text-slate-400 italic">No variants configured for this base product yet.</div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}

      {/* Tab: Taxes */}
      {tab === 'taxes' && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Tax Configuration</CardTitle>
          </CardHeader>
          <CardContent>
            <table className="w-full text-left text-sm">
              <thead className="bg-slate-50 text-slate-500 text-xs uppercase border-b">
                <tr>
                  <th className="py-2.5 px-4 font-semibold">Tax Name</th>
                  <th className="py-2.5 px-4 font-semibold">Rate</th>
                  <th className="py-2.5 px-4 font-semibold">Country</th>
                  <th className="py-2.5 px-4 font-semibold">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {taxes.map((t) => (
                  <tr key={t.id}>
                    <td className="py-3 px-4 font-medium text-slate-900">{t.name}</td>
                    <td className="py-3 px-4 text-indigo-600 font-semibold">{(t.rateBps / 100).toFixed(2)}%</td>
                    <td className="py-3 px-4 text-slate-600">{t.country}</td>
                    <td className="py-3 px-4">
                      <Badge variant="success">Active</Badge>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      )}

      {/* Tab: Discounts */}
      {tab === 'discounts' && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Administrative Discount Rules</CardTitle>
          </CardHeader>
          <CardContent>
            <table className="w-full text-left text-sm">
              <thead className="bg-slate-50 text-slate-500 text-xs uppercase border-b">
                <tr>
                  <th className="py-2.5 px-4 font-semibold">Code</th>
                  <th className="py-2.5 px-4 font-semibold">Type</th>
                  <th className="py-2.5 px-4 font-semibold">Discount Value</th>
                  <th className="py-2.5 px-4 font-semibold">Usage</th>
                  <th className="py-2.5 px-4 font-semibold">Applies To</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {discounts.length > 0 ? (
                  discounts.map((d) => (
                    <tr key={d.id}>
                      <td className="py-3 px-4 font-mono font-bold text-indigo-600">{d.code}</td>
                      <td className="py-3 px-4 text-slate-600">{d.discountType}</td>
                      <td className="py-3 px-4 font-semibold">
                        {d.discountType === 'PERCENTAGE'
                          ? `${(d.valueCents / 100).toFixed(0)}%`
                          : formatCurrency(d.valueCents)}
                      </td>
                      <td className="py-3 px-4 text-slate-500">
                        {d.currentUsageCount} {d.maxUsageCount ? `/ ${d.maxUsageCount}` : 'unlimited'}
                      </td>
                      <td className="py-3 px-4">
                        <Badge variant="indigo">{d.appliesTo}</Badge>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={5} className="py-6 text-center text-xs text-slate-400">
                      No discount rules configured yet. Click 'New Discount Code' to create one.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </CardContent>
        </Card>
      )}

      {/* Modal: New Product */}
      {showProductModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-md w-full shadow-lg">
            <CardHeader>
              <CardTitle className="text-lg">Create New Product</CardTitle>
            </CardHeader>
            <form onSubmit={handleCreateProduct}>
              <CardContent className="space-y-4">
                <div>
                  <Label>Product Name</Label>
                  <Input
                    required
                    placeholder="e.g. Brightloop Analytics Suite"
                    value={newProductName}
                    onChange={(e) => setNewProductName(e.target.value)}
                  />
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <Label>Sales Price (₹)</Label>
                    <Input
                      type="number"
                      required
                      value={newProductSalesPrice}
                      onChange={(e) => setNewProductSalesPrice(e.target.value)}
                    />
                  </div>
                  <div>
                    <Label>Cost Price (₹)</Label>
                    <Input
                      type="number"
                      required
                      value={newProductCostPrice}
                      onChange={(e) => setNewProductCostPrice(e.target.value)}
                    />
                  </div>
                </div>
                <div>
                  <Label>Notes / Description</Label>
                  <Input
                    placeholder="Optional description"
                    value={newProductNotes}
                    onChange={(e) => setNewProductNotes(e.target.value)}
                  />
                </div>
              </CardContent>
              <div className="p-4 border-t flex justify-end gap-2 bg-slate-50 rounded-b-lg">
                <Button type="button" variant="outline" onClick={() => setShowProductModal(false)}>
                  Cancel
                </Button>
                <Button type="submit" variant="primary">
                  Save Product
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}

      {/* Modal: Add Variant */}
      {selectedProductForVariant && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-md w-full shadow-lg">
            <CardHeader>
              <CardTitle className="text-lg">Add Variant to {selectedProductForVariant.name}</CardTitle>
            </CardHeader>
            <form onSubmit={handleCreateVariant}>
              <CardContent className="space-y-4">
                <div>
                  <Label>Variant Name</Label>
                  <Input
                    required
                    placeholder="e.g. Odoo Brand Connector"
                    value={variantName}
                    onChange={(e) => setVariantName(e.target.value)}
                  />
                </div>
                <div>
                  <Label>Price Adjustment (₹)</Label>
                  <Input
                    type="number"
                    required
                    placeholder="560"
                    value={variantAdjustment}
                    onChange={(e) => setVariantAdjustment(e.target.value)}
                  />
                  <p className="text-[11px] text-slate-500 mt-1">
                    Adds this amount to the base product sales price (e.g., +₹560).
                  </p>
                </div>
              </CardContent>
              <div className="p-4 border-t flex justify-end gap-2 bg-slate-50 rounded-b-lg">
                <Button type="button" variant="outline" onClick={() => setSelectedProductForVariant(null)}>
                  Cancel
                </Button>
                <Button type="submit" variant="primary">
                  Add Variant
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}

      {/* Modal: New Tax */}
      {showTaxModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-md w-full shadow-lg">
            <CardHeader>
              <CardTitle className="text-lg">Add Tax Rule</CardTitle>
            </CardHeader>
            <form onSubmit={handleCreateTax}>
              <CardContent className="space-y-4">
                <div>
                  <Label>Tax Name</Label>
                  <Input required value={taxName} onChange={(e) => setTaxName(e.target.value)} />
                </div>
                <div>
                  <Label>Rate (Basis Points - 1800 = 18%)</Label>
                  <Input required type="number" value={taxRateBps} onChange={(e) => setTaxRateBps(e.target.value)} />
                </div>
              </CardContent>
              <div className="p-4 border-t flex justify-end gap-2 bg-slate-50 rounded-b-lg">
                <Button type="button" variant="outline" onClick={() => setShowTaxModal(false)}>
                  Cancel
                </Button>
                <Button type="submit" variant="primary">
                  Save Tax
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}

      {/* Modal: New Discount */}
      {showDiscountModal && (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <Card className="max-w-md w-full shadow-lg">
            <CardHeader>
              <CardTitle className="text-lg">Create Discount Rule</CardTitle>
            </CardHeader>
            <form onSubmit={handleCreateDiscount}>
              <CardContent className="space-y-4">
                <div>
                  <Label>Coupon Code</Label>
                  <Input required value={discountCode} onChange={(e) => setDiscountCode(e.target.value.toUpperCase())} />
                </div>
                <div>
                  <Label>Discount Type</Label>
                  <select
                    className="w-full border rounded px-3 py-2 text-sm bg-white"
                    value={discountType}
                    onChange={(e: any) => setDiscountType(e.target.value)}
                  >
                    <option value="PERCENTAGE">Percentage (%)</option>
                    <option value="FIXED_AMOUNT">Fixed Amount (cents)</option>
                  </select>
                </div>
                <div>
                  <Label>{discountType === 'PERCENTAGE' ? 'Percentage Basis Points (1000 = 10%)' : 'Amount in Cents'}</Label>
                  <Input required type="number" value={discountValue} onChange={(e) => setDiscountValue(e.target.value)} />
                </div>
              </CardContent>
              <div className="p-4 border-t flex justify-end gap-2 bg-slate-50 rounded-b-lg">
                <Button type="button" variant="outline" onClick={() => setShowDiscountModal(false)}>
                  Cancel
                </Button>
                <Button type="submit" variant="primary">
                  Save Discount
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  );
}
