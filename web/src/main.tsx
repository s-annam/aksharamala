import React from 'react';
import ReactDOM from 'react-dom/client';
import { CustomThemeProvider } from './ThemeProvider';
import App from './App';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <CustomThemeProvider>
      <App />
    </CustomThemeProvider>
  </React.StrictMode>
);
