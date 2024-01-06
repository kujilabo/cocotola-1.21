import { useState } from 'react';

import { Provider } from 'react-redux';

import viteLogo from '~/vite.svg';

import { store } from '@/app/store';
import reactLogo from '@/assets/react.svg';

import '@/App.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

// import AuthRoutes from '@/routes/auth_routes'
// import UnauthRoutes from '@/routes/unauth_routes'
import { Providers } from '@/components/Providers';
import { RouteAuthGuard } from '@/components/RouteAuthGuard';
import { Callback } from '@/features/auth/components/Callback';
import { Login } from '@/features/auth/components/Login';
import { Callback as CallbackZustand } from '@/features/auth_zustand/components/Callback';
import { Login as LoginZustand } from '@/features/auth_zustand/components/Login';
import { Test } from '@/features/private_workbook/components/Test';
import { Home } from '@/pages/Home';
import { Layout } from '@/pages/Layout';
// import { Login } from '@/pages/Login';
import { NotFound } from '@/pages/NotFound';
import { Page1 } from '@/pages/Page1';
import { Page2 } from '@/pages/Page2';
import { Page3, Page3Child1, Page3Child2 } from '@/pages/Page3';
import { RoleType } from '@/types';
function App() {
  return (
    <Provider store={store}>
      <BrowserRouter>
        <Routes>
          {/* <Route path="/" element={<Layout />} >
						<Route index element={<Home />} />
          <Route path="/page1" element={<Page1 />}/>

<Route path="/page2" element={
	  <RouteAuthGuard component={<Page2 />} redirect="/login" />} />

<Route path="/page3" element={
                <RouteAuthGuard component={<Page3 />} redirect="/login" 
                        allowroles={[RoleType.Admin, RoleType.Manager]} />} >
            <Route index element={<Page3Child1 />} />
            <Route path="child1" element={<Page3Child1 />} />
            <Route path="child2" element={<Page3Child2 />} />
          </Route>
						<Route path="/login" element={<Login />} />
          <Route path="*" element={<NotFound />} />
					</Route> */}
          <Route path="/" element={<Layout />}>
            <Route path="/app/login" element={<Login />} />
            <Route path="/app/login-zustand" element={<LoginZustand />} />
            <Route path="/app/callback" element={<Callback />} />
            <Route path="/app/callback-zustand" element={<CallbackZustand />} />
            <Route path="/test" element={<Test />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </Provider>
  );
}

export default App;
