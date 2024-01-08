import * as React from 'react';

import { Link } from 'react-router-dom';

import logo from '@/assets/react.svg';
// import { Link } from '@/components/Elements';
import { Head } from '@/components/head';

type LayoutProps = {
  children: React.ReactNode;
  title: string;
};

export const Layout = ({ children, title }: LayoutProps) => {
  return (
    <>
      <Head title={title} />
      <div>
        <div>
          <div>
            <Link to="/">
              <img className="h-24 w-auto" src={logo} alt="Workflow" />
            </Link>
          </div>

          <h2>{title}</h2>
        </div>

        <div>
          <div>{children}</div>
        </div>
      </div>
    </>
  );
};
