import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import DataTable from './DataTable';
import { Paginated } from 'service/common/paginate';
import { APIError } from 'service/error/model';
import { getGridStringOperators, GridColDef } from '@mui/x-data-grid';
import { ToastProvider } from 'context/toast';
import { FilterRequest } from 'service/common/filter';
import { SortRequest } from 'service/common/sort';

type Item = {
    id: string,
    name: string,
};

const fetcher = async (
    page?: number,
    size?: number,
    filter?: FilterRequest,
    sort?: SortRequest,
): Promise<Paginated<Item> | APIError> => {
    const items = [
        {id: '1', name: 'item-1'},
        {id: '2', name: 'item-2'},
        {id: '3', name: 'item-3'},
        {id: '4', name: 'item-4'},
        {id: '5', name: 'item-5'},
        {id: '6', name: 'item-6'},
        {id: '7', name: 'item-7'},
        {id: '8', name: 'item-8'},
        {id: '9', name: 'item-9'},
        {id: '10', name: 'item-10'},
    ];

    if (filter !== undefined) {
        switch (filter.operator) {
            case 'contains':
                items.filter((item: any) => (item[filter.field] as string).includes(filter.value));
            break;
        }
    }

    if (sort !== undefined) {
        switch (sort.field) {
            case 'name':
                if (sort.order === 'desc') {
                    items.reverse();
                }
                break;
        }
    }

    const subitems = items.slice((page!-1)*size!, size!*page!);

    return {
        total: items.length,
        page: page!,
        size: size!,
        data: subitems,
    };
};

test('DataTable: fetch data, render list and paginate', async () => {
    // Given
    const columns: GridColDef[] = [
        { field: 'id', headerName: 'ID' },
        { field: 'name', headerName: 'Nom' },
    ];

    const Wrapper = () => {
        return (
            <ToastProvider>
                <DataTable
                    title='Test list'
                    columns={columns}
                    fetcher={fetcher}
                    defaultSize={5}
                />
            </ToastProvider>
        )
    };

    // Render
    render(
        <Wrapper />
    );

    // When waiting for results to be displayed
    // the table should contain the first 5 elements.
    await waitFor(() => {
        expect(screen.getByText('item-5')).toBeInTheDocument();
    });

    let grid = screen.getByRole('grid');

    expect(grid).toHaveTextContent('item-1');
    expect(grid).toHaveTextContent('item-2');
    expect(grid).toHaveTextContent('item-3');
    expect(grid).toHaveTextContent('item-4');
    expect(grid).toHaveTextContent('item-5');

    expect(grid).not.toHaveTextContent('item-6');

    // When click on next page
    // the table should contain the next (and last) 5 elements
    await userEvent.click(screen.getByTitle('Go to next page'));
    await waitFor(() => {
        expect(screen.getByText('item-6')).toBeInTheDocument();
    });

    grid = screen.getByRole('grid');

    expect(grid).not.toHaveTextContent('item-5');

    expect(grid).toHaveTextContent('item-6');
    expect(grid).toHaveTextContent('item-7');
    expect(grid).toHaveTextContent('item-8');
    expect(grid).toHaveTextContent('item-9');
    expect(grid).toHaveTextContent('item-10');
});

test('DataTable: sort ascending / descending', async () => {
    // Given
    const columns: GridColDef[] = [
        { field: 'id', headerName: 'ID' },
        { field: 'name', headerName: 'Nom' },
    ];

    const Wrapper = () => {
        return (
            <ToastProvider>
                <DataTable
                    title='Test list'
                    columns={columns}
                    fetcher={fetcher}
                    defaultSize={5}
                />
            </ToastProvider>
        )
    };

    // Render
    render(
        <Wrapper />
    );

    // When waiting for results to be displayed
    // the table should contain the first 5 elements.
    await waitFor(() => {
        expect(screen.getByText('item-5')).toBeInTheDocument();
    });

    let rows = screen.getAllByRole('row');

    expect(rows[1].getElementsByClassName('MuiDataGrid-cellContent')[1].innerHTML).toEqual('item-1');
    expect(rows[2].getElementsByClassName('MuiDataGrid-cellContent')[1].innerHTML).toEqual('item-2');
    expect(rows[3].getElementsByClassName('MuiDataGrid-cellContent')[1].innerHTML).toEqual('item-3');
    expect(rows[4].getElementsByClassName('MuiDataGrid-cellContent')[1].innerHTML).toEqual('item-4');
    expect(rows[5].getElementsByClassName('MuiDataGrid-cellContent')[1].innerHTML).toEqual('item-5');

    // Open name field submenu
    const nameSubmenu = screen.getAllByRole('columnheader')
        .find(item => item.attributes.getNamedItem('data-field')?.value === 'name');

    await userEvent.click(nameSubmenu?.getElementsByClassName('MuiDataGrid-menuIconButton')[0]!);

    // Click on descending filter
    const descSort = screen.getAllByRole('menuitem')
        .find(item => item.attributes.getNamedItem('data-value')?.value === 'desc');
    await userEvent.click(descSort!);

    await waitFor(() => {
        rows = screen.getAllByRole('row');

        expect(rows[1].getElementsByClassName('MuiDataGrid-cellContent')[1].innerHTML).toEqual('item-10');
        expect(rows[2].getElementsByClassName('MuiDataGrid-cellContent')[1].innerHTML).toEqual('item-9');
        expect(rows[3].getElementsByClassName('MuiDataGrid-cellContent')[1].innerHTML).toEqual('item-8');
        expect(rows[4].getElementsByClassName('MuiDataGrid-cellContent')[1].innerHTML).toEqual('item-7');
        expect(rows[5].getElementsByClassName('MuiDataGrid-cellContent')[1].innerHTML).toEqual('item-6');
    });
});