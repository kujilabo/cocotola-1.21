// import { Dialog, Menu, Transition } from '@headlessui/react';
import * as React from 'react';

import { Dialog, Transition } from '@headlessui/react';
import {
  UserIcon,
  FolderIcon,
  HomeIcon,
  // MenuAlt2Icon,
  // UsersIcon,
  // XIcon,
} from '@heroicons/react/24/outline';
import { clsx } from 'clsx';
import { NavLink, Link } from 'react-router-dom';

import logo from '@/assets/react.svg';
// import { useAuth } from '@/lib/auth';
// import { useAuthorization, ROLES } from '@/lib/authorization';
import { useAuthStore } from '@/stores/auth';
import { Box, Flex, Text, Button, Stack } from '@chakra-ui/react';
import {
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  MenuItemOption,
  MenuGroup,
  MenuOptionGroup,
  Spacer,
  MenuDivider,
} from '@chakra-ui/react';
const NavBar = () => {
  const [isOpen, setIsOpen] = React.useState(false);
  const logout = useAuthStore((state) => state.resetTokens);
  const toggle = () => setIsOpen(!isOpen);

  return (
    <NavBarContainer>
      <Logo />
      <Logo />
      <Spacer />
      <Menu>
        <MenuButton as={Button} colorScheme="pink">
          Profile
        </MenuButton>
        <MenuList>
          <MenuGroup title="Profile">
            <MenuItem>My Account</MenuItem>
            <MenuItem onClick={() => logout()}>Logout </MenuItem>
          </MenuGroup>
          <MenuDivider />
          <MenuGroup title="Help">
            <MenuItem>Docs</MenuItem>
            <MenuItem>FAQ</MenuItem>
          </MenuGroup>
        </MenuList>
      </Menu>
      {/* 
      <MenuLinks isOpen={isOpen} /> */}
    </NavBarContainer>
  );
};

const Logo = () => {
  return (
    <Box w="100px" color={'white'}>
      <Text fontSize="lg" fontWeight="bold">
        cocotola
      </Text>
    </Box>
  );
};

type NavBarContainerProps = {
  children: React.ReactNode;
};

const NavBarContainer = ({ children }: NavBarContainerProps) => {
  return (
    <Flex
      as="nav"
      // align="center"
      // justify="space-between"
      // wrap="wrap"
      // w="100%"
      mb={0}
      p={4}
      background={'blue.500'}
      // color={'white'}
    >
      {children}
    </Flex>
  );
};
type MainLayoutProps = {
  children: React.ReactNode;
};

export const MainLayout = ({ children }: MainLayoutProps) => {
  const [sidebarOpen, setSidebarOpen] = React.useState(false);

  return (
    <div>
      <NavBar />
      <div>
        <main>{children}</main>
      </div>
    </div>
  );
};
