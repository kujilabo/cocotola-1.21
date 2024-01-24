import { clientId, frontendUrl } from '@/lib/base';
import { useAuthStore } from '@/stores/auth';

export const LoginForm = () => {
  const resetTokens = useAuthStore((state) => state.resetTokens);
  const generateState = useAuthStore((state) => state.generateState);
  const handleLogin = () => {
    resetTokens();
    const callback = (state: string): void => {
      let url = 'https://accounts.google.com/o/oauth2/auth';
      url += '?client_id=';
      url += clientId;
      url += '&redirect_uri=';
      url += frontendUrl;
      url += '/app/auth/callback';
      url += '&scope=profile email';
      url += '&response_type=';
      url += 'code';
      url += '&access_type=';
      url += 'offline';
      url += '&state=';
      url += state;
      window.location.href = url;
    };
    const f = async () => {
      await generateState(callback);
    };
    f().catch(console.error);
  };

  return (
    <div>
      <h1>Login</h1>
      <button onClick={handleLogin}>Login</button>
    </div>
  );
};
