import { Route, Routes } from 'react-router-dom';

import { Callback } from '@/features/auth/components/Callback';
import { Login } from '@/features/auth/routes/Login';

export const AuthRoutes = () => {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/callback" element={<Callback />} />
    </Routes>
  );
};
