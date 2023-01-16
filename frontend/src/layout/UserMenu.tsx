import { MouseEvent, useContext, useState } from "react";
import Avatar from '@mui/material/Avatar';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import Divider from '@mui/material/Divider';
import Logout from '@mui/icons-material/Logout';
import { IconButton, ListItemText, Typography } from "@mui/material";
import { stringAvatar } from "./helper/avatar";
import { AuthContext } from "../context/auth";

export default function UserMenu() {
    const { user, logout } = useContext(AuthContext);

    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const anchorIsOpened = Boolean(anchorEl);

    const handleClickAvatar = (event: MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => setAnchorEl(null);

    return user ? (
        <>
            <IconButton
                onClick={handleClickAvatar}
                size='small'
                sx={{ ml: 2 }}
                role='user-menu'
                aria-controls={anchorIsOpened ? 'account-menu' : undefined}
                aria-haspopup="true"
                aria-expanded={anchorIsOpened ? 'true' : undefined}
            >
                <Avatar {...stringAvatar(`${user.username}`)} />
            </IconButton>

            <Menu
                anchorEl={anchorEl}
                id="account-menu"
                open={anchorIsOpened}
                onClose={handleClose}
                onClick={handleClose}
                PaperProps={{
                elevation: 0,
                sx: {
                    overflow: 'visible',
                    filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.32))',
                    mt: 1.5,
                    '& .MuiAvatar-root': {
                        width: 32,
                        height: 32,
                        ml: -0.5,
                        mr: 1,
                    },
                    '&:before': {
                        content: '""',
                        display: 'block',
                        position: 'absolute',
                        top: 0,
                        right: 19,
                        width: 10,
                        height: 10,
                        bgcolor: 'background.paper',
                        transform: 'translateY(-50%) rotate(45deg)',
                        zIndex: 0,
                    },
                },
                }}
                transformOrigin={{ horizontal: 'right', vertical: 'top' }}
                anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
            >
                <MenuItem sx={{ height: 40 }}>
                    <Typography variant='overline' display='block' gutterBottom>
                        {user.username}
                    </Typography>
                </MenuItem>

                <Divider />

                <MenuItem role='user-logout' onClick={() => logout()}>
                    <ListItemIcon>
                        <Logout fontSize='small' />
                    </ListItemIcon>
                    <ListItemText>
                        Logout
                    </ListItemText>
                </MenuItem>
            </Menu>
        </>
    ) : null
};