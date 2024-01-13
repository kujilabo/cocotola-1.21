import { Suspense } from 'react';

import { Routes, Route, Outlet } from 'react-router-dom';

import { MainLayout } from '@/components/layout';
import { PrivateRoute } from '@/components/PrivateRoute';
import { Landing } from '@/features/misc';
import { lazyImport } from '@/utils/lazyImport';

const { Dashboard } = lazyImport(() => import('@/features/misc'), 'Dashboard');
const { AuthRoutes } = lazyImport(() => import('@/features/auth'), 'AuthRoutes');
const { WorkbookRoutes } = lazyImport(() => import('@/features/workbook'), 'WorkbookRoutes');
const { DiscussionsRoutes } = lazyImport(
  () => import('@/features/discussions'),
  'DiscussionsRoutes'
);

const App = () => {
  console.log('app');
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
export const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<Landing />} />
      <Route path="/app/auth/*" element={<AuthRoutes />} />
      <Route path="/app" element={<PrivateRoute element={<App />} />}>
        <Route path="/app/dis/*" element={<DiscussionsRoutes />} />
        <Route path="/app/workbook/*" element={<WorkbookRoutes />} />
        <Route path="/app" element={<Dashboard />} />
      </Route>
    </Routes>
  );
};
