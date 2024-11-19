import React, { createContext, useContext, useEffect, useState } from 'react';

type UserType = {
  email: string;
  name: string;
};

const AuthProvider = createContext<{
  user: UserType | null | undefined;
  setUser: (u: UserType) => void;
  logout: () => void;
  isLoggedIn: undefined | boolean;
}>({
  user: null,
  setUser: (_u: UserType) => {},
  logout: () => {},
  isLoggedIn: undefined,
});

const USER_KEY = '__USER_AUTHED__';

export const AuthCtx = ({ children }: { children: React.ReactNode }) => {
  const [user, _setUser] = useState<UserType | null>();
  const [isLoggedIn, setIsLoggedIn] = useState<boolean | undefined>(undefined);

  const setUser = (user: UserType) => {
    _setUser(user);

    localStorage.setItem(USER_KEY, JSON.stringify(user));
    setIsLoggedIn(true);
  };

  const logout = () => {
    _setUser(null);

    localStorage.removeItem(USER_KEY);
    setIsLoggedIn(false);
  };

  // on mount, if user present in local storage set the user detail
  useEffect(() => {
    const userStored = localStorage.getItem(USER_KEY);
    if (!userStored) {
      setIsLoggedIn(false);
      return;
    }
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
