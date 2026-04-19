'use client';
import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import apiClient from '@/lib/api-client';
import type { Enrollment, PaginatedResponse } from '@/types';

export default function EnrollmentsPage() {
  const [enrollments, setEnrollments] = useState<Enrollment[]>([]);
  const [filter, setFilter] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const params = new URLSearchParams();
    if (filter) params.set('reviewStatus', filter);
    apiClient.get<PaginatedResponse<Enrollment>>(`/enrollments?${params}`)
      .then(({ data }) => setEnrollments(data.data || [])).catch(() => {}).finally(() => setLoading(false));
  }, [filter]);

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Inscriptions</h1>
      <div className="flex gap-2 mb-4">
        {['', 'pending', 'approved', 'rejected', 'needs_correction'].map((s) => (
          <button key={s} onClick={() => setFilter(s)}
            className={`px-3 py-1 rounded-full text-xs ${filter === s ? 'bg-emerald-800 text-white' : 'bg-gray-100 text-gray-600'}`}>
            {s || 'Toutes'}
          </button>
        ))}
      </div>
      <Card>
        <CardHeader><p className="text-sm text-gray-500">{enrollments.length} inscription(s)</p></CardHeader>
        <CardContent>
          {loading ? <p className="text-gray-400 py-8 text-center">Chargement...</p> : enrollments.length === 0 ?
            <p className="text-gray-400 py-8 text-center">Aucune inscription</p> : (
            <table className="w-full text-sm">
              <thead><tr className="border-b text-left text-gray-500">
                <th className="pb-3 font-medium">Date</th><th className="pb-3 font-medium">Agent</th>
                <th className="pb-3 font-medium">Lieu</th><th className="pb-3 font-medium">Sync</th>
                <th className="pb-3 font-medium">Revue</th>
              </tr></thead>
              <tbody>{enrollments.map((e) => (
                <tr key={e.id} className="border-b last:border-0">
                  <td className="py-3">{new Date(e.enrolledAt).toLocaleDateString('fr-FR')}</td>
                  <td className="py-3">{e.agentId.slice(0, 8)}</td>
                  <td className="py-3">{e.locationName || '—'}</td>
                  <td className="py-3"><Badge status={e.syncStatus} /></td>
                  <td className="py-3"><Badge status={e.reviewStatus} /></td>
                </tr>
              ))}</tbody>
            </table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
