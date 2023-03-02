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
  forcedSort?: SortRequest
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
  forcedSort,
  title,
  sx,
  getRowId,
  getRowHeight,
}: DataTableProps<T>) {
  const toast = useToast();
  const [loading, setLoading] = useState(false);
  const [sort, setSort] = useState<SortRequest>();
  const [filter, setFilter] = useState<FilterRequest>();
  const [total, setTotal] = useState(0);
  const [rows, setRows] = useState<T[]>([]);
  const [forcedSortModel, setForcedSortModel] = useState<GridSortModel>([]);
  const [paginationModel, setPaginationModel] = useState({
    page: 0,
    pageSize: defaultSize,
  });

  useEffect(() => {
    const fetch = async () => {
      setLoading(true);
      const response = await fetcher(paginationModel.page+1, paginationModel.pageSize, filter, sort);

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
  }, [paginationModel, sort, filter]);

  useEffect(() => {
    if (forcedSort === undefined) {
      return;
    }

    setSort(forcedSort);
    setForcedSortModel([
      {field: forcedSort?.field, sort: forcedSort?.order},
    ]);
  }, [forcedSort]);

  const handleOnSortModelChange = useCallback((sortModel: GridSortModel) => {
    setForcedSortModel(sortModel);

    if (sortModel.length !== 1) {
      setSort(undefined);
      return;
    }

    const currentSort = sortModel[0];
    setSort({ field: currentSort.field, order: currentSort.sort ?? 'asc' });
  }, []);

  const handleOnFilterModelChange = useCallback((filterModel: GridFilterModel) => {
    if (filterModel.items.length !== 1) {
      setFilter(undefined);
      return;
    }

    const currentFilter = filterModel.items[0];

    if (currentFilter.value === undefined) {
      return;
    }

    setFilter({
      field: currentFilter.field,
      operator: currentFilter.operator as FilterOperator,
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
            paginationModel={paginationModel}
            onPaginationModelChange={setPaginationModel}
            onSortModelChange={handleOnSortModelChange}
            paginationMode='server'
            rowCount={total}
            rows={rows}
            sortModel={forcedSortModel}
            sortingMode='server'
            sx={sx}
        />
        </Paper>
      </Grid>
  );
}