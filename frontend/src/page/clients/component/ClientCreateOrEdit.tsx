import { useContext, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { Button, Paper, TextField, Typography } from '@mui/material';
import Grid from '@mui/material/Grid';
import SaveIcon from '@mui/icons-material/Save';
import useClientForm, { mapClientToFormData } from 'form/client';
import { AuthContext } from 'context/auth';
import { getClient } from 'service/model/client';
import { isAPIError } from 'service/error/model';
import { Client } from 'service/model/model';
import { useToast } from 'context/toast';

export default function ClientCreateOrEdit() {
  const toast = useToast();
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);
  const [client, setClient] = useState<Client>();

  const {
    onSubmit,
    register,
    defaultValues,
    errors,
    isSubmitting,
    reset,
  } = useClientForm(navigate, client);

  useEffect(() => {
    if (user === undefined) {
      return;
    }

    if (id === undefined) {
      return;
    }

    const fetch = async () => {
      const response = await getClient(user?.token!, id!);

      if (isAPIError(response)) {
        toast.error(`Unable to retrieve role: ${response.message}`);
      } else {
        setClient(response);
        reset(mapClientToFormData(response));  
      }
    };

    fetch();
  // eslint-disable-next-line
  }, [user, id, reset]);

  return (
    <div>
      <Typography
        variant="h3"
        gutterBottom
        sx={{
          marginTop: 1,
          marginBottom: 2
        }}>
        {client ? `Service account "${client?.name}"` : `Create new service account`}
      </Typography>

      <form onSubmit={onSubmit(user?.token!)}>
        <Paper sx={{ display: 'flex', flexDirection: 'column', p: 2 }}>
          <Typography variant='h5' sx={{ pb: 3 }}>
              Informations
          </Typography>

          <Grid container >
            <Grid
              size={{
                xs: 12,
                md: 6
              }}>
              <TextField {...register('name')}
                label='name'
                defaultValue={defaultValues?.name}
                error={errors?.name ? true : false}
                helperText={errors?.name?.message}
                sx={{ mb: 2, width: '100%' }}
                slotProps={{
                  inputLabel: { shrink: defaultValues?.name !== '' ? true : undefined }
                }}
              />

              {client ? (
                <TextField
                  disabled
                  label='client_id'
                  value={client?.client_id}
                  sx={{ mb: 2, width: '100%' }}
                  slotProps={{
                    inputLabel: { shrink: client?.client_id !== '' ? true : undefined }
                  }}
                />
              ) : null}

              {client ? (
                <TextField
                  disabled
                  label='client_secret'
                  value={client?.client_secret}
                  sx={{ mb: 2, width: '100%' }}
                  slotProps={{
                    inputLabel: { shrink: client?.client_secret !== '' ? true : undefined }
                  }}
                />
              ) : null}
            </Grid>
          </Grid>
        </Paper>

        <Grid
          size={{
            xs: 12,
            md: 8
          }}>
          <Button
            type='submit'
            loadingPosition='start'
            variant='contained'
            loading={isSubmitting}
            startIcon={<SaveIcon />}
            sx={{ marginTop: 2 }}
          >
            Submit
          </Button>

          <Button
            type='button'
            variant='text'
            onClick={() => navigate('/clients')}
            sx={{ marginLeft: 2, marginTop: 2 }}
          >
            Cancel
          </Button>
        </Grid>
      </form>
    </div>
  );
}