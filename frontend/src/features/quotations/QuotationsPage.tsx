import { ModulePlaceholder } from '@/components/ui/ModulePlaceholder';
import { FileText } from 'lucide-react';

export function QuotationsPage() {
  return (
    <ModulePlaceholder
      title="Quotations & Templates"
      description="Create reusable quotation templates and client-specific proposals with automatic conversion to subscriptions."
      icon={FileText}
      owner="Member 4 (Quotations & Experience)"
      branch="feat/quotations-portal"
      deliverables={[
        'Reusable quotation templates with preconfigured lines and plans',
        'Customer quotation creation with validity period and tax calculation',
        'One-click conversion from accepted quotation to active subscription',
        'Quotation expiration tracking and customer approval workflow',
      ]}
    />
  );
}
