import React, { ReactElement, useEffect } from 'react';

import { Button, ButtonGroup } from '@chakra-ui/react';
import { useSelector, useDispatch } from 'react-redux';

import { useAppSelector, useAppDispatch } from '@/app/hooks';
import { testABC, selectTestStatus } from '@/features/private_workbook/api/test';
import { clientId, frontendUrl } from '@/lib/base';
import { emptyFunction } from '@/lib/util';

console.log(clientId, 'clientId');
export const Test = (): ReactElement => {
  const dispatch = useAppDispatch();
  console.log('test');

  const status = useAppSelector(selectTestStatus);

  useEffect(() => {
    if (status === 'idle') {
      dispatch(
        testABC({
          param: {
            workbookId: 1,
          },
        })
      );
    }
  }, [status, dispatch]);

  return <div>Test</div>;
};
