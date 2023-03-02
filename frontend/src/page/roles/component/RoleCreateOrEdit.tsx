import { useContext, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import moment from "moment";
import LoadingButton from '@mui/lab/LoadingButton';
import { Button, Divider, Paper, TextField, Tooltip, Typography } from '@mui/material';
import Grid from '@mui/material/Grid';
import SaveIcon from '@mui/icons-material/Save';
import useRoleForm, { mapRoleToFormData } from 'form/role';
import MultipleAutocompleteInput, { ItemType } from 'component/MultipleAutocompleteInput';
import { AuthContext } from 'context/auth';
import { getRole } from 'service/model/role';
import { Role } from 'service/model/model';
import { isAPIError } from 'service/error/model';
import { useToast } from 'context/toast';
import { getPolicies } from 'service/model/policy';

export default function RoleCreateOrEdit() {
  const toast = useToast();
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);
  const [role, setRole] = useState<Role>();

  const {
    onSubmit,
    register,
    setValue,
    defaultValues,
    errors,
    isSubmitting,
    reset,
  } = useRoleForm(navigate, role);

  useEffect(() => {
    if (user === undefined) {
      return;
    }

    if (id === undefined) {
      return;
    }

    const fetch = async () => {
      const response = await getRole(user?.token!, id!);

      if (isAPIError(response)) {
        toast.error(`Unable to retrieve role: ${response.message}`);
      } else {
        setRole(response);
        reset(mapRoleToFormData(response));  
      }
    };

    fetch();
  // eslint-disable-next-line
  }, [user, id, reset]);

  const policiesFetcher = async (input: string): Promise<ItemType[]> => {
    const response = await getPolicies(user?.token!, 1, 50, {
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

  const haveUpdatedAt = !role?.updated_at.toString().startsWith('0001-01-01');

  return (
    <div>
      <Typography variant="h3" gutterBottom marginTop={1} marginBottom={2}>
          {role?.id ? `Role "${role?.id}"` : `Create new role`}
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

              {id ? (
                <>
                  <Tooltip title='This field is locked' placement='right'>
                    <TextField
                      disabled
                      label="Creation date"
                      value={moment(role?.created_at).format('LLL')}
                      InputLabelProps={{ shrink: true }}
                      sx={{ mr: 2, mb: 2, width: '200px' }}
                    />
                  </Tooltip>

                  {haveUpdatedAt ? (
                    <Tooltip title='This field is locked' placement='right'>
                      <TextField
                        disabled
                        label="Update date"
                        value={moment(role?.updated_at).format('LLL')}
                        InputLabelProps={{ shrink: true }}
                        sx={{ width: '200px' }}
                      />
                    </Tooltip>
                ) : null}
                </>
              ) : null}
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
                label='Associated policies'
                placeholder='Search for a policy...'
                defaultValues={defaultValues?.policies as any}
                errorField={errors?.policies}
                fetcher={policiesFetcher}
                setValue={(items: ItemType[]) => setValue('policies', items)}
                style={{ marginBottom: 2, marginTop: 2 }}
                inputSx={{ width: '100%' }}
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
            onClick={() => navigate('/roles')}
            sx={{ marginLeft: 2, marginTop: 2 }}
          >
            Retour
          </Button>
        </Grid>
      </form>
    </div>
  );
}