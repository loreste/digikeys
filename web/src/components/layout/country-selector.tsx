'use client';
import { useCountryStore } from '@/stores/country-store';

// Country is fixed per deployment (BF or CD), not user-selectable.
// This component displays which country this deployment serves.
export function CountrySelector() {
  const { country } = useCountryStore();
  return (
    <div className="bg-emerald-900 text-emerald-200 text-sm rounded-lg px-3 py-1.5 border border-emerald-700 text-center">
      {country.flag} {country.name}
    </div>
  );
}
