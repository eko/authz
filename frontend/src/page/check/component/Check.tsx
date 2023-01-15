import { useContext, useState } from 'react';
import { useNavigate } from 'react-router';
import LoadingButton from '@mui/lab/LoadingButton';
import { Alert, Paper, Typography } from '@mui/material';
import Grid from '@mui/material/Grid';
import SaveIcon from '@mui/icons-material/Save';
import { ItemType } from 'component/MultipleAutocompleteInput';
import { AuthContext } from 'context/auth';
import { isAPIError } from 'service/error/model';
import { getResources } from 'service/model/resource';
import { getActions } from 'service/model/action';
import useCheckForm from 'form/check';
import SingleAutocompleteInput from 'component/SingleAutocompleteInput';
import { getPrincipals } from 'service/model/principal';
import { Check as CheckItem, CheckResponse } from 'service/model/check';

export default function Check() {
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);
  const [check, setCheck] = useState<CheckItem>();

  const {
    onSubmit,
    setValue,
    defaultValues,
    errors,
    isSubmitting,
  } = useCheckForm(navigate);

  const principalsFetcher = async (input: string): Promise<ItemType[]> => {
    const response = await getPrincipals(user?.token!, 1, 50, {
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

  const resourcesFetcher = async (input: string): Promise<ItemType[]> => {
    const response = await getResources(user?.token!, 1, 50, {
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
      raw: item,
    }));
  }

  const actionsFetcher = async (input: string): Promise<ItemType[]> => {
    const response = await getActions(user?.token!, 1, 50, {
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

  const onSubmitSuccess = (response: CheckResponse) => {
    setCheck(response.checks[0]);
  }

  return (
    <div>
      <form onSubmit={onSubmit(user?.token!, onSubmitSuccess)}>
        <Grid container spacing={2}>
          <Grid item xs={12} lg={4}>
            <Paper sx={{ display: 'flex', flexDirection: 'column', p: 2 }}>
              <Typography variant='h5' sx={{ pb: 3 }}>
                  Who?
              </Typography>

              <Grid container spacing={2}>
                <Grid item xs={12} md={12}>
                  <SingleAutocompleteInput
                    label='Principal'
                    placeholder='Search for a principal...'
                    defaultValue={defaultValues?.principal as any}
                    errorField={errors?.principal}
                    fetcher={principalsFetcher}
                    setValue={(item: ItemType) => setValue('principal', item)}
                    style={{ marginBottom: 2, marginTop: 2 }}
                    inputSx={{ width: '100%' }}
                  />
                </Grid>
              </Grid>
            </Paper>
          </Grid>

          <Grid item xs={12} lg={4}>
            <Paper sx={{ display: 'flex', flexDirection: 'column', p: 2 }}>
              <Typography variant='h5' sx={{ pb: 3 }}>
                  Which resource?
              </Typography>

              <Grid container spacing={2}>
                <Grid item xs={12} md={12}>
                  <SingleAutocompleteInput
                    allowAdd
                    label='Resource'
                    placeholder='Search for a resource...'
                    defaultValue={defaultValues?.resource as any}
                    errorField={errors?.resource}
                    fetcher={resourcesFetcher}
                    setValue={(item: ItemType) => setValue('resource', item)}
                    style={{ marginBottom: 2, marginTop: 2 }}
                    inputSx={{ width: '100%' }}
                  />
                </Grid>
              </Grid>
            </Paper>
          </Grid>

          <Grid item xs={12} lg={4}>
            <Paper sx={{ display: 'flex', flexDirection: 'column', p: 2 }}>
              <Typography variant='h5' sx={{ pb: 3 }}>
                  Which action?
              </Typography>

              <Grid container spacing={2}>
                <Grid item xs={12} md={12}>
                  <SingleAutocompleteInput
                      label='Action'
                      placeholder='Search for an action...'
                      defaultValue={defaultValues?.action as any}
                      errorField={errors?.action}
                      fetcher={actionsFetcher}
                      setValue={(item: ItemType) => setValue('action', item)}
                      style={{ marginBottom: 2, marginTop: 2 }}
                      inputSx={{ width: '100%' }}
                  />
                </Grid>
              </Grid>
            </Paper>
          </Grid>
        </Grid>

        {check ? (
          <div style={{ marginTop: '20px' }}>
            {check.is_allowed ? (
              <Alert variant='filled' severity='success' sx={{ color: '#ffffff' }}>
                Access is allowed
              </Alert>
            ) : (
              <Alert variant='filled' severity='error'>
                Access is denied
              </Alert>
            )}
          </div>
        ) : null}

        <Grid item xs={12} md={8}>
          <LoadingButton
            type='submit'
            loadingPosition='start'
            variant='contained'
            loading={isSubmitting}
            startIcon={<SaveIcon />}
            sx={{ marginTop: 2 }}
          >
            Check
          </LoadingButton>
        </Grid>
      </form>
    </div>
  );
}