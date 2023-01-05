import { Button } from '@mui/material';
import Grid from '@mui/material/Grid';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import { useNavigate } from 'react-router';
import DataTable from 'component/DataTable';
import { deletePolicy, getPolicies } from 'service/model/policy';
import { ListColumns } from 'page/policies/component/columns';
import { Box } from '@mui/system';
import { useContext, useState } from 'react';
import { AuthContext } from 'context/auth';
import { useToast } from 'context/toast';
import { GridRenderCellParams, GridRowHeightParams } from '@mui/x-data-grid';
import useConfirm from 'hook/confirm';
import { isAPIError } from 'service/error/model';
import { SortRequest } from 'service/common/sort';
import { FilterRequest } from 'service/common/filter';

export default function PolicyList() {
  const navigate = useNavigate();
  const toast = useToast();
  const { user } = useContext(AuthContext);

  const [updatedKey, setUpdatedKey] = useState(0);

  const { ConfirmationDialog, confirm: confirmDelete } = useConfirm();

  const handleOnDelete = async (params: GridRenderCellParams) => {
    const confirmed = await confirmDelete(
      'Delete confirmation',
      `Do you really want to delete policy ${params.row.id}?`,
    );

    if (confirmed) {
      const response = await deletePolicy(user?.token!, params.row.id);

      if (isAPIError(response)) {
        toast.error(`An error occurred while trying to delete policy ${params.row.id}: ${response.message}.`);
      } else if (response) {
          toast.success(`Policy ${params.row.id} has been successfully deleted.`);
          setUpdatedKey(updatedKey+1);
      } else {
        toast.error(`An error occurred while trying to delete policy ${params.row.id}.`);
      }
    }
  };

  const columns = ListColumns({
    navigate,
    onDelete: handleOnDelete,
  });

  const fetcher = (page?: number, size?: number, filter?: FilterRequest, sort?: SortRequest) => {
    return getPolicies(user?.token!, page, size, filter, sort);
  };

  return (
    <>
      <ConfirmationDialog />
      <Box mb={3} display='flex' justifyContent='flex-start'>
        <Button
          variant='contained'
          size='medium'
          color='primary'
          startIcon={<AddCircleOutlineIcon />}
          onClick={() => navigate('/policies/create')}
        >
          Create a new policy
        </Button>
      </Box>

      <Grid container spacing={3}>
        <DataTable
          key={updatedKey}
          title='Policies list'
          columns={columns}
          fetcher={fetcher}
          sx={{ p: 2 }}
          getRowHeight={(params: GridRowHeightParams) => {
            const resourcesNumber = params.model.resources?.length || 1;
            const actionsNumber = params.model.actions?.length || 1;

            const max = resourcesNumber > actionsNumber ? resourcesNumber : actionsNumber;

            return params.densityFactor * (45 * max);
          }}
        />
      </Grid>
    </>
  );
}