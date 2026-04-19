'use client';
import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import apiClient from '@/lib/api-client';
import type { Card as CardType, PaginatedResponse } from '@/types';

export default function CardsPage() {
  const [cards, setCards] = useState<CardType[]>([]);
  const [statusFilter, setStatusFilter] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const params = new URLSearchParams();
    if (statusFilter) params.set('status', statusFilter);
    apiClient.get<PaginatedResponse<CardType>>(`/cards?${params}`)
      .then(({ data }) => setCards(data.data || [])).catch(() => {}).finally(() => setLoading(false));
  }, [statusFilter]);

  const statuses = ['', 'pending', 'approved', 'printing', 'printed', 'delivered', 'active', 'suspended', 'revoked'];

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Cartes Consulaires</h1>
      <div className="flex gap-2 mb-4 flex-wrap">
        {statuses.map((s) => (
          <button key={s} onClick={() => setStatusFilter(s)}
            className={`px-3 py-1 rounded-full text-xs ${statusFilter === s ? 'bg-emerald-800 text-white' : 'bg-gray-100 text-gray-600'}`}>
            {s || 'Toutes'}
          </button>
        ))}
      </div>
      <Card>
        <CardHeader><p className="text-sm text-gray-500">{cards.length} carte(s)</p></CardHeader>
        <CardContent>
          {loading ? <p className="text-gray-400 py-8 text-center">Chargement...</p> : cards.length === 0 ?
            <p className="text-gray-400 py-8 text-center">Aucune carte</p> : (
            <table className="w-full text-sm">
              <thead><tr className="border-b text-left text-gray-500">
                <th className="pb-3 font-medium">N° Carte</th><th className="pb-3 font-medium">Ambassade</th>
                <th className="pb-3 font-medium">Expire le</th><th className="pb-3 font-medium">Statut</th>
                <th className="pb-3 font-medium">Actions</th>
              </tr></thead>
              <tbody>{cards.map((c) => (
                <tr key={c.id} className="border-b last:border-0">
                  <td className="py-3 font-mono text-xs">{c.cardNumber}</td>
                  <td className="py-3">{c.embassyId}</td>
                  <td className="py-3">{c.expiresAt ? new Date(c.expiresAt).toLocaleDateString('fr-FR') : '—'}</td>
                  <td className="py-3"><Badge status={c.status} /></td>
                  <td className="py-3 flex gap-2">
                    {c.status === 'pending' && <Button size="sm" variant="secondary" onClick={() => {}}>Approuver</Button>}
                    {c.status === 'approved' && <Button size="sm" onClick={() => {}}>Imprimer</Button>}
                    {c.status === 'active' && <Button size="sm" variant="danger" onClick={() => {}}>Suspendre</Button>}
                  </td>
                </tr>
              ))}</tbody>
            </table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
