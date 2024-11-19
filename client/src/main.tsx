import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';

import { AppRoutes } from './routes';
import { ThemeProvider } from './lib/theme-provider';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider defaultTheme='dark' storageKey='vite-ui-theme'>
      <AppRoutes />
    </ThemeProvider>
  </StrictMode>
);
