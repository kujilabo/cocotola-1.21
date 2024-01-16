import { ReactElement, useEffect, useRef } from 'react';

import jwt_decode, { JwtPayload } from 'jwt-decode';
import { Navigate } from 'react-router-dom';

import { useAuthStore } from '@/stores/auth';

type PrivateRouteProps = {
  element: JSX.Element;
};

export const PrivateRoute = (props: PrivateRouteProps): ReactElement => {
  console.log('PrivateRoute');
  const once = useRef(false);

  const accessToken = useAuthStore((state) => state.accessToken);
  const refreshToken = useAuthStore((state) => state.refreshToken);
  const error = useAuthStore((state) => state.error);
  const reauthenticate = useAuthStore((state) => state.reauthenticate);

  let isAccessTokenExpired = true;
  if (accessToken && accessToken != null && accessToken !== '') {
    const decoded = jwt_decode<JwtPayload>(accessToken) || null;
    if (decoded.exp) {
      isAccessTokenExpired = decoded.exp < new Date().getTime() / 1000;
    }
  }

  let isRefreshTokenExpired = true;
  if (refreshToken && refreshToken != null && refreshToken !== '') {
    const decoded = jwt_decode<JwtPayload>(refreshToken) || null;
    if (decoded.exp) {
      isRefreshTokenExpired = decoded.exp < new Date().getTime() / 1000;
    }
  }

  //   setError(authError);

  useEffect(() => {
    if (once.current === false) {
      once.current = true;
      if (refreshToken && isAccessTokenExpired && !isRefreshTokenExpired) {
        // onsole.log('xxx refreshAccessToken');
        const f = async () => {
          await reauthenticate(refreshToken);
        };
        f().catch(console.error);
      }
    }
  }, [refreshToken, isAccessTokenExpired, isRefreshTokenExpired, reauthenticate]);

  // eslint-disable-line react-hooks/exhaustive-deps

  if (error) {
    return <div>{error}</div>;
  } else if (!error && isAccessTokenExpired && !isRefreshTokenExpired) {
    return <div>Refreshing...</div>;
  } else if (isRefreshTokenExpired) {
    return <Navigate replace to={`/app/auth/login`} />;
  } else {
    console.log('children');
    return <>{props.element}</>;
  }
};
