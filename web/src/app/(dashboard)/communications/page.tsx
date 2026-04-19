'use client';
import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader } from '@/components/ui/card';

export default function CommunicationsPage() {
  const [form, setForm] = useState({ subject: '', body: '', channel: 'email', targetCountry: '' });

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Communications</h1>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader><h2 className="font-semibold">Envoyer une communication</h2></CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">Canal</label>
              <select value={form.channel} onChange={(e) => setForm(f => ({ ...f, channel: e.target.value }))}
                className="block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm">
                <option value="email">Email</option><option value="sms">SMS</option><option value="push">Notification push</option>
              </select>
            </div>
            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">Pays cible (vide = tous)</label>
              <select value={form.targetCountry} onChange={(e) => setForm(f => ({ ...f, targetCountry: e.target.value }))}
                className="block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm">
                <option value="">Tous les pays</option><option value="FR">France</option><option value="US">États-Unis</option>
                <option value="IT">Italie</option><option value="DE">Allemagne</option><option value="CA">Canada</option>
              </select>
            </div>
            <Input id="subject" label="Sujet" value={form.subject} onChange={(e) => setForm(f => ({ ...f, subject: e.target.value }))} />
            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">Message</label>
              <textarea value={form.body} onChange={(e) => setForm(f => ({ ...f, body: e.target.value }))}
                rows={6} className="block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm" />
            </div>
            <Button>Envoyer</Button>
          </CardContent>
        </Card>
        <Card>
          <CardHeader><h2 className="font-semibold">Communications récentes</h2></CardHeader>
          <CardContent className="py-8 text-center text-gray-400">Aucune communication envoyée</CardContent>
        </Card>
      </div>
    </div>
  );
}
