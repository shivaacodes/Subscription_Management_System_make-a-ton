import React from 'react';
import { cn } from '@/lib/utils';
import { AlertCircle, CheckCircle2, Info, AlertTriangle } from 'lucide-react';

export interface AlertProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'info' | 'success' | 'warning' | 'danger';
  title?: string;
}

export function Alert({ variant = 'info', title, className, children, ...props }: AlertProps) {
  const styles = {
    info: {
      container: 'bg-indigo-50/70 border-indigo-200 text-indigo-900',
      icon: <Info className="h-5 w-5 text-indigo-600 flex-shrink-0" />,
    },
    success: {
      container: 'bg-emerald-50/70 border-emerald-200 text-emerald-900',
      icon: <CheckCircle2 className="h-5 w-5 text-emerald-600 flex-shrink-0" />,
    },
    warning: {
      container: 'bg-amber-50/70 border-amber-200 text-amber-900',
      icon: <AlertTriangle className="h-5 w-5 text-amber-600 flex-shrink-0" />,
    },
    danger: {
      container: 'bg-red-50/70 border-red-200 text-red-900',
      icon: <AlertCircle className="h-5 w-5 text-red-600 flex-shrink-0" />,
    },
  };

  const current = styles[variant];

  return (
    <div
      role="alert"
      className={cn('flex gap-3 p-4 rounded-md border text-sm', current.container, className)}
      {...props}
    >
      {current.icon}
      <div className="flex-1">
        {title && <h5 className="font-semibold mb-0.5">{title}</h5>}
        <div className="text-slate-700">{children}</div>
      </div>
    </div>
  );
}
