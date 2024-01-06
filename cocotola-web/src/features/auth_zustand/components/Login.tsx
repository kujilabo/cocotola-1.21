import { clientId, frontendUrl } from '@/lib/base';
export const Login = () => {
  const handleLogin = () => {
    let url = 'https://accounts.google.com/o/oauth2/auth';
    url += '?client_id=';
    url += clientId;
    url += '&redirect_uri=';
    url += frontendUrl;
    url += '/app/callback-zustand';
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
