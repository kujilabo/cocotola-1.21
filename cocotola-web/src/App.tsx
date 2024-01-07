import '@/App.css';
import { BrowserRouter, Routes } from 'react-router-dom';

// import { PrivateRoute } from '@/components/PrivateRoute';
// import { Callback } from '@/features/auth/components/Callback';
// import { Login } from '@/features/auth/components/Login';
// import { Test } from '@/features/private_workbook/components/Test';
// import { Layout } from '@/pages/Layout';

import { AppRoutes } from '@/routes';
function App() {
  return (
    <BrowserRouter>
      <AppRoutes />
      {/* <Route path="/" element={<Layout />}>
          <Route path="/app/login" element={<Login />} />
          <Route path="/app/callback" element={<Callback />} />
          <Route path="/test" element={<PrivateRoute element={<Test />} />} />
        </Route> */}
    </BrowserRouter>
  );
}

export default App;
