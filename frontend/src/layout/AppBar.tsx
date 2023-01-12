import { styled } from '@mui/material/styles';
import IconButton from '@mui/material/IconButton';
import MenuIcon from '@mui/icons-material/Menu';
import MuiAppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import { Link, MenuItem, Theme, useMediaQuery } from '@mui/material';
import UserMenu from 'layout/UserMenu';

type AppBarProps = {
    drawerWidth: number
    menuOpened: boolean
    setMenuOpened: Function
    title?: string
};

export default function AppBar(props: AppBarProps) {
    const { drawerWidth, menuOpened, setMenuOpened, title } = props;
    const matchesSmall = useMediaQuery((theme: Theme) => theme.breakpoints.up('sm'));

    const AppBar = styled(MuiAppBar, {
        shouldForwardProp: (prop) => prop !== 'open',
    })(({ theme }) => ({
        zIndex: theme.zIndex.drawer + 1,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
        ...(menuOpened && {
            marginLeft: drawerWidth,
            width: `calc(100% - ${drawerWidth}px)`,
            transition: theme.transitions.create(['width', 'margin'], {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.enteringScreen,
            }),
        }),
    }));

    return (
        <AppBar position="absolute">
            <Toolbar
            sx={{
                pr: '24px',
            }}
            >
                <IconButton
                    edge="start"
                    color="inherit"
                    aria-label="open drawer"
                    onClick={() => setMenuOpened(!menuOpened)}
                    sx={{
                        marginLeft: matchesSmall ? '-10px' : '-7px',
                        marginRight: '36px',
                        ...(menuOpened && { display: 'none' }),
                    }}
                >
                    <MenuIcon />
                </IconButton>
                <Typography
                    component="h1"
                    variant="h6"
                    color="inherit"
                    noWrap
                    sx={{ flexGrow: 1 }}
                >
                    {title || 'Authz'}
                </Typography>

                <MenuItem component={Link} href='https://github.com/eko/authz' target='_blank'>
                  GitHub
                </MenuItem>

                <UserMenu />
            </Toolbar>
        </AppBar>
    )
};