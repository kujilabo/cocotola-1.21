import axios from 'axios';
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

type State = {
  accessToken: string | null;
  refreshToken: string | null;
  error: string | null;
};
type Action = {
  resetTokens: () => void;
  authenticate: (code: string) => Promise<void>;
  reauthenticate: (refreshToken: string) => Promise<void>;
};
export const useAuthStore = create<State & Action>()(
  devtools(
    persist(
      (set) => ({
        accessToken: null,
        refreshToken: null,
        error: null,
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
        }),
      }
    )
  )
);
