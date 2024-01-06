import { ReactElement, useEffect, useRef } from 'react';

import axios from 'axios';
import queryString from 'query-string';
import { useNavigate } from 'react-router-dom';
import { create } from 'zustand';
import { persist, devtools } from 'zustand/middleware';

import jwt_decode, { JwtPayload } from 'jwt-decode';
import { Navigate } from 'react-router';
import { backendAuthUrl, extractErrorMessage } from '@/lib/base';
type State = {
  accessToken: string | null;
  refreshToken: string | null;
  error: string | null;
};

type Action = {
  fetchToken: (code: string) => void;
};
export type GoogleAuthorizeParameter = {
  organizationName: string;
  code: string;
};

type GoogleAuthorizeResponse = {
  accessToken: string;
  refreshToken: string;
};
const useAuthStore = create<State & Action>()(devtools(
  persist(
    (set) => ({
      bears: 0,
      fishies: {},
      accessToken: null,
      refreshToken: null,
      error: null,
      fetchToken: async (code: string) => {
        set({ error: null });
        console.log('aaa');
        const param: GoogleAuthorizeParameter = {
          organizationName: 'cocotola',
          code: code,
        };
        await axios
          .post(`${backendAuthUrl}/v1/google/authorize`, param)
          .then((resp) => {
            console.log('callback then');
            const tokens = resp.data as GoogleAuthorizeResponse;
            set({
              accessToken: tokens.accessToken,
               refreshToken: tokens.refreshToken,
               });
          })
          .catch((err: Error) => {
            console.log('callback err');
            const errorMessage = extractErrorMessage(err);
            //   arg.postFailureProcess(errorMessage);
            //   return thunkAPI.rejectWithValue(errorMessage);
            set({ error: errorMessage });
          });
      },
    }),
    {
      name: 'auth-storage',

      partialize: (state) => ({ accessToken: state.accessToken, refreshToken: state.refreshToken }),
    }
  )
));
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
        fetchToken(code);
      }
    }
  }, [accessToken]);

  if (error) {
    return <div>{error}</div>
  }


  if (!accessToken) {
    <div>
      <div>Loading</div>
    </div>
  }

  if (isAccessTokenExpired) {
    return <div>Expired</div>;
  }

  return <Navigate replace to="/" />;
};
