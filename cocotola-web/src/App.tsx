import { ChakraProvider } from '@chakra-ui/react';
import { BrowserRouter } from 'react-router-dom';

import { AppRoutes } from '@/routes';
import '@/App.css';

function App() {
  return (
    <BrowserRouter>
      <ChakraProvider>
        <AppRoutes />
        {/* <Route path="/" element={<Layout />}>
          <Route path="/app/login" element={<Login />} />
          <Route path="/app/callback" element={<Callback />} />
          <Route path="/test" element={<PrivateRoute element={<Test />} />} />
        </Route> */}
        {/* <Routes>
          <Route path="/" element={<Landing />} />
          <Route path="/app/login" element={<Login />} />
          <Route path="/app/callback" element={<Callback />} />
          <Route path="/test" element={<PrivateRoute element={<Test />} />} />
        </Routes> */}
      </ChakraProvider>
    </BrowserRouter>
  );
}

export default App;
