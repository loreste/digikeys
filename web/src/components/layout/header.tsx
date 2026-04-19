'use client';
import { useAuthStore } from '@/stores/auth-store';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';

export function Header() {
  const { user, logout } = useAuthStore();
  const router = useRouter();
  return (
    <header className="h-16 border-b border-gray-200 bg-white flex items-center justify-between px-6">
      <h2 className="text-lg font-semibold text-gray-800">
        {user?.role === 'super_admin' && 'Administration Centrale'}
        {user?.role === 'embassy_admin' && 'Gestion Ambassade'}
        {user?.role === 'print_operator' && 'Impression des Cartes'}
        {user?.role === 'enrollment_agent' && 'Agent d\'Enrôlement'}
        {user?.role === 'verifier' && 'Vérification'}
      </h2>
      <div className="flex items-center gap-4">
        <span className="text-sm text-gray-600">{user?.firstName} {user?.lastName}</span>
        <Button variant="ghost" size="sm" onClick={() => { logout(); router.push('/login'); }}>Déconnexion</Button>
      </div>
    </header>
  );
}
