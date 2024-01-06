import '@/App.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import { Callback as CallbackZustand } from '@/features/auth_zustand/components/Callback';
import { Login as LoginZustand } from '@/features/auth_zustand/components/Login';
import { Test } from '@/features/private_workbook/components/Test';
import { Layout } from '@/pages/Layout';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route path="/app/login-zustand" element={<LoginZustand />} />
          <Route path="/app/callback-zustand" element={<CallbackZustand />} />
          <Route path="/test" element={<Test />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
