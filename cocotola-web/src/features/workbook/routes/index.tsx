import { Route, Routes } from 'react-router-dom';

import { WorkbookList } from '@/features/workbook/components/WorkbookList';
import { WorkbookView } from '@/features/workbook/components/WorkbookView';

export const WorkbookRoutes = () => {
  return (
    <Routes>
      <Route path="" element={<WorkbookList />} />
      <Route path=":_workbookId" element={<WorkbookView />} />
    </Routes>
  );
};
