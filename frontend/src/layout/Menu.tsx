import { useEffect } from 'react';
import { styled } from '@mui/material/styles';
import MuiDrawer from '@mui/material/Drawer';
import Toolbar from '@mui/material/Toolbar';
import IconButton from '@mui/material/IconButton';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import { Divider, List, useMediaQuery, Theme } from '@mui/material';

import ClassIcon from '@mui/icons-material/Class';
import ContactsIcon from '@mui/icons-material/Contacts';
import GridViewIcon from '@mui/icons-material/GridView';
import GppGoodIcon from '@mui/icons-material/GppGood';
// import GroupIcon from '@mui/icons-material/Group';
import SecurityIcon from '@mui/icons-material/Security';
import SchoolIcon from '@mui/icons-material/School';
import PersonIcon from '@mui/icons-material/Person';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';

import MenuItem from 'layout/MenuItem';

type MenuProps = {
    drawerWidth: number
    menuOpened: boolean
    setMenuOpened: Function
};

export default function Menu(props: MenuProps) {
  const matchesSmallOrGreater = useMediaQuery((theme: Theme) => theme.breakpoints.up('sm'));
  const { drawerWidth, menuOpened, setMenuOpened } = props;

  useEffect(() => {
    if (menuOpened && !matchesSmallOrGreater) {
      setMenuOpened(false);
    } else if (!menuOpened && matchesSmallOrGreater) {
      setMenuOpened(true);
    }
  // eslint-disable-next-line
  }, []);

  const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
    ({ theme, open }) => ({
      '& .MuiDrawer-paper': {
        position: 'relative',
        whiteSpace: 'nowrap',
        width: drawerWidth,
        transition: theme.transitions.create('width', {
          easing: theme.transitions.easing.sharp,
          duration: theme.transitions.duration.enteringScreen,
        }),
        boxSizing: 'border-box',
        ...(!open && {
          overflowX: 'hidden',
          transition: theme.transitions.create('width', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
          }),
          width: theme.spacing(7),
          [theme.breakpoints.up('sm')]: {
            width: theme.spacing(9),
          },
        }),
      },
    }),
  );

  return (
    <Drawer variant="permanent" open={menuOpened}>
        <Toolbar
        sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'flex-end',
            px: [1],
        }}
        >
            <a href='/' title='Authz' style={{ marginRight: '50px', marginTop: '5px' }}>
              <img width='40' src='/logo.png' alt='Authz' title='Authz' style={{ borderRadius: '20px' }} />
            </a>

            <IconButton onClick={() => setMenuOpened(!menuOpened)}>
                <ChevronLeftIcon />
            </IconButton>
        </Toolbar>

        <Divider />

        <List component="nav">
          <MenuItem
            label='Dashboard'
            path='/'
            icon={<GridViewIcon />}
            menuOpened={menuOpened}
          />

          <Divider sx={{ margin: 2 }} />

          {/* <MenuItem
            label='Groups'
            path='/groups'
            icon={<GroupIcon />}
            menuOpened={menuOpened}
          /> */}

          <MenuItem
            label='Roles'
            path='/roles'
            icon={<SchoolIcon />}
            menuOpened={menuOpened}
          />

          <MenuItem
            label='Policies'
            path='/policies'
            icon={<SecurityIcon />}
            menuOpened={menuOpened}
          />

          <Divider sx={{ margin: 2 }} />

          <MenuItem
            label='Resources'
            path='/resources'
            icon={<ClassIcon />}
            menuOpened={menuOpened}
          />

          <MenuItem
            label='Principals'
            path='/principals'
            icon={<ContactsIcon />}
            menuOpened={menuOpened}
          />

          <Divider sx={{ margin: 2 }} />

          <MenuItem
            label='Users'
            path='/users'
            icon={<PersonIcon />}
            menuOpened={menuOpened}
          />

          <MenuItem
            label='Service accounts'
            path='/clients'
            icon={<AdminPanelSettingsIcon />}
            menuOpened={menuOpened}
          />

          <Divider sx={{ margin: 2 }} />

          <MenuItem
            label='Check access'
            path='/check'
            icon={<GppGoodIcon />}
            menuOpened={menuOpened}
          />
        </List>
    </Drawer>
 );
}