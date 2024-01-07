import { useRoutes } from 'react-router-dom';

import { Landing } from '@/features/misc';
import { useAuthStore } from '@/stores/auth';

import { privateRoutes } from './private';
import { publicRoutes } from './public';

export const AppRoutes = () => {
  const userInfo = useAuthStore((state) => state.userInfo);

  const commonRoutes = [{ path: '/', element: <Landing /> }];

  const routes = userInfo ? privateRoutes : publicRoutes;

  const element = useRoutes([...routes, ...commonRoutes]);

  return <>{element}</>;
};
