import React from 'react';

import { useNavigate, useLocation } from 'react-router-dom';

import { AuthUserContextType, useAuthUserContext } from '@/providers/AuthUser';
import { UserType, RoleType } from '@/types';

type CustomLocation = {
  state: { from: { pathname: string } };
};

export const Login = () => {
  const navigate = useNavigate();
  const location: CustomLocation = useLocation() as CustomLocation;
  const fromPathName: string = location.state.from.pathname;
  const authUser: AuthUserContextType = useAuthUserContext();

  const signin = (role: RoleType) => {
    const user: UserType = {
      name: 'no-name',
      role: role,
    };
    authUser.signin(user, () => {
      navigate(fromPathName, { replace: true });
    });
  };

  return (
    <div>
      <h3>Login</h3>
      <button onClick={() => signin(RoleType.Admin)}>admin権限でログイン</button>
      <button onClick={() => signin(RoleType.Manager)}>manager権限でログイン</button>
      <button onClick={() => signin(RoleType.User)}>user権限でログイン</button>
    </div>
  );
};
