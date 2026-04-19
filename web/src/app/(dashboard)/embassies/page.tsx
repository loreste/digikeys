'use client';
import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import apiClient from '@/lib/api-client';
import type { Embassy } from '@/types';

export default function EmbassiesPage() {
  const [embassies, setEmbassies] = useState<Embassy[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiClient.get<Embassy[]>('/embassies').then(({ data }) => setEmbassies(data || [])).catch(() => {}).finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Ambassades et Consulats</h1>
        <Button>Ajouter</Button>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {loading ? <p className="text-gray-400">Chargement...</p> : embassies.length === 0 ?
          <p className="text-gray-400">Aucune ambassade enregistrée</p> :
          embassies.map((e) => (
            <Card key={e.id}>
              <CardHeader className="flex flex-row items-center justify-between">
                <div><h3 className="font-semibold">{e.name}</h3><p className="text-xs text-gray-500">{e.city} - {e.countryCode}</p></div>
                <Badge status={e.status} />
              </CardHeader>
              <CardContent>
                <div className="text-sm space-y-1">
                  <p><span className="text-gray-500">Consul:</span> {e.consulName || '—'}</p>
                  <p><span className="text-gray-500">Préfixe carte:</span> <span className="font-mono">{e.cardPrefix}</span></p>
                  {e.phone && <p><span className="text-gray-500">Tél:</span> {e.phone}</p>}
                  {e.email && <p><span className="text-gray-500">Email:</span> {e.email}</p>}
                </div>
              </CardContent>
            </Card>
          ))}
      </div>
    </div>
  );
}
