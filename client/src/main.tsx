import { createRoot } from 'react-dom/client';
import './index.css';

import { AppRoutes } from './routes';
import { ThemeProvider } from './lib/theme-provider';
import { AuthCtx } from './lib/authCTX';

createRoot(document.getElementById('root')!).render(
  <>
    <ThemeProvider defaultTheme='dark' storageKey='vite-ui-theme'>
      <AuthCtx>
        <AppRoutes />
      </AuthCtx>
    </ThemeProvider>
  </>
);
