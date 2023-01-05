export type FilterOperator = 'contains' | 'is';

export type FilterRequest = {
    field: string
    value: string
    operator: FilterOperator
};

export const filterRequestToValue = (request: FilterRequest): string => {
    return `${request.field}:${request.operator}:${request.value}`;
}