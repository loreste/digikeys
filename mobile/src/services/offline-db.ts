import type { EnrollmentData } from '../types';

// Simple in-memory offline store (replace with SQLite in production)
let enrollments: EnrollmentData[] = [];

export const saveEnrollment = (enrollment: EnrollmentData) => {
  enrollments.push(enrollment);
};

export const getEnrollments = (): EnrollmentData[] => [...enrollments];

export const updateSyncStatus = (id: string, status: EnrollmentData['syncStatus']) => {
  const idx = enrollments.findIndex((e) => e.id === id);
  if (idx >= 0) enrollments[idx].syncStatus = status;
};

export const getPendingSync = (): EnrollmentData[] =>
  enrollments.filter((e) => e.syncStatus === 'pending' || e.syncStatus === 'failed');

export const clearSynced = () => {
  enrollments = enrollments.filter((e) => e.syncStatus !== 'synced');
};
