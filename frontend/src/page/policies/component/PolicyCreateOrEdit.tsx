import { useContext, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import dayjs from 'helper/date';
import { Button, Divider, Paper, TextField, Tooltip, Typography } from '@mui/material';
import Grid from '@mui/material/Grid';
import SaveIcon from '@mui/icons-material/Save';
import usePolicyForm, { mapPolicyToFormData } from 'form/policy';
import MultipleAutocompleteInput, { ItemType } from 'component/MultipleAutocompleteInput';
import { AuthContext } from 'context/auth';
import { Policy } from 'service/model/model';
import { isAPIError } from 'service/error/model';
import { useToast } from 'context/toast';
import { getPolicy } from 'service/model/policy';
import { getResources } from 'service/model/resource';
import { getActions } from 'service/model/action';
import AttributeRules from './AttributeRules';

export default function PolicyCreateOrEdit() {
  const toast = useToast();
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);
  const [policy, setPolicy] = useState<Policy>();

  const {
    control,
    onSubmit,
    register,
    setValue,
    defaultValues,
    errors,
    isSubmitting,
    reset,
  } = usePolicyForm(navigate, policy);

  useEffect(() => {
    if (user === undefined) {
      return;
    }

    if (id === undefined) {
      return;
    }

    const fetch = async () => {
      const response = await getPolicy(user?.token!, id!);

      if (isAPIError(response)) {
        toast.error(`Unable to retrieve policy: ${response.message}`);
      } else {
        setPolicy(response);
        reset(mapPolicyToFormData(response));  
      }
    };

    fetch();
  // eslint-disable-next-line
  }, [user, id, reset]);

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

  const haveUpdatedAt = !policy?.updated_at.toString().startsWith('0001-01-01');

  return (
    <div>
      <Typography
        variant="h3"
        gutterBottom
        sx={{
          marginTop: 1,
          marginBottom: 2
        }}>
          {policy?.id ? `Policy "${policy?.id}"` : `Create new policy`}
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

              {id ? (
                <>
                  <Tooltip title='This field is locked' placement='right'>
                    <TextField
                      disabled
                      label="Creation date"
                      value={dayjs(policy?.created_at).format('LLL')}
                      sx={{ mr: 2, mb: 2, width: '200px' }}
                      slotProps={{
                        inputLabel: { shrink: true }
                      }}
                    />
                  </Tooltip>

                  {haveUpdatedAt ? (
                    <Tooltip title='This field is locked' placement='right'>
                      <TextField
                        disabled
                        label="Update date"
                        value={dayjs(policy?.updated_at).format('LLL')}
                        sx={{ width: '200px' }}
                        slotProps={{
                          inputLabel: { shrink: true }
                        }}
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
            <Grid
              size={{
                xs: 12,
                md: 6
              }}>
              <MultipleAutocompleteInput
                allowAdd
                label='Associated resources'
                placeholder='Search for a resource...'
                defaultValues={defaultValues?.resources as any}
                errorField={errors?.resources}
                fetcher={resourcesFetcher}
                setValue={(items: ItemType[]) => setValue('resources', items)}
                style={{ marginBottom: 2, marginTop: 2 }}
                inputSx={{ width: '100%' }}
              />
            </Grid>

            <Grid
              size={{
                xs: 12,
                md: 6
              }}>
              <MultipleAutocompleteInput
                label='Associated actions'
                placeholder='Search for an action...'
                defaultValues={defaultValues?.actions as any}
                errorField={errors?.actions}
                fetcher={actionsFetcher}
                setValue={(items: ItemType[]) => setValue('actions', items)}
                style={{ marginBottom: 2, marginTop: 2 }}
                inputSx={{ width: '100%' }}
              />
            </Grid>
          </Grid>
        </Paper>

        <Divider style={{ marginTop: 20, marginBottom: 20 }} />

        <Paper sx={{ display: 'flex', flexDirection: 'column', p: 2 }}>
          <Typography variant='h5' sx={{ pb: 3 }}>
              Attribute rules
          </Typography>

          <Grid container spacing={2}>
            <Grid
              size={{
                xs: 12,
                md: 6
              }}>
              <AttributeRules
                fieldName='attribute_rules'
                register={register}
                control={control}
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
            onClick={() => navigate('/policies')}
            sx={{ marginLeft: 2, marginTop: 2 }}
          >
            Cancel
          </Button>
        </Grid>
      </form>
    </div>
  );
}