import { create } from 'zustand';
import type { EnrollmentData } from '../types';

interface EnrollmentState {
  enrollments: EnrollmentData[];
  currentForm: Partial<EnrollmentData>;
  isSyncing: boolean;
  addEnrollment: (e: EnrollmentData) => void;
  updateForm: (fields: Partial<EnrollmentData>) => void;
  resetForm: () => void;
  setSyncing: (s: boolean) => void;
  updateSyncStatus: (id: string, status: EnrollmentData['syncStatus']) => void;
}

export const useEnrollmentStore = create<EnrollmentState>((set) => ({
  enrollments: [],
  currentForm: {},
  isSyncing: false,
  addEnrollment: (e) => set((s) => ({ enrollments: [...s.enrollments, e] })),
  updateForm: (fields) => set((s) => ({ currentForm: { ...s.currentForm, ...fields } })),
  resetForm: () => set({ currentForm: {} }),
  setSyncing: (isSyncing) => set({ isSyncing }),
  updateSyncStatus: (id, status) => set((s) => ({
    enrollments: s.enrollments.map((e) => e.id === id ? { ...e, syncStatus: status } : e),
  })),
}));
