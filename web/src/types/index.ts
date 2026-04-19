export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  phone?: string;
  role: 'super_admin' | 'embassy_admin' | 'enrollment_agent' | 'print_operator' | 'bank_agent' | 'verifier' | 'readonly';
  embassyId?: string;
  status: string;
  createdAt: string;
}

export interface Citizen {
  id: string;
  firstName: string;
  lastName: string;
  maidenName?: string;
  dateOfBirth: string;
  placeOfBirth: string;
  gender: 'M' | 'F';
  nationality: string;
  nationalId?: string;
  uniqueIdentifier?: string;
  passportNumber?: string;
  phone?: string;
  email?: string;
  countryOfResidence: string;
  cityOfResidence?: string;
  addressAbroad?: string;
  provinceOfOrigin?: string;
  communeOfOrigin?: string;
  embassyId: string;
  photoKey?: string;
  status: 'registered' | 'verified' | 'active' | 'suspended' | 'deceased';
  createdAt: string;
}

export interface Card {
  id: string;
  citizenId: string;
  cardNumber: string;
  mrzLine1: string;
  mrzLine2: string;
  mrzLine3?: string;
  embassyId: string;
  issuedBy?: string;
  issuedAt?: string;
  expiresAt?: string;
  status: 'pending' | 'approved' | 'printing' | 'printed' | 'delivered' | 'active' | 'suspended' | 'revoked' | 'expired';
  printBatchId?: string;
  printedAt?: string;
  deliveredAt?: string;
  createdAt: string;
}

export interface Embassy {
  id: string;
  countryCode: string;
  name: string;
  city?: string;
  address?: string;
  phone?: string;
  email?: string;
  consulName?: string;
  cardPrefix: string;
  status: string;
}

export interface Enrollment {
  id: string;
  citizenId?: string;
  embassyId: string;
  agentId: string;
  teamId?: string;
  locationName?: string;
  syncStatus: 'pending' | 'syncing' | 'synced' | 'failed' | 'processed';
  reviewStatus: 'pending' | 'approved' | 'rejected' | 'needs_correction';
  reviewNotes?: string;
  enrolledAt: string;
}

export interface Transfer {
  id: string;
  citizenId: string;
  amount: number;
  currency: string;
  type: 'savings' | 'fsb_contribution' | 'remittance' | 'withdrawal';
  sourceProvider?: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  createdAt: string;
}

export interface BankAccount {
  id: string;
  citizenId: string;
  bankName: string;
  accountNumber?: string;
  status: string;
}

export interface FSBContribution {
  id: string;
  citizenId: string;
  cardId: string;
  amount: number;
  status: string;
  createdAt: string;
}

export interface VerificationResult {
  valid: boolean;
  status: string;
  cardNumber?: string;
  holderName?: string;
  embassy?: string;
  expiresAt?: string;
  message: string;
}

export interface Pagination {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: Pagination;
}

export interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  user: User;
}
