import { clientId, frontendUrl } from '@/lib/base';
import { useAuthStore } from '@/stores/auth';

export const Login = () => {
  const resetTokens = useAuthStore((state) => state.resetTokens);
  const handleLogin = () => {
    resetTokens();
    let url = 'https://accounts.google.com/o/oauth2/auth';
    url += '?client_id=';
    url += clientId;
    url += '&redirect_uri=';
    url += frontendUrl;
    url += '/app/callback';
    url += '&scope=profile email';
    url += '&response_type=';
    url += 'code';
    url += '&access_type=';
    url += 'offline';
    url += '&state=';
    url += 'STATE';
    window.location.href = url;
  };

  return (
    <div>
      <h1>Login</h1>
      <button onClick={handleLogin}>Login</button>
    </div>
  );
};
