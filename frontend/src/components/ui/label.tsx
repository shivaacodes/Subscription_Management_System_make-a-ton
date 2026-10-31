import React from 'react';
import { cn } from '@/lib/utils';

export function Label({ className, children, ...props }: React.LabelHTMLAttributes<HTMLLabelElement>) {
  return (
    <label
      className={cn('block text-xs font-semibold uppercase tracking-wider text-slate-700 mb-1.5', className)}
      {...props}
    >
      {children}
    </label>
  );
}
