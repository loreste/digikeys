import axios from 'axios';

const API_BASE = 'http://localhost:8081/api/v1';

const api = axios.create({ baseURL: API_BASE, headers: { 'Content-Type': 'application/json' } });

let authToken: string | null = null;
export const setAuthToken = (token: string) => { authToken = token; };

api.interceptors.request.use((config) => {
  if (authToken) config.headers.Authorization = `Bearer ${authToken}`;
  return config;
});

export const login = async (email: string, password: string) => {
  const { data } = await api.post('/auth/login', { email, password });
  return data;
};

export const syncEnrollments = async (enrollments: object[]) => {
  const { data } = await api.post('/enrollments/sync', { enrollments });
  return data;
};

export const uploadBiometrics = async (enrollmentId: string, biometrics: object) => {
  const { data } = await api.post(`/enrollments/${enrollmentId}/biometrics`, biometrics);
  return data;
};

export const uploadPhoto = async (enrollmentId: string, photoBase64: string) => {
  const { data } = await api.post(`/enrollments/${enrollmentId}/photo`, { photo: photoBase64 });
  return data;
};

export default api;
