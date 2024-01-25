import { ReactElement, useEffect, useRef } from 'react';

import jwt_decode, { JwtPayload } from 'jwt-decode';
import queryString from 'query-string';
import { Navigate } from 'react-router';

import { useAuthStore } from '@/stores/auth';
export const Callback = (): ReactElement => {
  const once = useRef(false);

  const location = window.location.search;
  const parsed = queryString.parse(location);
  const sessionState = useAuthStore((state) => state.sessionState) || '';
  const paramState = parsed ? String(parsed.state) : '';
  const code = parsed ? String(parsed.code) : '';

  const authenticate = useAuthStore((state) => state.authenticate);
  const accessToken = useAuthStore((state) => state.accessToken);
  const error = useAuthStore((state) => state.error);

  console.log('Callback', accessToken);

  let isAccessTokenExpired = true;
  if (accessToken && accessToken != null && accessToken !== '') {
    console.log('jwt_decode');
    const decoded = jwt_decode<JwtPayload>(accessToken) || null;
    if (decoded.exp) {
      isAccessTokenExpired = decoded.exp < new Date().getTime() / 1000;
    }
  }

  useEffect(() => {
    if (once.current === false) {
      once.current = true;
      if (!accessToken) {
        const f = async () => {
          await authenticate(sessionState, paramState, code);
        };
        f().catch(console.error);
      }
    }
  }, [accessToken, authenticate, sessionState, paramState, code]);

  if (error) {
    return <div>{error}</div>;
  }

  if (!accessToken) {
    return (
      <div>
        <div>Loading</div>
      </div>
    );
  }

  if (isAccessTokenExpired) {
    return <div>Expired</div>;
  }

  return <Navigate replace to="/" />;
};
