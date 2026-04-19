import { create } from 'zustand';

export interface Country {
  code: string;
  name: string;
  currency: string;
  motto: string;
  flag: string;
  fundName: string;
  fundAbbrev: string;
  contributionAmount: number;
  nationalityCode: string;
  providers: string[];
  banks: string[];
}

const COUNTRIES: Record<string, Country> = {
  BF: {
    code: 'BF', name: 'Burkina Faso', currency: 'XOF',
    motto: 'Unité - Progrès - Justice', flag: '🇧🇫',
    fundName: 'Fonds de Solidarité Burkinabè', fundAbbrev: 'FSB',
    contributionAmount: 1500, nationalityCode: 'BFA',
    providers: ['ORANGE_MONEY', 'MOOV_MONEY'],
    banks: ['Coris Bank International', 'Bank of Africa', 'Ecobank'],
  },
  CD: {
    code: 'CD', name: 'République Démocratique du Congo', currency: 'CDF',
    motto: 'Justice - Paix - Travail', flag: '🇨🇩',
    fundName: 'Fonds de Solidarité Congolais', fundAbbrev: 'FSC',
    contributionAmount: 5000, nationalityCode: 'COD',
    providers: ['ORANGE_MONEY_RDC', 'AIRTEL_MONEY', 'AFRICELL_MONEY', 'VODACOM'],
    banks: ['Rawbank', 'Equity BCDC', 'TMB', 'FBN Bank RDC'],
  },
};

// Country is determined by deployment, not user choice
const deployedCountryCode = process.env.NEXT_PUBLIC_COUNTRY || 'BF';

interface CountryState {
  country: Country;
}

export const useCountryStore = create<CountryState>(() => ({
  country: COUNTRIES[deployedCountryCode] || COUNTRIES.BF,
}));
