import axios from 'axios';
import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';

import { WorkbookWithProblems } from 'src/features/workbook/types';

import { backendCoreUrl, extractErrorMessage } from '@/lib/base';
import { jsonHeaders } from '@/lib/util';
import { useAuthStore, isTokenExpired } from '@/stores/auth';

type State = {
  workbooks: { [key: number]: WorkbookWithProblems };
  state: string;
  error: string | null;
};
type Action = {
  retrieveWorkbook: (id: number) => Promise<void>;
};
type WorkbookRetrieveResponse = {
  id: number;
  name: string;
  description: string;
  problems: {
    type: string;
    properties: { [key: string]: string };
  }[];
};

export const useWorkbookRetrieveStore = create<State & Action>()(
  devtools(
    persist(
      (set, get) => ({
        workbooks: {},
        state: 'idle',
        error: null,
        retrieveWorkbook: async (id: number): Promise<void> => {
          set({
            state: 'pending',
            error: null,
          });
          const url = `${backendCoreUrl}/v1/workbook/${id}`;
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
                const workbook = resp.data as WorkbookRetrieveResponse;
                console.log(workbook);
                const workbooks = get().workbooks;
                set({
                  state: 'fulfilled',
                  workbooks: { ...workbooks, [workbook.id]: workbook },
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
        name: 'workbook-retrieve-storage',

        partialize: (state) => ({
          workbooks: state.workbooks,
        }),
      }
    )
  )
);
