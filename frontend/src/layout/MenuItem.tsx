import { useMatch, useNavigate } from "react-router-dom";
import { Collapse, ListItemButton, ListItemIcon, ListItemText, Theme, Tooltip } from '@mui/material';
import useMediaQuery from '@mui/material/useMediaQuery';
import { ExpandLess, ExpandMore } from "@mui/icons-material";

type MenuItemProps = {
    label: string
    path?: string
    handleClick?: Function
    icon?: React.ReactElement
    menuOpened: boolean
    nested?: React.ReactElement
    nestedOpened?: boolean
};

export default function MenuItem(props: MenuItemProps) {
    const { label, path, handleClick, icon, menuOpened, nested, nestedOpened } = props;
    const navigate = useNavigate();
    const matchesSmall = useMediaQuery((theme: Theme) => theme.breakpoints.up('sm'));

    return (
        <>
            <Tooltip title={!menuOpened ? label : ''} placement='right' arrow>
                <ListItemButton
                    selected={useMatch(path === '/' ? path : path + '/*' || '') && !nested ? true : false}
                    onClick={() => handleClick ? handleClick() : navigate(path || '')}
                >
                    {icon ? (<ListItemIcon sx={{
                        marginLeft: matchesSmall && !menuOpened ? '6px' : 0,
                    }}>
                        {icon}
                    </ListItemIcon>): null}
                    <ListItemText primary={label} />

                    {nested ? (
                        nestedOpened ? <ExpandLess /> : <ExpandMore />
                    ) : null}
                </ListItemButton>
            </Tooltip>

            {nested ? (
                <Collapse in={nestedOpened} timeout="auto" unmountOnExit>
                    {nested}
                </Collapse>
            ) : null}
        </>
    );
}