import axios from "axios";
import { baseUrl } from "service/common/api"
import { paginate, Paginated } from "service/common/paginate";
import { catchError } from "service/error/catch";
import { APIError } from "service/error/model";
import { Policy } from "./model";
import { FilterRequest } from "service/common/filter";
import { SortRequest } from "service/common/sort";
import { PolicyFormData } from "form/policy";

export const getPolicies = async (
    token: string,
    page?: number,
    size?: number,
    filter?: FilterRequest,
    sort?: SortRequest,
): Promise<Paginated<Policy> | APIError> => {
    return await catchError<Paginated<Policy>>(async () => {
        return await paginate({
            url: baseUrl() + '/policies',
            token: token,
            page: page,
            size: size,
            filter: filter,
            sort: sort,
        });    
    });
}

export const getPolicy = async (
    token: string,
    identifier: string,
): Promise<Policy | APIError> => {
    return await catchError<Policy>(async () => {
        const url = new URL(baseUrl() + `/policies/${identifier}`);

        const response = await axios.get<Policy>(url.href, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}

export const createPolicy = async (
    token: string,
    data: PolicyFormData,
): Promise<Policy | APIError> => {
    return await catchError<Policy>(async () => {
        const url = new URL(baseUrl() + '/policies/');

        const response = await axios.post<Policy | APIError>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        return response.data;
    });
}

export const updatePolicy = async (
    token: string,
    identifier: string,
    data: PolicyFormData,
): Promise<Policy | APIError> => {
    return await catchError<Policy>(async () => {
        const url = new URL(baseUrl() + `/policies/${identifier}`);

        const response = await axios.put<Policy>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}

export const deletePolicy = async (token: string, identifier: string): Promise<boolean | APIError> => {
    return await catchError<boolean>(async () => {
        const url = new URL(baseUrl() + `/policies/${identifier}`);

        const response = await axios.delete<Policy>(url.href, {
            validateStatus: () => true,
            headers: {'Authorization': `Bearer ${token}`},
        });
    
        return response.status === 200;
    });
}
