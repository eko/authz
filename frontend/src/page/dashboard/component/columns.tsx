import { Chip, List, ListItemButton, ListItemIcon, ListItemText } from "@mui/material";
import { getGridBooleanOperators, getGridStringOperators, GridColDef, GridRenderCellParams } from "@mui/x-data-grid";
import LinkIcon from '@mui/icons-material/Link';
import moment from "moment";

export const AuditColumns = (): GridColDef[] => [
    {
      field: 'date',
      headerName: 'Date',
      width: 200,
      sortable: true,
      filterable: false,
      renderCell: (params: GridRenderCellParams) => {
        if (params.row.date.startsWith('0001-01-01')) {
          return (<i>Unknown</i>);
        }

        const date = moment(params.row.date);
        return (
          <div>
            {date.format('L')} at {date.format('LTS')}
          </div>
        )
      },
    },
    {
      field: 'is_allowed',
      headerName: 'Is allowed?',
      width: 100,
      sortable: true,
      filterable: true,
      filterOperators: getGridBooleanOperators(),
      renderCell: (params: GridRenderCellParams) => {
        return (
          <>
            {params.row.is_allowed ? (
              <Chip label='Allowed' color='success' />
            ) : (
              <Chip label='Denied' color='error' />
            )}
          </>
        )
      },
    },
    {
      field: 'policy_id',
      headerName: 'Matched policy',
      width: 250,
      sortable: true,
      filterable: true,
      filterOperators: getGridStringOperators().filter(
        (operator) => operator.value === 'contains',
      ),
      renderCell: (params: GridRenderCellParams) => {
        return params.row.policy_id === '' ? 'None' : (
          <List dense>
            <ListItemButton dense component='a' href={`/policies/edit/${params.row.policy_id}`}>
              <ListItemIcon>
                <LinkIcon fontSize='small' color='primary' />
              </ListItemIcon>
              <ListItemText primary={params.row.policy_id} />
            </ListItemButton>
          </List>
        )
      },
    },
    {
      field: 'principal',
      headerName: 'Principal',
      width: 250,
      sortable: true,
      filterable: true,
      filterOperators: getGridStringOperators().filter(
        (operator) => operator.value === 'contains',
      ),
    },
    {
      field: 'resource_kind',
      headerName: 'Resource kind',
      width: 200,
      sortable: true,
      filterable: true,
      filterOperators: getGridStringOperators().filter(
        (operator) => operator.value === 'contains',
      ),
    },
    {
      field: 'resource_value',
      headerName: 'Resource value',
      width: 200,
      sortable: true,
      filterable: true,
      filterOperators: getGridStringOperators().filter(
        (operator) => operator.value === 'contains',
      ),
    },
    {
      field: 'action',
      headerName: 'Action',
      width: 150,
      sortable: true,
      filterable: true,
      filterOperators: getGridStringOperators().filter(
        (operator) => operator.value === 'contains',
      ),
    },
];