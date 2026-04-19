import { create } from 'zustand';
import type { User } from '@/types';
import { getToken, getUser, setToken, setRefreshToken, setUser as saveUser, removeTokens } from '@/lib/auth';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (token: string, refreshToken: string, user: User) => void;
  logout: () => void;
  initialize: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null, token: null, isAuthenticated: false,
  login: (token, refreshToken, user) => {
    setToken(token); setRefreshToken(refreshToken); saveUser(user);
    set({ token, user, isAuthenticated: true });
  },
  logout: () => { removeTokens(); set({ token: null, user: null, isAuthenticated: false }); },
  initialize: () => { set({ token: getToken(), user: getUser(), isAuthenticated: !!getToken() }); },
}));
