import axios from "axios";
import { FilterRequest, filterRequestToValue } from "./filter";
import { SortRequest, sortRequestToValue } from "./sort";

export type PaginateRequest = {
    url: string
    token: string
    page?: number
    size?: number
    filter?: FilterRequest
    sort?: SortRequest
}

export type Paginated<T> = {
    total: number
    page: number
    size: number
    data: T[]
};

export const paginate = async <T>(request: PaginateRequest): Promise<Paginated<T>> => {
    const url = new URL(request.url);

    if (request.page !== undefined) {
        url.searchParams.set('page', request.page.toString());
    }

    if (request.size !== undefined) {
        url.searchParams.set('size', request.size.toString());
    }

    if (request.filter) {
        url.searchParams.set('filter', filterRequestToValue(request.filter));
    }

    if (request.sort) {
        url.searchParams.set('sort', sortRequestToValue(request.sort));
    }

    const response = await axios.get<Paginated<T>>(url.href, {
        validateStatus: () => true,
        headers: {
            'Authorization': `Bearer ${request.token}`,
            'Content-Type': 'application/json',
        },
    });
    
    return response.data;
}