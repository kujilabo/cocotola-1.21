import { configureStore, ThunkAction, Action } from '@reduxjs/toolkit';
import { combineReducers } from 'redux';

import authReducer from '@/features/auth/api/auth';
import workbookFindReducer from '@/features/private_workbook/api/workbook_find';
import workbookTestReducer from '@/features/private_workbook/api/test';
export const rootReducer = combineReducers({
  // ...reducers,
  auth: authReducer,
  // // router: routerReducer,
  // stat: statReducer,

  workbookTest: workbookTestReducer,
  workbookFind: workbookFindReducer,
  // workbookAdd: workbookAddReducer,
  // workbookUpdate: workbookUpdateReducer,
  // workbookGet: workbookGetReducer,
  // workbookRemove: workbookRemoveReducer,
  // problemFind: problemFindReducer,
  // problemAdd: problemAddReducer,
  // problemUpdate: problemUpdateReducer,
  // // problemGet: problemGetReducer,
  // problemRemove: problemRemoveReducer,
  // problemImport: problemImportReducer,
  // recordbookGet: recordbookGetReducer,
  // recordAdd: recordAddReducer,
  // audio: audioReducer,
  // // plugin
  // translationFind: translationFindReducer,
  // translationGet: translationGetReducer,
  // translationGetList: translationGetListReducer,
  // translationAdd: translationAddReducer,
  // translationUpdate: translationUpdateReducer,
  // translationRemove: translationRemoveReducer,
  // translationImport: translationImportReducer,
  // translationExport: translationExportReducer,
  // tatoebaImport: tatoebaImportReducer,
  // tatoebaSentenceFindSlice: tatoebaSentenceFindReducer,
});
export const store = configureStore({
  reducer: rootReducer,
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;
export type BaseThunkApiConfig = {
  dispatch: AppDispatch;
  state: RootState;
};
