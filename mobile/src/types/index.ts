export interface EnrollmentData {
  id: string;
  firstName: string;
  lastName: string;
  maidenName?: string;
  dateOfBirth: string;
  placeOfBirth: string;
  gender: 'M' | 'F';
  nationalId?: string;
  passportNumber?: string;
  phone?: string;
  email?: string;
  countryOfResidence: string;
  cityOfResidence?: string;
  addressAbroad?: string;
  provinceOfOrigin?: string;
  communeOfOrigin?: string;
  photoUri?: string;
  fingerprintsCapured: boolean;
  syncStatus: 'pending' | 'syncing' | 'synced' | 'failed';
  createdAt: string;
}

export interface BiometricData {
  rightThumb?: string; // base64
  rightIndex?: string;
  leftThumb?: string;
  leftIndex?: string;
  qualityScores: { rightThumb: number; rightIndex: number; leftThumb: number; leftIndex: number };
}
