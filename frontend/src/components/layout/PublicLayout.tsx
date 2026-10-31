import { Outlet, Link } from 'react-router-dom';
import { Layers } from 'lucide-react';

export function PublicLayout() {
  return (
    <div className="min-h-screen bg-slate-50 flex flex-col justify-between selection:bg-indigo-500 selection:text-white">
      <header className="border-b border-slate-200 bg-white/80 backdrop-blur-sm sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2 text-indigo-600 font-bold text-xl tracking-tight">
            <div className="h-8 w-8 rounded-lg bg-indigo-600 flex items-center justify-center text-white shadow-sm">
              <Layers className="h-5 w-5" />
            </div>
            <span>Veyra</span>
          </Link>
          <div className="flex items-center gap-4 text-sm font-medium">
            <Link to="/login" className="text-slate-600 hover:text-slate-900 transition-colors">
              Sign In
            </Link>
            <Link
              to="/register"
              className="bg-indigo-600 text-white px-3.5 py-1.5 rounded-md hover:bg-indigo-700 transition-colors shadow-sm"
            >
              Get Started
            </Link>
          </div>
        </div>
      </header>

      <main className="flex-1 flex items-center justify-center p-4 sm:p-6 lg:p-8">
        <Outlet />
      </main>

      <footer className="border-t border-slate-200 bg-white py-6 text-center text-xs text-slate-500">
        <div className="max-w-7xl mx-auto px-4">
          Brightloop Software &bull; Veyra Subscription Platform &bull; Phase 0 Foundation
        </div>
      </footer>
    </div>
  );
}
