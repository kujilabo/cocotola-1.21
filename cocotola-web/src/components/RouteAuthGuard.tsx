import React from 'react';

import { Navigate, useLocation } from 'react-router-dom';

import { useAuthUserContext } from '@/providers';
import { RoleType } from '@/types';

type Props = {
  component: React.ReactNode;
  redirect: string;
  allowroles?: RoleType[];
};

export const RouteAuthGuard: React.VFC<Props> = (props) => {
  const authUser = useAuthUserContext().user;

  let allowRoute = false;
  if (authUser) {
    allowRoute = props.allowroles ? props.allowroles.includes(authUser.role) : true;
  }

  if (!allowRoute) {
    return <Navigate to={props.redirect} state={{ from: useLocation() }} replace={false} />;
  }

  return <>{props.component}</>;
};
