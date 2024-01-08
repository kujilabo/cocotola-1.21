// import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
// import axios from 'axios';

// import { RootState, BaseThunkApiConfig } from '@/app/store';
// import { refreshAccessToken } from '@/features/auth/api/auth';
// import { WorkbookModel } from '@/features/private_workbook/models/model';
// import { backendCoreUrl, extractErrorMessage } from '@/lib/base';
// import { jsonRequestConfig } from '@/lib/util';

// const baseUrl = `${backendCoreUrl}/v1/private/workbook`;

// export type TestParameter = {
//   workbookId: number;
// };

// export type TestArg = {
//   param: TestParameter;
//   // postSuccessProcess: () => void;
//   // postFailureProcess: (error: string) => void;
// };
// type TestResponse = {
//   id: number;
// };
// type TestResult = {
//   response: TestResponse;
// };

// export const testABC = createAsyncThunk<TestResult, TestArg, BaseThunkApiConfig>(
//   'private/workbook/test',
//   async (arg: TestArg, thunkAPI) => {
//     const url = `${baseUrl}/test`;
//     const { refreshToken } = thunkAPI.getState().auth;
//     return await thunkAPI.dispatch(refreshAccessToken({ refreshToken: refreshToken })).then(() => {
//       // onsole.log('accessToken1');
//       const { accessToken } = thunkAPI.getState().auth;
//       // onsole.log('accessToken', accessToken);
//       return axios
//         .get(url, jsonRequestConfig(accessToken))
//         .then((resp) => {
//           // onsole.log('the1', resp);
//           const response = resp.data as TestResponse;
//           // onsole.log('the2', response);
//           // arg.postSuccessProcess();
//           return {
//             response: response,
//           } as TestResult;
//         })
//         .catch((err: Error) => {
//           const errorMessage = extractErrorMessage(err);
//           // arg.postFailureProcess(errorMessage);
//           return thunkAPI.rejectWithValue(errorMessage);
//         });
//     });
//   }
// );

// export interface TestState {
//   loading: boolean;
//   failed: boolean;
//   status: string;
// }

// const initialState: TestState = {
//   loading: false,
//   failed: false,
//   status: 'idle',
// };

// export const workbookTestSlice = createSlice({
//   name: 'workbook_test',
//   initialState: initialState,
//   reducers: {},
//   extraReducers: (builder) => {
//     builder
//       .addCase(testABC.pending, (state) => {
//         state.loading = true;
//         state.status = 'loading';
//       })
//       .addCase(testABC.fulfilled, (state, action) => {
//         state.loading = false;
//         state.failed = false;
//         state.status = 'fulfilled';
//         // onsole.log('workbooksLoadedMap', state.workbooksLoadedMap);
//       })
//       .addCase(testABC.rejected, (state) => {
//         // onsole.log('rejected', action);
//         state.loading = false;
//         state.failed = true;
//         state.status = 'rejected';
//       });
//   },
// });

// export const selectTestStatus = (state: RootState): string => state.workbookTest.status;

// export const selectWorkbookFindLoading = (state: RootState): boolean => state.workbookFind.loading;

// export const selectWorkbookFindFailed = (state: RootState): boolean => state.workbookFind.failed;

// export const selectWorkbooksLoadedMap = (state: RootState): { [key: string]: boolean } =>
//   state.workbookFind.workbooksLoadedMap;

// export const selectWorkbooksMap = (state: RootState): { [key: string]: WorkbookModel[] } =>
//   state.workbookFind.workbooksMap;

// export default workbookTestSlice.reducer;
