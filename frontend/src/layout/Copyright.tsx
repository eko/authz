import { Typography } from "@mui/material";

export default function Copyright(props: any) {
    return (
        <Typography
            variant="body2"
            align="center"
            {...props}
            sx={[{
                color: "text.secondary"
            }, ...(Array.isArray(props.sx) ? props.sx : [props.sx])]}>
            {'Copyright © '}Authz {new Date().getFullYear()}.
                    </Typography>
    );
}