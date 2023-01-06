import axios from "axios";
import { baseUrl } from "service/common/api"
import { paginate, Paginated } from "service/common/paginate";
import { catchError } from "service/error/catch";
import { APIError } from "service/error/model";
import { Resource } from "./model";
import { FilterRequest } from "service/common/filter";
import { SortRequest } from "service/common/sort";
import { ResourceFormData } from "form/resource";

export const getResources = async (
    token: string,
    page?: number,
    size?: number,
    filter?: FilterRequest,
    sort?: SortRequest,
): Promise<Paginated<Resource> | APIError> => {
    return await catchError<Paginated<Resource>>(async () => {
        return await paginate({
            url: baseUrl() + '/resources',
            token: token,
            page: page,
            size: size,
            filter: filter,
            sort: sort,
        });    
    });
}

export const createResource = async (
    token: string,
    data: ResourceFormData,
): Promise<Resource | APIError> => {
    return await catchError<Resource>(async () => {
        const url = new URL(baseUrl() + '/resources');

        const response = await axios.post<Resource | APIError>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        return response.data;
    });
}

export const updateResource = async (
    token: string,
    identifier: string,
    data: ResourceFormData,
): Promise<Resource | APIError> => {
    return await catchError<Resource>(async () => {
        const url = new URL(baseUrl() + `/resources/${identifier}`);

        const response = await axios.put<Resource>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}

export const getResource = async (
    token: string,
    identifier: string,
): Promise<Resource | APIError> => {
    return await catchError<Resource>(async () => {
        const url = new URL(baseUrl() + `/resources/${identifier}`);

        const response = await axios.get<Resource>(url.href, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}

export const deleteResource = async (token: string, identifier: string): Promise<boolean | APIError> => {
    return await catchError<boolean>(async () => {
        const url = new URL(baseUrl() + `/resources/${identifier}`);

        const response = await axios.delete<Resource>(url.href, {
            validateStatus: () => true,
            headers: {'Authorization': `Bearer ${token}`},
        });
    
        return response.status === 200;
    });
}
