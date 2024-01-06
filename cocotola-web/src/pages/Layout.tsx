import React from 'react';

import { Outlet, Link } from 'react-router-dom';

import { useAuthStore } from '@/stores/auth';

export const Layout: React.VFC = () => {
  const resetTokens = useAuthStore((state) => state.resetTokens);
  return (
    <>
      <h3>Layout</h3>
      <ul>
        <li>
          <Link to="/">show home</Link>
        </li>
        <li>
          <Link to="page1">show page1</Link>
        </li>
        <li>
          <Link to="test">Test</Link>
        </li>
        <li>
          <Link to="app/login">Login</Link>
        </li>
        <li>
          <Link to="app/login-zustand">ZustandLogin</Link>
        </li>
        <li>
          <Link to="page3">show page3</Link>
          <span> ( Admin または Manager 権限のユーザーのみ可 )</span>
        </li>
      </ul>
      <button onClick={resetTokens}>Logout</button>

      <Outlet />
    </>
  );
};
