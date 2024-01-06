import { ReactElement, useEffect, useRef } from 'react';

import jwt_decode, { JwtPayload } from 'jwt-decode';
import queryString from 'query-string';
import { Navigate } from 'react-router';

import { useAuthStore } from '@/stores/auth';
export const Callback = (): ReactElement => {
  const location = window.location.search;
  const fetchToken = useAuthStore((state) => state.fetchToken);
  const accessToken = useAuthStore((state) => state.accessToken);
  const error = useAuthStore((state) => state.error);
  const parsed = queryString.parse(location);
  const code = parsed ? String(parsed.code) : '';
  const once = useRef(false);
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
    if (!accessToken) {
      if (once.current === false) {
        once.current = true;

        const f = async () => {
          await fetchToken(code);
        };
        f().catch(console.error);
      }
    }
  }, [accessToken, fetchToken, code]);

  if (error) {
    return <div>{error}</div>;
  }

  if (!accessToken) {
    <div>
      <div>Loading</div>
    </div>;
  }

  if (isAccessTokenExpired) {
    return <div>Expired</div>;
  }

  return <Navigate replace to="/" />;
};
