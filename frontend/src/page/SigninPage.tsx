import { Grid, InputAdornment, TextField, Typography } from '@mui/material';
import Paper from '@mui/material/Paper';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import LoginIcon from '@mui/icons-material/Login';
import PinIcon from '@mui/icons-material/Pin';
import { useNavigate } from 'react-router';

import AnonymousLayout from 'layout/AnonymousLayout';
import useSigninForm from 'form/signin';
import { LoadingButton } from '@mui/lab';

export default function SigninPage() {
  const navigate = useNavigate();

  const {
    onSubmit,
    register,
    defaultValues,
    errors,
    isSubmitting,
  } = useSigninForm(navigate);

  return (
    <AnonymousLayout>
      <form onSubmit={onSubmit()}>
        <Paper sx={{ display: 'flex', flexDirection: 'column', p: 6 }}>
          <Grid container spacing={4}>
            <Grid container item display='flex' direction='column' alignSelf='center' alignItems='center' xs={12} md={6}>
              <img width='300' src='/logo-full.png' alt='Authz' title='Authz' />
            </Grid>

            <Grid container item display='flex' direction='column' alignSelf='center' xs={12} md={6}>
              <Typography variant="h3" gutterBottom sx={{ textAlign: 'center', marginBottom: 4 }}>
                  Sign In
              </Typography>

              <TextField {...register('username')}
                label='Username'
                defaultValue={defaultValues?.username}
                error={errors?.username ? true : false}
                helperText={errors?.username?.message}
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="start">
                      <AccountCircleIcon />
                    </InputAdornment>
                  ),
                }}
                sx={{ mb: 2, width: '100%' }}
              />

              <TextField {...register('password')}
                type='password'
                label='Password'
                defaultValue={defaultValues?.password}
                error={errors?.password ? true : false}
                helperText={errors?.password?.message}
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="start">
                      <PinIcon />
                    </InputAdornment>
                  ),
                }}
                sx={{ mb: 2, width: '100%' }}
              />

              <LoadingButton
                type='submit'
                loadingPosition='start'
                variant='contained'
                loading={isSubmitting}
                startIcon={<LoginIcon />}
                sx={{ marginTop: 2 }}
              >
                Log in
              </LoadingButton>
            </Grid>
          </Grid>
        </Paper>
      </form>
   </AnonymousLayout>
  );
}