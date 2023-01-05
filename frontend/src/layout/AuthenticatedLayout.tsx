import { useState, ReactNode } from 'react';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';

import '@fontsource/roboto/300.css';

import AppBar from 'layout/AppBar';
import ErrorBoundary from 'layout/ErrorBoundary';
import Menu from 'layout/Menu';
import theme from 'layout/theme';
import { AuthContextProvider } from 'context/auth';

const drawerWidth = 240;

const Copyright = (props: any) => (
    <Typography variant="body2" color="text.secondary" align="center" {...props}>
        {'Copyright © '}
        Authz {new Date().getFullYear()}.
    </Typography>
);

type LayoutProps = {
    children?: ReactNode
    title?: string
}

export default function AuthenticatedLayout(props: LayoutProps) {
  const { children, title } = props;
  const [menuOpened, setMenuOpened] = useState(true);

  return (
    <ThemeProvider theme={theme}>
      <AuthContextProvider>
        <Box sx={{ display: 'flex' }}>
          <CssBaseline />

          <AppBar
            drawerWidth={drawerWidth}
            menuOpened={menuOpened}
            setMenuOpened={setMenuOpened}
            title={title}
          />

          <Menu
              drawerWidth={drawerWidth}
              menuOpened={menuOpened}
              setMenuOpened={setMenuOpened}
          />

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
            <ErrorBoundary>
              <Container maxWidth="xl" sx={{ mt: 2, mb: 2 }}>
                {children}
                <Copyright sx={{ pt: 4 }} />
              </Container>
            </ErrorBoundary>
          </Box>
        </Box>
      </AuthContextProvider>
    </ThemeProvider>
  );
}