import { useContext, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import LoadingButton from '@mui/lab/LoadingButton';
import { Button, Divider, Paper, TextField, Typography } from '@mui/material';
import Grid from '@mui/material/Grid';
import SaveIcon from '@mui/icons-material/Save';
import usePrincipalForm, { mapPrincipalToFormData } from 'form/principal';
import { AuthContext } from 'context/auth';
import Attributes from './Attributes';
import { getPrincipal } from 'service/model/principal';
import { isAPIError } from 'service/error/model';
import { Principal } from 'service/model/model';
import { useToast } from 'context/toast';
import { getRoles } from 'service/model/role';
import MultipleAutocompleteInput, { ItemType } from 'component/MultipleAutocompleteInput';

export default function PrincipalCreateOrEdit() {
  const toast = useToast();
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);
  const [principal, setPrincipal] = useState<Principal>();

  const {
    control,
    onSubmit,
    register,
    defaultValues,
    errors,
    isSubmitting,
    reset,
    setValue,
  } = usePrincipalForm(navigate, principal);

  useEffect(() => {
    if (user === undefined) {
      return;
    }

    if (id === undefined) {
      return;
    }

    const fetch = async () => {
      const response = await getPrincipal(user?.token!, id!);

      if (isAPIError(response)) {
        toast.error(`Unable to retrieve role: ${response.message}`);
      } else {
        setPrincipal(response);
        reset(mapPrincipalToFormData(response));  
      }
    };

    fetch();
  // eslint-disable-next-line
  }, [user, id, reset]);

  const rolesFetcher = async (input: string): Promise<ItemType[]> => {
    const response = await getRoles(user?.token!, 1, 50, {
      field: 'id',
      operator: 'contains',
      value: input,
    }, {
      field: 'id',
      order: 'asc',
    });

    if (isAPIError(response)) {
      return Promise.resolve([]);
    }

    return response.data.map(item => ({
      id: item.id,
      label: item.id,
    }));
  }

  return (
    <div>
      <Typography variant="h3" gutterBottom marginTop={1} marginBottom={2}>
        {principal?.id ? `Principal "${principal?.id}"` : `Create new principal`}
      </Typography>

      <form onSubmit={onSubmit(user?.token!)}>
        <Paper sx={{ display: 'flex', flexDirection: 'column', p: 2 }}>
          <Typography variant='h5' sx={{ pb: 3 }}>
              Informations
          </Typography>

          <Grid container >
            <Grid item xs={12} md={6}>
              <TextField {...register('id')}
                label='Name'
                defaultValue={defaultValues?.id}
                error={errors?.id ? true : false}
                helperText={errors?.id?.message}
                InputLabelProps={{ shrink: defaultValues?.id !== '' ? true : undefined }}
                sx={{ mb: 2, width: '100%' }}
              />
            </Grid>
          </Grid>
        </Paper>

        <Divider style={{ marginTop: 20, marginBottom: 20 }} />

        <Paper sx={{ display: 'flex', flexDirection: 'column', p: 2 }}>
          <Typography variant='h5' sx={{ pb: 3 }}>
              Associations
          </Typography>

          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <MultipleAutocompleteInput
                label='Associated roles'
                placeholder='Search for a role...'
                defaultValues={defaultValues?.roles as any}
                errorField={errors?.roles}
                fetcher={rolesFetcher}
                setValue={(items: ItemType[]) => setValue('roles', items)}
                style={{ marginBottom: 2, marginTop: 2 }}
                inputSx={{ width: '100%' }}
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
            <Grid item xs={12} md={6}>
              <Attributes
                control={control}
                defaultValues={defaultValues}
                fieldName='attributes'
                register={register}
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
            Submit
          </LoadingButton>

          <Button
            type='button'
            variant='text'
            onClick={() => navigate('/principals')}
            sx={{ marginLeft: 2, marginTop: 2 }}
          >
            Cancel
          </Button>
        </Grid>
      </form>
    </div>
  );
}