import { Alert, AlertTitle, Button, Typography } from '@mui/material';
import React from 'react';

type ErrorBoundaryProps = {
    children: React.ReactNode
}

type ErrorBoundaryState = {
    showDetails: boolean
    hasError: boolean
    error?: Error
}

export default class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
    constructor(props: ErrorBoundaryProps) {
        super(props);
        this.state = {
            showDetails: false,
            hasError: false,
        };
    }
    
    static getDerivedStateFromError(error: Error) {
        return {
            hasError: true,
            error: error,
        };
    }
    
    componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
        console.error("Uncaught error:", error, errorInfo);

        this.setState({
            error: error,
            ...this.state
        });
    }
    
    render() {
        const { children } = this.props;
        const { hasError, error } = this.state;

        const reload = () => window.location.reload();

        return hasError ? (
            <>
                <Alert severity="error" action={
                    <Button color="inherit" size="small" onClick={reload}>
                        Reload
                    </Button>
                }>
                    <AlertTitle>An error occurred.</AlertTitle>
                    {error?.message}
                </Alert>

                <Typography textAlign='center' color='#ebebeb' sx={{
                    fontSize: '500px',
                    marginTop: '100px',
                }}>
                    Error
                </Typography>
            </>
        ) : children
    }
}