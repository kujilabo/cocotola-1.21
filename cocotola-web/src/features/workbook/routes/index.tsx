import { Route, Routes } from 'react-router-dom';

import { WorkbookList } from '@/features/workbook/components/WorkbookList';

export const WorkbookRoutes = () => {
  return (
    <Routes>
      <Route path="" element={<WorkbookList />} />
    </Routes>
  );
};
