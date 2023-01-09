import { useContext, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import LoadingButton from '@mui/lab/LoadingButton';
import { Alert, Button, Paper, TextField, Typography } from '@mui/material';
import Grid from '@mui/material/Grid';
import SaveIcon from '@mui/icons-material/Save';
import useUserForm, { mapUserToFormData } from 'form/user';
import { AuthContext } from 'context/auth';
import { getUser } from 'service/model/user';
import { isAPIError } from 'service/error/model';
import { User } from 'service/model/model';
import { useToast } from 'context/toast';
import CopyToClipboardButton from 'component/CopyToClipboardButton';

export default function UserCreateOrEdit() {
  const toast = useToast();
  const { id } = useParams();
  const navigate = useNavigate();
  const { user: currentUser } = useContext(AuthContext);
  const [user, setUser] = useState<User>();
  const [password, setPassword] = useState<string>();

  const {
    onSubmit,
    register,
    defaultValues,
    errors,
    isSubmitting,
    reset,
  } = useUserForm(navigate, user);

  useEffect(() => {
    if (currentUser === undefined) {
      return;
    }

    if (id === undefined) {
      return;
    }

    const fetch = async () => {
      const response = await getUser(currentUser?.token!, id!);

      if (isAPIError(response)) {
        toast.error(`Unable to retrieve role: ${response.message}`);
      } else {
        setUser(response);
        reset(mapUserToFormData(response));  
      }
    };

    fetch();
  // eslint-disable-next-line
  }, [currentUser, id, reset]);

  const onUserCreated = (user: User) => {
    setPassword(user?.password);
  }

  return (
    <div>
      <Typography variant="h3" gutterBottom marginTop={1} marginBottom={2}>
        {user ? `User "${user?.username}"` : `Create new user`}
      </Typography>

      {password ? (
        <Alert severity="success">
          User created with password: <CopyToClipboardButton color='default' text={password} />
        </Alert>
      ) : null}

      <form onSubmit={onSubmit(currentUser?.token!, onUserCreated)}>
        <Paper sx={{ display: 'flex', flexDirection: 'column', p: 2 }}>
          <Typography variant='h5' sx={{ pb: 3 }}>
              Informations
          </Typography>

          <Grid container >
            <Grid item xs={12} md={6}>
              <TextField {...register('username')}
                label='Username'
                defaultValue={defaultValues?.username}
                error={errors?.username ? true : false}
                helperText={errors?.username?.message}
                InputLabelProps={{ shrink: defaultValues?.username !== '' ? true : undefined }}
                sx={{ mb: 2, width: '100%' }}
              />
            </Grid>
          </Grid>
        </Paper>

        <Grid item xs={12} md={8}>
          <LoadingButton
            type='submit'
            loadingPosition='start'
            variant='contained'
            loading={isSubmitting}
            startIcon={<SaveIcon />}
            sx={{ marginTop: 2 }}
          >
            Enregistrer
          </LoadingButton>

          <Button
            type='button'
            variant='text'
            onClick={() => navigate('/users')}
            sx={{ marginLeft: 2, marginTop: 2 }}
          >
            Retour
          </Button>
        </Grid>
      </form>
    </div>
  );
}