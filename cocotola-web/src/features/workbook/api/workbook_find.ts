import axios from 'axios';
import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';

import { Workbook } from 'src/features/workbook/types';

import { backendCoreUrl, extractErrorMessage } from '@/lib/base';
import { jsonHeaders } from '@/lib/util';
import { useAuthStore, isTokenExpired } from '@/stores/auth';

type State = {
  workbooks: Workbook[];
  state: string;
  error: string | null;
};
type Action = {
  findWorkbooks: () => Promise<void>;
};
type WorkbookFindResponse = {
  accessToken: string;
  refreshToken: string;
};

export const useWorkbookFindStore = create<State & Action>()(
  devtools(
    persist(
      (set) => ({
        workbooks: [],
        state: 'idle',
        error: null,
        findWorkbooks: async (): Promise<void> => {
          set({
            state: 'pending',
            error: null,
          });
          const url = `${backendCoreUrl}/v1/workbook`;
          const { refreshToken, reauthenticate } = useAuthStore.getState();
          const isExpired = isTokenExpired(refreshToken);
          if (isExpired) {
            set({ state: 'unauthenticated' });
            return;
          }
          await reauthenticate(refreshToken || '').then(() => {
            const { accessToken } = useAuthStore.getState();
            return axios
              .get(url, { headers: jsonHeaders(accessToken || ''), data: {} })
              .then((resp) => {
                console.log('callback then');
                const tokens = resp.data as WorkbookFindResponse;
                console.log(tokens);
                set({
                  state: 'fulfilled',
                });
              })
              .catch((err: Error) => {
                console.log('callback err');
                const errorMessage = extractErrorMessage(err);
                //   arg.postFailureProcess(errorMessage);
                //   return thunkAPI.rejectWithValue(errorMessage);
                set({
                  state: 'rejected',
                  error: errorMessage,
                });
              });
          });
        },
      }),
      {
        name: 'workbook-find-storage',

        partialize: (state) => ({
          workbooks: state.workbooks,
        }),
      }
    )
  )
);
