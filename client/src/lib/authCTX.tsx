import React, { createContext, useContext, useEffect, useState } from 'react';
import { axiosInstance } from './axiosInstance';

type UserType = {
  email: string;
  name: string;
};

const AuthProvider = createContext<{
  user: UserType | null | undefined;
  setUser: (u: UserType, jwt: string) => void;
  logout: () => void;
  isLoggedIn: undefined | boolean;
}>({
  user: null,
  setUser: (_u: UserType, _jwt: string) => {},
  logout: () => {},
  isLoggedIn: undefined,
});

const USER_KEY = '__USER_AUTHED__';

export const AUTH_TOKEN = '__AUTH_TOKEN__';

export const AuthCtx = ({ children }: { children: React.ReactNode }) => {
  const [user, _setUser] = useState<UserType | null>(null);
  const [isLoggedIn, setIsLoggedIn] = useState<boolean | undefined>(undefined);

  const setUser = (user: UserType, jwt: string) => {
    _setUser(user);
    localStorage.setItem(AUTH_TOKEN, jwt);
    localStorage.setItem(USER_KEY, JSON.stringify(user));
    setIsLoggedIn(true);
  };

  const logout = async () => {
    try {
      await axiosInstance.post('/auth/logout');
    } catch (error) {
    } finally {
      _setUser(null);

      localStorage.removeItem(USER_KEY);
      localStorage.removeItem(AUTH_TOKEN);

      setIsLoggedIn(false);
    }
  };

  useEffect(() => {
    const userStored = localStorage.getItem(USER_KEY);
    if (!userStored) {
      setIsLoggedIn(false);
      return;
    }
    // on mount, if user persisted in local storage set the user detail
   setIsLoggedIn(!false);

    _setUser(JSON.parse(userStored));
  }, []);

  useEffect(() => {
    if (!user) return;
    setIsLoggedIn(true);
  }, [user]);

  return (
    <AuthProvider.Provider value={{ user, setUser, logout, isLoggedIn }}>
      {children}
    </AuthProvider.Provider>
  );
};

export const useAuth = () => useContext(AuthProvider);
