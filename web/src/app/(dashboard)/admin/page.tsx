'use client';
import { Card, CardContent } from '@/components/ui/card';

const stats = [
  { label: 'Citoyens enregistrés', value: '—', color: 'text-emerald-600' },
  { label: 'Cartes actives', value: '—', color: 'text-blue-600' },
  { label: 'Inscriptions en attente', value: '—', color: 'text-amber-600' },
  { label: 'Contribution FSB', value: '—', color: 'text-purple-600' },
  { label: 'Transferts (mois)', value: '—', color: 'text-cyan-600' },
  { label: 'Ambassades', value: '—', color: 'text-gray-600' },
];

export default function AdminDashboard() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Tableau de bord - Administration Centrale</h1>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4 mb-8">
        {stats.map((s) => (
          <Card key={s.label}>
            <CardContent className="pt-4 pb-4">
              <p className="text-xs text-gray-500">{s.label}</p>
              <p className={`text-2xl font-bold mt-1 ${s.color}`}>{s.value}</p>
            </CardContent>
          </Card>
        ))}
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card><CardContent className="py-12 text-center text-gray-400">Répartition par pays - données en attente</CardContent></Card>
        <Card><CardContent className="py-12 text-center text-gray-400">Évolution des inscriptions - données en attente</CardContent></Card>
      </div>
    </div>
  );
}
