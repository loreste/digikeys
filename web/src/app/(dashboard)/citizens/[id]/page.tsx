'use client';
import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import apiClient from '@/lib/api-client';
import type { Citizen } from '@/types';

export default function CitizenDetailPage() {
  const params = useParams();
  const [citizen, setCitizen] = useState<Citizen | null>(null);

  useEffect(() => {
    apiClient.get<Citizen>(`/citizens/${params.id}`).then(({ data }) => setCitizen(data)).catch(() => {});
  }, [params.id]);

  if (!citizen) return <p className="text-gray-400">Chargement...</p>;

  const fields = [
    ['Prénom', citizen.firstName], ['Nom', citizen.lastName], ['Nom de jeune fille', citizen.maidenName],
    ['Date de naissance', new Date(citizen.dateOfBirth).toLocaleDateString('fr-FR')],
    ['Lieu de naissance', citizen.placeOfBirth], ['Genre', citizen.gender === 'M' ? 'Masculin' : 'Féminin'],
    ['CNIB', citizen.nationalId], ['Passeport', citizen.passportNumber],
    ['Téléphone', citizen.phone], ['Email', citizen.email],
    ['Pays de résidence', citizen.countryOfResidence], ['Ville', citizen.cityOfResidence],
    ['Province d\'origine', citizen.provinceOfOrigin], ['Commune d\'origine', citizen.communeOfOrigin],
  ];

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <h1 className="text-2xl font-bold text-gray-800">{citizen.firstName} {citizen.lastName}</h1>
        <Badge status={citizen.status} />
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader><h2 className="font-semibold">Informations personnelles</h2></CardHeader>
          <CardContent>
            <dl className="space-y-2">
              {fields.map(([label, value]) => value && (
                <div key={label} className="flex justify-between text-sm">
                  <dt className="text-gray-500">{label}</dt>
                  <dd className="font-medium text-gray-800">{value}</dd>
                </div>
              ))}
            </dl>
          </CardContent>
        </Card>
        <Card>
          <CardHeader><h2 className="font-semibold">Cartes et historique</h2></CardHeader>
          <CardContent className="text-gray-400 text-sm">Les cartes associées apparaîtront ici.</CardContent>
        </Card>
      </div>
    </div>
  );
}
