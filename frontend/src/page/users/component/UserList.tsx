import { Button } from '@mui/material';
import Grid from '@mui/material/Grid';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import { useNavigate } from 'react-router';
import DataTable from 'component/DataTable';
import { deleteUser, getUsers } from 'service/model/user';
import { ListColumns } from 'page/users/component/columns';
import { Box } from '@mui/system';
import { useContext, useState } from 'react';
import { AuthContext } from 'context/auth';
import { useToast } from 'context/toast';
import { GridRenderCellParams } from '@mui/x-data-grid';
import useConfirm from 'hook/confirm';
import { isAPIError } from 'service/error/model';
import { SortRequest } from 'service/common/sort';
import { FilterRequest } from 'service/common/filter';

export default function UserList() {
  const navigate = useNavigate();
  const toast = useToast();
  const { user } = useContext(AuthContext);

  const [updatedKey, setUpdatedKey] = useState(0);

  const { ConfirmationDialog, confirm: confirmDelete } = useConfirm();

  const handleOnDelete = async (params: GridRenderCellParams) => {
    const confirmed = await confirmDelete(
      'Delete confirmation',
      `Do you really want to delete user ${params.row.username}?`,
    );

    if (confirmed) {
      const response = await deleteUser(user?.token!, params.row.username);

      if (isAPIError(response)) {
        toast.error(`An error occurred while trying to delete user ${params.row.username}: ${response.message}.`);
      } else if (response) {
          toast.success(`User ${params.row.username} has been successfully deleted.`);
          setUpdatedKey(updatedKey+1);
      } else {
        toast.error(`An error occurred while trying to delete user ${params.row.username}.`);
      }
    }
  };

  const columns = ListColumns({
    onDelete: handleOnDelete,
    navigate: navigate,
  });

  const fetcher = (page?: number, size?: number, filter?: FilterRequest, sort?: SortRequest) => {
    return getUsers(user?.token!, page, size, filter, sort);
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
          onClick={() => navigate('/users/create')}
        >
          Create new user
        </Button>
      </Box>

      <Grid container spacing={3}>
        <DataTable
          key={updatedKey}
          title='User list'
          columns={columns}
          fetcher={fetcher}
          sx={{ p: 2 }}
          getRowId={(row) => row.username}
        />
      </Grid>
    </>
  );
}