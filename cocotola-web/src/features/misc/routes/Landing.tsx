import { Navigate } from 'react-router';

import logo from '@/assets/react.svg';
import { Head } from '@/components/head';
import { LoginForm } from '@/features/auth/components/LoginForm';
import { useAuthStore } from '@/stores/auth';

export const Landing = () => {
  const getUserInfo = useAuthStore((state) => state.getUserInfo);
  const userInfo = getUserInfo();

  if (userInfo) {
    return <Navigate to="/app" />;
  }

  return (
    <>
      <Head description="Welcome to bulletproof react" />
      <div className="bg-white h-[100vh] flex items-center">
        <div className="max-w-7xl mx-auto text-center py-12 px-4 sm:px-6 lg:py-16 lg:px-8">
          <h2 className="text-3xl font-extrabold tracking-tight text-gray-900 sm:text-4xl">
            <span className="block">Bulletproof React</span>
          </h2>
          <img src={logo} alt="react" />
          <LoginForm />
          Landing
        </div>
      </div>
    </>
  );
};
