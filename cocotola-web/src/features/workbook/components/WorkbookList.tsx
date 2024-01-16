import { useEffect, useRef } from 'react';

import {
  Card,
  CardHeader,
  CardBody,
  Heading,
  StackDivider,
  Stack,
  Box,
  Text,
  Container,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
} from '@chakra-ui/react';
import { Link } from 'react-router-dom';

import { useWorkbookFindStore } from '../api/workbook_find';
import { Workbook } from '../types';
export const WorkbookList = (): JSX.Element => {
  const once = useRef(false);
  const workbooks = useWorkbookFindStore((state) => state.workbooks);
  const state = useWorkbookFindStore((state) => state.state);
  const findWorkbooks = useWorkbookFindStore((state) => state.findWorkbooks);

  useEffect(() => {
    if (once.current === false) {
      once.current = true;
      if (state === 'idle') {
        const f = async () => {
          await findWorkbooks();
        };
        f().catch(console.error);
      }
    }
  }, [state, findWorkbooks]);
  return (
    <Container>
      <Breadcrumb>
        <BreadcrumbItem>
          <BreadcrumbLink as={Link} to="/">
            Home
          </BreadcrumbLink>
        </BreadcrumbItem>

        <BreadcrumbItem>
          <BreadcrumbLink href="#">Docs</BreadcrumbLink>
        </BreadcrumbItem>

        <BreadcrumbItem isCurrentPage>
          <BreadcrumbLink href="#">Breadcrumb</BreadcrumbLink>
        </BreadcrumbItem>
      </Breadcrumb>
      {workbooks.map((workbook: Workbook) => {
        return (
          <Card key={workbook.id}>
            {workbook.createdAt}
            <CardHeader>
              <Heading size="md">Client Report</Heading>
            </CardHeader>

            <CardBody>
              <Stack divider={<StackDivider />} spacing="4">
                <Box>
                  <Heading size="xs" textTransform="uppercase">
                    Summary
                  </Heading>
                  <Text pt="2" fontSize="sm">
                    View a summary of all your clients over the last month.
                  </Text>
                </Box>
                <Box>
                  <Heading size="xs" textTransform="uppercase">
                    Overview
                  </Heading>
                  <Text pt="2" fontSize="sm">
                    Check out the overview of your clients.
                  </Text>
                </Box>
                <Box>
                  <Heading size="xs" textTransform="uppercase">
                    Analysis
                  </Heading>
                  <Text pt="2" fontSize="sm">
                    See a detailed analysis of all your business clients.
                  </Text>
                </Box>
              </Stack>
            </CardBody>
          </Card>
        );
      })}
    </Container>
  );
};
