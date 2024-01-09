// import { useRoutes } from 'react-router-dom';

import { Landing } from '@/features/misc';
// import { useAuthStore } from '@/stores/auth';

import { Routes, Route } from 'react-router-dom';
// import { ChakraProvider } from '@chakra-ui/react';
// import { privateRoutes } from './private';
// import { publicRoutes } from './public';
import { PrivateRoute } from '@/components/PrivateRoute';
import { Callback } from '@/features/auth/components/Callback';
import { Login } from '@/features/auth/routes/Login';
import { Test } from '@/features/private_workbook/components/Test';

export const AppRoutes = () => {
  // const getUserInfo = useAuthStore((state) => state.getUserInfo);
  // const userInfo = getUserInfo();

  // const commonRoutes = [{ path: '/', element: <Landing /> }];

  // const routes = userInfo ? privateRoutes : publicRoutes;

  // const element = useRoutes([...routes, ...commonRoutes]);

  // return <>{element}</>;
  return (
    <Routes>
      <Route path="/" element={<Landing />} />
      <Route path="/app/login" element={<Login />} />
      <Route path="/app/callback" element={<Callback />} />
      <Route path="/test" element={<PrivateRoute element={<Test />} />} />
    </Routes>
  );
};
