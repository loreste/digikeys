'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';
import { useAuthStore } from '@/stores/auth-store';

const menuItems: Record<string, { label: string; href: string; icon: string }[]> = {
  super_admin: [
    { label: 'Tableau de bord', href: '/admin', icon: '📊' },
    { label: 'Ambassades', href: '/embassies', icon: '🏛️' },
    { label: 'Citoyens', href: '/citizens', icon: '👥' },
    { label: 'Cartes', href: '/cards', icon: '💳' },
    { label: 'Inscriptions', href: '/enrollments', icon: '📝' },
    { label: 'Transferts', href: '/transfers', icon: '💰' },
    { label: 'FSB', href: '/fsb', icon: '🤝' },
    { label: 'Communications', href: '/communications', icon: '📢' },
  ],
  embassy_admin: [
    { label: 'Tableau de bord', href: '/admin', icon: '📊' },
    { label: 'Citoyens', href: '/citizens', icon: '👥' },
    { label: 'Cartes', href: '/cards', icon: '💳' },
    { label: 'Inscriptions', href: '/enrollments', icon: '📝' },
    { label: 'Communications', href: '/communications', icon: '📢' },
  ],
  print_operator: [
    { label: "File d'impression", href: '/cards/print-queue', icon: '🖨️' },
    { label: 'Cartes', href: '/cards', icon: '💳' },
  ],
  enrollment_agent: [
    { label: 'Mes inscriptions', href: '/enrollments', icon: '📝' },
  ],
  verifier: [
    { label: 'Vérification', href: '/verify', icon: '🔍' },
  ],
};

export function Sidebar() {
  const pathname = usePathname();
  const user = useAuthStore((s) => s.user);
  const role = user?.role || 'readonly';
  const items = menuItems[role] || menuItems.verifier;

  return (
    <aside className="w-64 bg-emerald-950 text-white min-h-screen flex flex-col">
      <div className="p-6 border-b border-emerald-800">
        <h1 className="text-xl font-bold text-amber-400">DIGIKEYS</h1>
        <p className="text-xs text-emerald-300 mt-1">Carte Consulaire Biométrique</p>
      </div>
      <nav className="flex-1 py-4">
        {items.map((item) => (
          <Link key={item.href} href={item.href}
            className={cn('flex items-center gap-3 px-6 py-3 text-sm transition-colors',
              pathname === item.href ? 'bg-emerald-800 text-white border-r-2 border-amber-400' : 'text-emerald-200 hover:bg-emerald-900')}>
            <span>{item.icon}</span>{item.label}
          </Link>
        ))}
      </nav>
      <div className="p-4 border-t border-emerald-800">
        <p className="text-xs text-emerald-500">v1.0.0 - Burkina Faso 🇧🇫</p>
      </div>
    </aside>
  );
}
