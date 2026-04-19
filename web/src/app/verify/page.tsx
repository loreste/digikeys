'use client';
import { useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import apiClient from '@/lib/api-client';
import type { VerificationResult } from '@/types';

export default function VerifyPage() {
  const [cardNumber, setCardNumber] = useState('');
  const [result, setResult] = useState<VerificationResult | null>(null);
  const [loading, setLoading] = useState(false);

  const handleVerify = async () => {
    if (!cardNumber.trim()) return;
    setLoading(true); setResult(null);
    try {
      const { data } = await apiClient.get<VerificationResult>(`/verify/card/${cardNumber.trim()}`);
      setResult(data);
    } catch { setResult({ valid: false, status: 'NOT_FOUND', message: 'Carte non trouvée' }); }
    finally { setLoading(false); }
  };

  return (
    <main className="min-h-screen bg-gray-50">
      <nav className="bg-emerald-900 px-8 py-4">
        <div className="max-w-3xl mx-auto flex justify-between items-center">
          <Link href="/" className="text-xl font-bold text-amber-400">DIGIKEYS</Link>
          <Link href="/login" className="text-emerald-200 hover:text-white text-sm">Se Connecter</Link>
        </div>
      </nav>
      <div className="max-w-xl mx-auto px-4 py-12">
        <h1 className="text-3xl font-bold text-gray-800 text-center mb-2">Vérifier une Carte Consulaire</h1>
        <p className="text-gray-500 text-center mb-8">Entrez le numéro de la carte consulaire pour vérifier sa validité</p>
        <Card>
          <CardContent className="pt-6 space-y-4">
            <Input id="card" label="Numéro de carte" placeholder="CC-FR-2026-000001" value={cardNumber} onChange={(e) => setCardNumber(e.target.value)} />
            <Button onClick={handleVerify} className="w-full" disabled={loading}>{loading ? 'Vérification...' : 'Vérifier'}</Button>
          </CardContent>
        </Card>
        {result && (
          <Card className="mt-6">
            <CardContent className="pt-6 text-center">
              <div className="text-6xl mb-4">{result.valid ? '✅' : '❌'}</div>
              <h2 className={`text-2xl font-bold mb-2 ${result.valid ? 'text-emerald-600' : 'text-red-600'}`}>
                {result.valid ? 'Carte Valide' : result.status === 'expired' ? 'Carte Expirée' : result.status === 'revoked' ? 'Carte Révoquée' : 'Carte Non Trouvée'}
              </h2>
              <p className="text-gray-500 mb-4">{result.message}</p>
              {result.valid && (
                <div className="text-left bg-gray-50 rounded-lg p-4 space-y-2 text-sm">
                  {result.holderName && <p><span className="text-gray-500">Titulaire:</span> {result.holderName}</p>}
                  {result.embassy && <p><span className="text-gray-500">Ambassade:</span> {result.embassy}</p>}
                  {result.expiresAt && <p><span className="text-gray-500">Expire le:</span> {new Date(result.expiresAt).toLocaleDateString('fr-FR')}</p>}
                </div>
              )}
            </CardContent>
          </Card>
        )}
      </div>
    </main>
  );
}
