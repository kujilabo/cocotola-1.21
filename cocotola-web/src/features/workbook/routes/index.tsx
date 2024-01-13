import { Route, Routes } from 'react-router-dom';

import { Workbooks } from '@/features/workbook/components/Workbooks';

export const WorkbookRoutes = () => {
  return (
    <Routes>
      <Route path="" element={<Workbooks />} />
    </Routes>
  );
};
