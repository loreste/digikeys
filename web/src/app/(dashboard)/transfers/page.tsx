'use client';
import { Card, CardContent, CardHeader } from '@/components/ui/card';

export default function TransfersPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Transferts et Épargne</h1>
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <Card><CardContent className="pt-4 text-center"><p className="text-xs text-gray-500">Total épargne (mois)</p><p className="text-2xl font-bold text-emerald-600">— FCFA</p></CardContent></Card>
        <Card><CardContent className="pt-4 text-center"><p className="text-xs text-gray-500">Transferts réussis</p><p className="text-2xl font-bold text-blue-600">—</p></CardContent></Card>
        <Card><CardContent className="pt-4 text-center"><p className="text-xs text-gray-500">En attente</p><p className="text-2xl font-bold text-amber-600">—</p></CardContent></Card>
        <Card><CardContent className="pt-4 text-center"><p className="text-xs text-gray-500">Comptes ouverts</p><p className="text-2xl font-bold text-purple-600">—</p></CardContent></Card>
      </div>
      <Card>
        <CardHeader><h2 className="font-semibold">Historique des transferts</h2></CardHeader>
        <CardContent className="py-8 text-center text-gray-400">Les transferts apparaîtront ici une fois le système opérationnel.</CardContent>
      </Card>
    </div>
  );
}
