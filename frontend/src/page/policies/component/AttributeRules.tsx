import { Alert, Button, IconButton, TextField } from '@mui/material';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import DeleteForeverOutlinedIcon from '@mui/icons-material/DeleteForeverOutlined';
import { PolicyFormData } from 'form/policy';
import { Control, useFieldArray, UseFormRegister } from 'react-hook-form';

type PolicyFormDataKey = keyof PolicyFormData;

type AttributeRuleProps = {
    control: Control<PolicyFormData>
    fieldName: PolicyFormDataKey
    register: UseFormRegister<PolicyFormData>
};

export default function AttributeRules(props: AttributeRuleProps) {
  const { control, register, fieldName } = props;

  const { fields, append, remove } = useFieldArray({
    control,
    name: fieldName as any,
  });

  return (
    <>
      <Button
        variant='outlined'
        size='small'
        color='primary'
        startIcon={<AddCircleOutlineIcon />}
        onClick={() => append('')}
        sx={{ mb: 2 }}
      >
        Add new rule
      </Button>

      {fields.length === 0 ? (
        <Alert variant="standard" severity="info">
          No attribute rule defined. Your policy will be applied on all selected resources.
        </Alert>
      ) : null}

      {fields.map((item, index) => (
          <div key={item.id} style={{ marginTop: '10px' }}>
            <TextField {...register(fieldName + '.' + index as any)}
              label={`Attribute rule #${index + 1}`}
              placeholder='resource.my_attribute_id == principal.another_attribute_id'
              sx={{ mb: 2, width: '80%' }}
            />
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
          </div>
        ))}
    </>
  );
}