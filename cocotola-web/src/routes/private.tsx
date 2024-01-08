import { Suspense } from 'react';

import { Navigate, Outlet } from 'react-router-dom';

// import { Spinner } from '@/components/Elements';
import { MainLayout } from '@/components/layout';
import { lazyImport } from '@/utils/lazyImport';

import { PrivateRoute } from '@/components/PrivateRoute';
const { DiscussionsRoutes } = lazyImport(
  () => import('@/features/discussions'),
  'DiscussionsRoutes'
);
const { Dashboard } = lazyImport(() => import('@/features/misc'), 'Dashboard');
// const { Profile } = lazyImport(() => import('@/features/users'), 'Profile');
// const { Users } = lazyImport(() => import('@/features/users'), 'Users');

const App = () => {
  return (
    <MainLayout>
      <Suspense
        fallback={
          <div className="h-full w-full flex items-center justify-center">
            {/* <Spinner size="xl" /> */}
            SPINNER
          </div>
        }
      >
        <Outlet />
      </Suspense>
    </MainLayout>
  );
};

export const privateRoutes = [
  {
    path: '/app',
    element: <App />,
    // element: <PrivateRoute element={<App />} />,
    // children: [
    //   { path: '/app/discussions', element: <DiscussionsRoutes /> },
    //   // { path: '/users', element: <Users /> },
    //   // { path: '/profile', element: <Profile /> },
    //   { path: '/app', element: <Dashboard /> },
    //   { path: '/app/*', element: <Navigate to="." /> },
    // ],
  },
];
