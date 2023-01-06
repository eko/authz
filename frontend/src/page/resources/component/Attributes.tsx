import { Alert, Button, Grid, IconButton, TextField } from '@mui/material';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import DeleteForeverOutlinedIcon from '@mui/icons-material/DeleteForeverOutlined';
import { Control, DeepPartial, useFieldArray, UseFormRegister } from 'react-hook-form';
import { ResourceFormData } from 'form/resource';

type ResourceFormDataKey = keyof ResourceFormData;

type AttributeProps = {
    control: Control<ResourceFormData>
    disabled?: boolean
    defaultValues?: Readonly<DeepPartial<ResourceFormData>> | ResourceFormData
    fieldName: ResourceFormDataKey
    register: UseFormRegister<ResourceFormData>
};

export default function Attributes(props: AttributeProps) {
  const { control, disabled, fieldName, register } = props;

  const { fields, append, remove } = useFieldArray({
    control,
    name: fieldName as any,
  });

  return (
    <>
      {!disabled ? (
        <Button
          variant='outlined'
          size='small'
          color='primary'
          startIcon={<AddCircleOutlineIcon />}
          onClick={() => append('')}
          sx={{ mb: 2 }}
        >
          Add new attribute
        </Button>
      ) : null}

      {fields.length === 0 ? (
        <Alert variant="standard" severity="info">
          No attribute rule defined. Your policy will be applied on all selected resources.
        </Alert>
      ) : null}

        {fields.map((item, index) => (
          <Grid key={item.id} container spacing={2} sx={{ mt: '1px' }}>
            <Grid item xs={5} md={5}>
              <TextField {...register(`${fieldName}.${index}.key` as any)}
                disabled={disabled}
                label={`Key #${index + 1}`}
                placeholder='key'
                sx={{ width: '100%' }}
              />
            </Grid>

            <Grid item xs={5} md={5}>
              <TextField {...register(`${fieldName}.${index}.value` as any)}
                disabled={disabled}
                label={`Value #${index + 1}`}
                placeholder='value'
                sx={{ width: '100%' }}
              />
            </Grid>

            <Grid item xs={2} md={2}>
              <IconButton
                type='button'
                size='small'
                title='Delete this rule'
                color='error'
                sx={{ mt: '5px', ml: '4px', p: '10px' }}
                onClick={() => remove(index)}
              >
                <DeleteForeverOutlinedIcon />
              </IconButton>
            </Grid>
          </Grid>
        ))}
    </>
  );
}