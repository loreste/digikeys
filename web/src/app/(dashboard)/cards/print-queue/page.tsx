'use client';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export default function PrintQueuePage() {
  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">File d&apos;impression</h1>
        <Button>Créer un lot d&apos;impression</Button>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <Card><CardContent className="pt-4 text-center"><p className="text-xs text-gray-500">En attente</p><p className="text-3xl font-bold text-amber-600">—</p></CardContent></Card>
        <Card><CardContent className="pt-4 text-center"><p className="text-xs text-gray-500">En impression</p><p className="text-3xl font-bold text-purple-600">—</p></CardContent></Card>
        <Card><CardContent className="pt-4 text-center"><p className="text-xs text-gray-500">Imprimées</p><p className="text-3xl font-bold text-emerald-600">—</p></CardContent></Card>
      </div>
      <Card>
        <CardHeader><h2 className="font-semibold">Lots d&apos;impression</h2></CardHeader>
        <CardContent className="py-8 text-center text-gray-400">Aucun lot d&apos;impression en cours. Sélectionnez des cartes approuvées pour créer un lot.</CardContent>
      </Card>
    </div>
  );
}
