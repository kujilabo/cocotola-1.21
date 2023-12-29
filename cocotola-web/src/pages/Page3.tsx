import React from 'react';

import { Outlet, Link } from 'react-router-dom';

export const Page3: React.VFC = () => {
  return (
    <>
      <h3>Page 3</h3>
      <ul>
        <li>
          <Link to="child1">show child1 page1</Link>
        </li>
        <li>
          <Link to="child2">show child2 page2</Link>
        </li>
      </ul>
      <Outlet />
    </>
  );
};

export const Page3Child1: React.VFC = () => {
  return <h3>Page 3 Child1</h3>;
};
export const Page3Child2: React.VFC = () => {
  return <h3>Page 3 Child2</h3>;
};
