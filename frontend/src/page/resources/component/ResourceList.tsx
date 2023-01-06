import { Button } from '@mui/material';
import Grid from '@mui/material/Grid';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import { useNavigate } from 'react-router';
import DataTable from 'component/DataTable';
import { deleteResource, getResources } from 'service/model/resource';
import { ListColumns } from 'page/resources/component/columns';
import { Box } from '@mui/system';
import { useContext, useState } from 'react';
import { AuthContext } from 'context/auth';
import { useToast } from 'context/toast';
import { GridRenderCellParams } from '@mui/x-data-grid';
import useConfirm from 'hook/confirm';
import { isAPIError } from 'service/error/model';
import { SortRequest } from 'service/common/sort';
import { FilterRequest } from 'service/common/filter';

export default function ResourceList() {
  const navigate = useNavigate();
  const toast = useToast();
  const { user } = useContext(AuthContext);

  const [updatedKey, setUpdatedKey] = useState(0);

  const { ConfirmationDialog, confirm: confirmDelete } = useConfirm();

  const handleOnDelete = async (params: GridRenderCellParams) => {
    const confirmed = await confirmDelete(
      'Delete confirmation',
      `Do you really want to delete resource ${params.row.id}?`,
    );

    if (confirmed) {
      const response = await deleteResource(user?.token!, params.row.id);

      if (isAPIError(response)) {
        toast.error(`An error occurred while trying to delete resource ${params.row.id}: ${response.message}.`);
      } else if (response) {
          toast.success(`Resource ${params.row.id} has been successfully deleted.`);
          setUpdatedKey(updatedKey+1);
      } else {
        toast.error(`An error occurred while trying to delete resource ${params.row.id}.`);
      }
    }
  };

  const columns = ListColumns({
    onDelete: handleOnDelete,
    navigate: navigate,
  });

  const fetcher = (page?: number, size?: number, filter?: FilterRequest, sort?: SortRequest) => {
    return getResources(user?.token!, page, size, filter, sort);
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
          onClick={() => navigate('/resources/create')}
        >
          Create new resource
        </Button>
      </Box>

      <Grid container spacing={3}>
        <DataTable
          key={updatedKey}
          title='Resource list'
          columns={columns}
          fetcher={fetcher}
          sx={{ p: 2 }}
        />
      </Grid>
    </>
  );
}