'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { useAuthStore } from '@/stores/auth-store';
import apiClient from '@/lib/api-client';
import type { LoginResponse } from '@/types';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const login = useAuthStore((s) => s.login);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(''); setLoading(true);
    try {
      const { data } = await apiClient.post<LoginResponse>('/auth/login', { email, password });
      login(data.accessToken, data.refreshToken, data.user);
      router.push('/admin');
    } catch { setError('Identifiants incorrects'); }
    finally { setLoading(false); }
  };

  return (
    <Card>
      <CardHeader><h2 className="text-xl font-semibold text-gray-800">Connexion</h2></CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && <div className="bg-red-50 text-red-700 p-3 rounded-lg text-sm">{error}</div>}
          <Input id="email" label="Email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
          <Input id="password" label="Mot de passe" type="password" value={password} onChange={(e) => setPassword(e.target.value)} required />
          <Button type="submit" className="w-full" disabled={loading}>{loading ? 'Connexion...' : 'Se Connecter'}</Button>
        </form>
      </CardContent>
    </Card>
  );
}
