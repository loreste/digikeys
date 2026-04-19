import Link from 'next/link';

export default function Home() {
  return (
    <main className="min-h-screen bg-gradient-to-b from-emerald-900 to-emerald-800">
      <nav className="flex items-center justify-between px-8 py-6">
        <h1 className="text-2xl font-bold text-amber-400">DIGIKEYS</h1>
        <div className="flex gap-4">
          <Link href="/verify" className="text-emerald-200 hover:text-white transition-colors">Vérifier une Carte</Link>
          <Link href="/login" className="bg-white text-emerald-800 px-4 py-2 rounded-lg font-medium hover:bg-emerald-50">Se Connecter</Link>
        </div>
      </nav>

      <section className="max-w-5xl mx-auto px-8 py-20 text-center">
        <div className="text-6xl mb-6">🇧🇫</div>
        <h2 className="text-5xl font-bold text-white mb-6">Carte Consulaire<br />Biométrique</h2>
        <p className="text-xl text-emerald-200 mb-10 max-w-2xl mx-auto">
          Plus qu&apos;une carte consulaire, un véritable outil au service du développement
          du Burkina Faso par sa diaspora. Une carte pour tous, tous pour le Burkina.
        </p>
        <div className="flex justify-center gap-4">
          <Link href="/verify" className="bg-amber-600 text-white px-8 py-3 rounded-lg text-lg font-medium hover:bg-amber-700">Vérifier une Carte</Link>
          <Link href="/login" className="bg-white text-emerald-800 px-8 py-3 rounded-lg text-lg font-medium hover:bg-emerald-50">Espace Consulaire</Link>
        </div>
      </section>

      <section className="max-w-5xl mx-auto px-8 py-16">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="bg-white/10 backdrop-blur rounded-xl p-8 text-center">
            <div className="text-4xl mb-4">🛡️</div>
            <h3 className="text-xl font-bold text-white mb-3">Sécurité Biométrique</h3>
            <p className="text-emerald-200">Carte PVC avec données biométriques, MRZ aux normes ICAO, et zone de lecture automatique.</p>
          </div>
          <div className="bg-white/10 backdrop-blur rounded-xl p-8 text-center">
            <div className="text-4xl mb-4">💰</div>
            <h3 className="text-xl font-bold text-white mb-3">Inclusion Financière</h3>
            <p className="text-emerald-200">Compte bancaire burkinabè ouvert pour chaque carte. Transferts et épargne dirigée vers le Burkina.</p>
          </div>
          <div className="bg-white/10 backdrop-blur rounded-xl p-8 text-center">
            <div className="text-4xl mb-4">🌍</div>
            <h3 className="text-xl font-bold text-white mb-3">Diaspora Connectée</h3>
            <p className="text-emerald-200">Communication directe avec les autorités consulaires. Équipes mobiles d&apos;enrôlement partout dans le monde.</p>
          </div>
        </div>
      </section>

      <footer className="text-center py-8 text-emerald-300 text-sm">
        DIGIKEYS - Carte Consulaire du Burkina Faso - Fonds de Solidarité Burkinabè (FSB)
      </footer>
    </main>
  );
}
