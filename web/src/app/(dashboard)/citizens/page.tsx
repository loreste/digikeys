'use client';
import { useState, useEffect } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import apiClient from '@/lib/api-client';
import type { Citizen, PaginatedResponse } from '@/types';

export default function CitizensPage() {
  const [citizens, setCitizens] = useState<Citizen[]>([]);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const params = new URLSearchParams();
    if (search) params.set('q', search);
    apiClient.get<PaginatedResponse<Citizen>>(`/citizens?${params}`)
      .then(({ data }) => setCitizens(data.data || []))
      .catch(() => {}).finally(() => setLoading(false));
  }, [search]);

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Citoyens</h1>
      </div>
      <div className="mb-4 max-w-md">
        <Input id="search" placeholder="Rechercher par nom, CNIB, passeport..." value={search} onChange={(e) => setSearch(e.target.value)} />
      </div>
      <Card>
        <CardHeader><p className="text-sm text-gray-500">{citizens.length} citoyen(s)</p></CardHeader>
        <CardContent>
          {loading ? <p className="text-gray-400 py-8 text-center">Chargement...</p> : citizens.length === 0 ?
            <p className="text-gray-400 py-8 text-center">Aucun citoyen trouvé</p> : (
            <table className="w-full text-sm">
              <thead><tr className="border-b text-left text-gray-500">
                <th className="pb-3 font-medium">Nom</th><th className="pb-3 font-medium">Date de naissance</th>
                <th className="pb-3 font-medium">Pays de résidence</th><th className="pb-3 font-medium">CNIB</th>
                <th className="pb-3 font-medium">Statut</th><th className="pb-3 font-medium">Actions</th>
              </tr></thead>
              <tbody>{citizens.map((c) => (
                <tr key={c.id} className="border-b last:border-0">
                  <td className="py-3 font-medium">{c.firstName} {c.lastName}</td>
                  <td className="py-3">{new Date(c.dateOfBirth).toLocaleDateString('fr-FR')}</td>
                  <td className="py-3">{c.countryOfResidence}</td>
                  <td className="py-3 font-mono text-xs">{c.nationalId || '—'}</td>
                  <td className="py-3"><Badge status={c.status} /></td>
                  <td className="py-3"><Link href={`/citizens/${c.id}`} className="text-emerald-600 hover:underline text-xs">Détails</Link></td>
                </tr>
              ))}</tbody>
            </table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
