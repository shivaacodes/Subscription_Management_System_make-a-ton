import { ModulePlaceholder } from '@/components/ui/ModulePlaceholder';
import { Package } from 'lucide-react';

export function CatalogPage() {
  return (
    <ModulePlaceholder
      title="Products & Variants Catalog"
      description="Define base software products, pricing variants, and add-on products for Brightloop Software."
      icon={Package}
      owner="Member 2 (Catalog & Configuration)"
      branch="feat/catalog-plans"
      deliverables={[
        'Base products with sales price, cost price, and status flags',
        'Variant price adjustments (e.g., Odoo brand adds +₹560)',
        'Add-on product classifications (optional vs required)',
        'Catalog search, filtering, and admin pricing management',
      ]}
    />
  );
}
