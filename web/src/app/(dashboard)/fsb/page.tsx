'use client';
import { Card, CardContent, CardHeader } from '@/components/ui/card';

export default function FSBPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Fonds de Solidarité Burkinabè (FSB)</h1>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <Card><CardContent className="pt-6 text-center"><p className="text-sm text-gray-500">Total collecté</p><p className="text-3xl font-bold text-emerald-600">— FCFA</p><p className="text-xs text-gray-400">1 500 FCFA par carte délivrée</p></CardContent></Card>
        <Card><CardContent className="pt-6 text-center"><p className="text-sm text-gray-500">Cartes contributrices</p><p className="text-3xl font-bold text-blue-600">—</p></CardContent></Card>
        <Card><CardContent className="pt-6 text-center"><p className="text-sm text-gray-500">Objectif 3 ans</p><p className="text-3xl font-bold text-amber-600">9 Mds FCFA</p></CardContent></Card>
      </div>
      <Card>
        <CardHeader><h2 className="font-semibold">Contributions par période</h2></CardHeader>
        <CardContent className="py-8 text-center text-gray-400">Les rapports FSB apparaîtront ici.</CardContent>
      </Card>
    </div>
  );
}
