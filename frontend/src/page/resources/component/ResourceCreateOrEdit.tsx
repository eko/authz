import { useContext, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { Button, Divider, Paper, TextField, Typography } from '@mui/material';
import Grid from '@mui/material/Grid';
import SaveIcon from '@mui/icons-material/Save';
import useResourceForm, { mapResourceToFormData } from 'form/resource';
import { AuthContext } from 'context/auth';
import Attributes from './Attributes';
import { getResource } from 'service/model/resource';
import { isAPIError } from 'service/error/model';
import { Resource } from 'service/model/model';
import { useToast } from 'context/toast';

export default function ResourceCreateOrEdit() {
  const toast = useToast();
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);
  const [resource, setResource] = useState<Resource>();

  const {
    control,
    onSubmit,
    register,
    defaultValues,
    errors,
    isSubmitting,
    reset,
  } = useResourceForm(navigate, resource);

  useEffect(() => {
    if (user === undefined) {
      return;
    }

    if (id === undefined) {
      return;
    }

    const fetch = async () => {
      const response = await getResource(user?.token!, id!);

      if (isAPIError(response)) {
        toast.error(`Unable to retrieve role: ${response.message}`);
      } else {
        setResource(response);
        reset(mapResourceToFormData(response));  
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
        {resource?.id ? `Resource "${resource?.id}"` : `Create new resource`}
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
              <TextField {...register('id')}
                label='Name'
                defaultValue={defaultValues?.id}
                error={errors?.id ? true : false}
                helperText={errors?.id?.message}
                sx={{ mb: 2, width: '100%' }}
                slotProps={{
                  inputLabel: { shrink: defaultValues?.id !== '' ? true : undefined }
                }}
              />

              <TextField {...register('kind')}
                label='Kind'
                defaultValue={defaultValues?.kind}
                error={errors?.kind ? true : false}
                helperText={errors?.kind?.message}
                sx={{ mb: 2, width: '100%' }}
                slotProps={{
                  inputLabel: { shrink: defaultValues?.kind !== '' ? true : undefined }
                }}
              />

              <TextField {...register('value')}
                label='Identifier value'
                defaultValue={defaultValues?.value}
                error={errors?.value ? true : false}
                helperText={errors?.value?.message}
                sx={{ mb: 2, width: '100%' }}
                slotProps={{
                  inputLabel: { shrink: defaultValues?.value !== '' ? true : undefined }
                }}
              />
            </Grid>
          </Grid>
        </Paper>

        <Divider style={{ marginTop: 20, marginBottom: 20 }} />

        <Paper sx={{ display: 'flex', flexDirection: 'column', p: 2 }}>
          <Typography variant='h5' sx={{ pb: 3 }}>
            Attributes
          </Typography>

          <Grid container spacing={2}>
            <Grid
              size={{
                xs: 12,
                md: 6
              }}>
              <Attributes
                control={control}
                defaultValues={defaultValues}
                fieldName='attributes'
                register={register}
              />
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
            onClick={() => navigate('/resources')}
            sx={{ marginLeft: 2, marginTop: 2 }}
          >
            Cancel
          </Button>
        </Grid>
      </form>
    </div>
  );
}