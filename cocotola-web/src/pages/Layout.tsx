import React from 'react';

import { Outlet, Link } from 'react-router-dom';

export const Layout: React.VFC = () => {
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
          <Link to="page2">show page2</Link>
          <span> ( 認証済みのユーザーなら可 )</span>
        </li>
        <li>
          <Link to="app/login">Login</Link>
        </li>
        <li>
          <Link to="page3">show page3</Link>
          <span> ( Admin または Manager 権限のユーザーのみ可 )</span>
        </li>
      </ul>

      <Outlet />
    </>
  );
};
