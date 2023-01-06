import { Typography } from '@mui/material';
import { SxProps } from '@mui/system';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import { DataGrid, GridColDef, GridFilterModel, GridRowClassNameParams, GridRowHeightParams, GridRowHeightReturnValue, GridSortModel } from '@mui/x-data-grid';

import { useCallback, useEffect, useState } from 'react';
import { useToast } from 'context/toast';
import { Paginated } from 'service/common/paginate';
import { APIError, isAPIError } from 'service/error/model';

import './DataTable.css';
import { SortRequest } from 'service/common/sort';
import { FilterOperator, FilterRequest } from 'service/common/filter';

type FetcherFunc<T> = (
  page?: number,
  size?: number,
  filter?: FilterRequest,
  sort?: SortRequest,
) => Promise<Paginated<T> | APIError>

type DataTableProps<T> = {
  fetcher: FetcherFunc<T>
  columns: GridColDef[]
  defaultSize?: number
  title?: string
  sx?: SxProps
  getRowId?: getRowId<T>
  getRowHeight?: getRowHeight
}

type getRowHeight = (params: GridRowHeightParams) => GridRowHeightReturnValue;
type getRowId<T> = (row: T) => string;

export default function DataTable<T>({
  fetcher,
  columns,
  defaultSize = 10,
  title,
  sx,
  getRowId,
  getRowHeight,
}: DataTableProps<T>) {
  const toast = useToast();
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(0);
  const [size, setSize] = useState(defaultSize);
  const [sort, setSort] = useState<SortRequest>();
  const [filter, setFilter] = useState<FilterRequest>();
  const [total, setTotal] = useState(0);
  const [rows, setRows] = useState<T[]>([]);

  useEffect(() => {
    const fetch = async () => {
      setLoading(true);
      const response = await fetcher(page+1, size, filter, sort);

      if (isAPIError(response)) {
        toast.error(`Impossible de charger les données : ${response.message}`);
      } else {
        setTotal(response.total);
        setRows(response.data);
      }

      setLoading(false);
    };

    fetch();
  // eslint-disable-next-line
  }, [size, page, sort, filter]);

  const handleOnPageChange = (page: number) => setPage(page);

  const handleOnSortModelChange = useCallback((sortModel: GridSortModel) => {
    if (sortModel.length !== 1) {
      return;
    }

    const currentSort = sortModel[0];
    setSort({ field: currentSort.field, order: currentSort.sort ?? 'asc' });
  }, []);

  const handleOnFilterModelChange = useCallback((filterModel: GridFilterModel) => {
    if (filterModel.items.length !== 1) {
      return;
    }

    const currentFilter = filterModel.items[0];

    if (currentFilter.value === undefined) {
      return;
    }

    setFilter({
      field: currentFilter.columnField,
      operator: currentFilter.operatorValue as FilterOperator,
      value: currentFilter.value,
    });
  }, []);

  return (
      <Grid item xs={12} md={12} lg={12}>
        <Paper
        sx={{
            display: 'flex',
            flexDirection: 'column',
        }}
        >

        {title ? (<Typography variant='h5' sx={{ pt: 2, pl: 2, pb: 2 }}>
            {title}
        </Typography>) : null}

        <DataGrid
            autoHeight
            columns={columns}
            density='standard'
            filterMode='server'
            getRowClassName={(params: GridRowClassNameParams) => {
              return params.indexRelativeToCurrentPage % 2 === 0 ? 'datagrid-row-lightgrey' : ''
            }}
            getRowId={getRowId}
            getRowHeight={getRowHeight}
            hideFooterSelectedRowCount
            loading={loading}
            onFilterModelChange={handleOnFilterModelChange}
            onPageChange={handleOnPageChange}
            onPageSizeChange={(newPageSize: number) => setSize(newPageSize)}
            onSortModelChange={handleOnSortModelChange}
            page={page}
            pageSize={size}
            paginationMode='server'
            rowCount={total}
            rows={rows}
            rowsPerPageOptions={[5,10,20,30,50,100]}
            sortingMode='server'
            sx={sx}
        />
        </Paper>
      </Grid>
  );
}