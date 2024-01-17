import axios from 'axios';
import jwt_decode, { JwtPayload } from 'jwt-decode';
import { create } from 'zustand';
import { persist, devtools } from 'zustand/middleware';

import { backendAuthUrl, extractErrorMessage } from '@/lib/base';

type GoogleAuthorizeParameter = {
  organizationName: string;
  code: string;
};
type GoogleAuthorizeResponse = {
  accessToken: string;
  refreshToken: string;
};

type RefreshTokenParameter = {
  refreshToken: string;
};
type RefreshTokenResponse = {
  accessToken: string;
};

interface AppJwtPayload extends JwtPayload {
  loginId: string;
  username: string;
}

type UserInfo = {
  loginId: string;
  username: string;
};

type State = {
  accessToken: string | null;
  refreshToken: string | null;
  error: string | null;
};
type Action = {
  resetTokens: () => void;
  authenticate: (code: string) => Promise<void>;
  reauthenticate: (refreshToken: string) => Promise<void>;
  getUserInfo: () => UserInfo | null;
};

const decodeJwt = (accessToken: string | null): UserInfo | null => {
  if (!accessToken) {
    return null;
  }

  const decoded = jwt_decode<AppJwtPayload>(accessToken) || null;
  if (!decoded) {
    return null;
  }

  const loginId = decoded ? decoded.loginId : '';
  const username = decoded ? decoded.username : '';
  return { loginId: loginId, username: username };
};

export const useAuthStore = create<State & Action>()(
  devtools(
    persist(
      (set, get) => ({
        accessToken: null,
        refreshToken: null,
        error: null,
        getUserInfo: (): UserInfo | null => {
          const accessToken = get().accessToken;
          return decodeJwt(accessToken);
        },
        resetTokens: (): void => {
          set({
            accessToken: null,
            refreshToken: null,
          });
        },
        authenticate: async (code: string): Promise<void> => {
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
        reauthenticate: async (refreshToken: string): Promise<void> => {
          const accessToken = get().accessToken;
          if (!isTokenExpired(accessToken)) {
            console.log('token is not expired');
            return new Promise(function (resolve) {
              resolve();
            });
          }
          const param: RefreshTokenParameter = {
            refreshToken: refreshToken,
          };
          await axios
            .post(`${backendAuthUrl}/v1/auth/refresh_token`, param)
            .then((resp) => {
              console.log('callback then');
              const token = resp.data as RefreshTokenResponse;
              set({
                accessToken: token.accessToken,
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

        partialize: (state) => ({
          accessToken: state.accessToken,
          refreshToken: state.refreshToken,
          userInfo: decodeJwt(state.accessToken),
        }),
      }
    )
  )
);

export const isTokenExpired = (token: string | null): boolean => {
  let isExpired = true;
  if (token && token != null && token !== '') {
    const decoded = jwt_decode<JwtPayload>(token) || null;
    if (decoded.exp) {
      isExpired = decoded.exp < new Date().getTime() / 1000;
    }
  }
  return isExpired;
};
