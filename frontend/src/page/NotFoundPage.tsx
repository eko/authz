import Paper from '@mui/material/Paper';
import { Typography } from '@mui/material';

import AuthenticatedLayout from 'layout/AuthenticatedLayout';

export default function NotFoundPage() {
  return (
    <AuthenticatedLayout title='We are unable to find what you are looking for'>
      <Paper
        sx={{
          p: 10,
          marginTop: 3,
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <div style={{ textAlign: 'center' }}>
          <Typography variant='body1' component='h1'>
            Page not found
          </Typography>
        </div>
      </Paper>
    </AuthenticatedLayout>
  );
}