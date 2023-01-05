import { Typography } from "@mui/material";

export default function Copyright(props: any) {
    return (
        <Typography variant="body2" color="text.secondary" align="center" {...props}>
            {'Copyright © '}
            Authz {new Date().getFullYear()}.
        </Typography>
    )
}