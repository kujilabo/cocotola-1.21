import { ReactElement } from 'react';

import jwt_decode, { JwtPayload } from 'jwt-decode';
import queryString from 'query-string';
import { Navigate } from 'react-router';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { AppDimmer } from '@/components/AppDimmer';
import {
  selectAccessToken,
  selectAuthLoading,
  selectAuthFailed,
  googleAuthorize,
} from '@/features/auth/api/auth';
import { emptyFunction } from '@/lib/util';

export const Callback = (): ReactElement => {
  const dispatch = useAppDispatch();
  const accessToken = useAppSelector(selectAccessToken);
  console.log('decode acc', accessToken);
  let isAccessTokenExpired = true;
  if (accessToken && accessToken != null && accessToken !== '') {
    console.log('jwt_decode');
    const decoded = jwt_decode<JwtPayload>(accessToken) || null;
    if (decoded.exp) {
      isAccessTokenExpired = decoded.exp < new Date().getTime() / 1000;
    }
  }

  const authLoading = useAppSelector(selectAuthLoading);
  const authFailed = useAppSelector(selectAuthFailed);
  const location = window.location.search;
  console.log('Callback', authLoading, authFailed, isAccessTokenExpired);
  if (authFailed) {
    return <div>Failed</div>;
  } else if (authLoading === false && isAccessTokenExpired) {
    const parsed = queryString.parse(location);
    const code = parsed ? String(parsed.code) : '';

    const f = async () => {
      await dispatch(
        googleAuthorize({
          param: {
            organizationName: 'cocotola',
            code: code,
          },
          postSuccessProcess: emptyFunction,
          postFailureProcess: (error: string) => {
            console.log('callback error', error);
            return;
          },
        })
      );
    };
    f().catch(console.error);
    return <AppDimmer />;
  } else if (!isAccessTokenExpired) {
    return <Navigate replace to="/" />;
  } else {
    return <AppDimmer />;
  }
};
