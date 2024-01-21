import { ChakraProvider } from '@chakra-ui/react';
import { BrowserRouter } from 'react-router-dom';

import { AppRoutes } from '@/routes';
import '@/App.css';

function App() {
  return (
    <BrowserRouter>
      <ChakraProvider>
        <AppRoutes />
      </ChakraProvider>
    </BrowserRouter>
  );
}

export default App;
