const T = 'digikeys_token', R = 'digikeys_refresh', U = 'digikeys_user';
export const getToken = () => typeof window !== 'undefined' ? localStorage.getItem(T) : null;
export const setToken = (t: string) => localStorage.setItem(T, t);
export const setRefreshToken = (t: string) => localStorage.setItem(R, t);
export const removeTokens = () => { localStorage.removeItem(T); localStorage.removeItem(R); localStorage.removeItem(U); };
export const isAuthenticated = () => !!getToken();
export const setUser = (u: object) => localStorage.setItem(U, JSON.stringify(u));
export const getUser = () => { if (typeof window === 'undefined') return null; const r = localStorage.getItem(U); return r ? JSON.parse(r) : null; };
