import { ReactNode } from 'react';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import Container from '@mui/material/Container';

import '@fontsource/roboto/300.css';

import Copyright from 'layout/Copyright';
import theme from 'layout/theme';

type LayoutProps = {
    children?: ReactNode;
}

export default function AnonymousLayout(props: LayoutProps) {
  const { children } = props;

  return (
    <ThemeProvider theme={theme}>
      <Box sx={{ display: 'flex' }}>
          <CssBaseline />

          <Box
          component="main"
          sx={{
              backgroundColor: (theme) =>
              theme.palette.mode === 'light'
                  ? theme.palette.grey[100]
                  : theme.palette.grey[900],
              flexGrow: 1,
              height: '100vh',
              overflow: 'auto',
          }}
          >
          <Toolbar />
          <Container maxWidth="lg" sx={{ mt: 2, mb: 2 }}>
              {children}
              <Copyright sx={{ pt: 4 }} />
          </Container>
          </Box>
      </Box>
    </ThemeProvider>
  );
}