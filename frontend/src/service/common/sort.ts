export type OrderType = 'asc' | 'desc';

export type SortRequest = {
    field: string
    order: OrderType
};

export const sortRequestToValue = (request: SortRequest): string => {
    return `${request.field}:${request.order}`;
}