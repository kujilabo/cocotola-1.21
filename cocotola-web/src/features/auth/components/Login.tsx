import { ReactElement } from 'react';

import { Button, ButtonGroup } from '@chakra-ui/react';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { selectRedirectUrl, redirectTo } from '@/features/auth/api/auth';
import { clientId, frontendUrl } from '@/lib/base';
import { emptyFunction } from '@/lib/util';

console.log(clientId, 'clientId');
export const Login = (): ReactElement => {
  console.log('a');
  const dispatch = useAppDispatch();
  const redirectUrl = useAppSelector(selectRedirectUrl);
  const googleAuth = () => {
    let url = 'https://accounts.google.com/o/oauth2/auth';
    url += '?client_id=';
    url += clientId;
    url += '&redirect_uri=';
    url += frontendUrl;
    url += '/app/callback';
    url += '&scope=profile email';
    url += '&response_type=';
    url += 'code';
    url += '&access_type=';
    url += 'offline';
    url += '&state=';
    url += 'STATE';
    console.log(url);
    dispatch(redirectTo({ url: url }));
  };

  console.log('redirectUrl', redirectUrl);
  if (redirectUrl && redirectUrl !== '') {
    console.log('redirect');
    window.location.href = redirectUrl;
    // return <Navigate replace to={redirectUrl} />;
  }

  return (
    <div>
      <Button color="teal" onClick={googleAuth}>
        Sign in with Google
      </Button>
    </div>
  );
};
