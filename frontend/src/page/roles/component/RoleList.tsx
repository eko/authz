import { Button } from '@mui/material';
import Grid from '@mui/material/Grid';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import { useNavigate } from 'react-router';
import DataTable from 'component/DataTable';
import { deleteRole, getRoles } from 'service/model/role';
import { ListColumns } from 'page/roles/component/columns';
import { Box } from '@mui/system';
import { useContext, useState } from 'react';
import { AuthContext } from 'context/auth';
import { useToast } from 'context/toast';
import { GridRenderCellParams, GridRowHeightParams } from '@mui/x-data-grid';
import useConfirm from 'hook/confirm';
import { isAPIError } from 'service/error/model';
import { SortRequest } from 'service/common/sort';
import { FilterRequest } from 'service/common/filter';

export default function RoleList() {
  const navigate = useNavigate();
  const toast = useToast();
  const { user } = useContext(AuthContext);

  const [updatedKey, setUpdatedKey] = useState(0);

  const { ConfirmationDialog, confirm: confirmDelete } = useConfirm();

  const handleOnDelete = async (params: GridRenderCellParams) => {
    const confirmed = await confirmDelete(
      'Delete confirmation',
      `Do you really want to delete role ${params.row.id}?`,
    );

    if (confirmed) {
      const response = await deleteRole(user?.token!, params.row.id);

      if (isAPIError(response)) {
        toast.error(`An error occurred while trying to delete role ${params.row.id}: ${response.message}.`);
      } else if (response) {
          toast.success(`Role ${params.row.id} has been successfully deleted.`);
          setUpdatedKey(updatedKey+1);
      } else {
        toast.error(`An error occurred while trying to delete role ${params.row.id}.`);
      }
    }
  };

  const columns = ListColumns({
    navigate,
    onDelete: handleOnDelete,
  });

  const fetcher = (page?: number, size?: number, filter?: FilterRequest, sort?: SortRequest) => {
    return getRoles(user?.token!, page, size, filter, sort);
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
          onClick={() => navigate('/roles/create')}
        >
          Create new role
        </Button>
      </Box>

      <Grid container spacing={3}>
        <DataTable
          key={updatedKey}
          title='Role list'
          columns={columns}
          fetcher={fetcher}
          sx={{ p: 2 }}
          getRowHeight={(params: GridRowHeightParams) => {
            const policiesNumber = params.model.policies?.length || 1;
            return params.densityFactor * (45 * policiesNumber);
          }}
        />
      </Grid>
    </>
  );
}