import axios from "axios";
import { baseUrl } from "service/common/api"
import { paginate, Paginated } from "service/common/paginate";
import { catchError } from "service/error/catch";
import { APIError } from "service/error/model";
import { Client } from "./model";
import { FilterRequest } from "service/common/filter";
import { SortRequest } from "service/common/sort";
import { ClientFormData } from "form/client";

export const getClients = async (
    token: string,
    page?: number,
    size?: number,
    filter?: FilterRequest,
    sort?: SortRequest,
): Promise<Paginated<Client> | APIError> => {
    return await catchError<Paginated<Client>>(async () => {
        return await paginate({
            url: baseUrl() + '/clients',
            token: token,
            page: page,
            size: size,
            filter: filter,
            sort: sort,
        });    
    });
}

export const createClient = async (
    token: string,
    data: ClientFormData,
): Promise<Client | APIError> => {
    return await catchError<Client>(async () => {
        const url = new URL(baseUrl() + '/clients');

        const response = await axios.post<Client | APIError>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        return response.data;
    });
}

export const getClient = async (
    token: string,
    identifier: string,
): Promise<Client | APIError> => {
    return await catchError<Client>(async () => {
        const url = new URL(baseUrl() + `/clients/${identifier}`);

        const response = await axios.get<Client>(url.href, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}

export const deleteClient = async (token: string, identifier: string): Promise<boolean | APIError> => {
    return await catchError<boolean>(async () => {
        const url = new URL(baseUrl() + `/clients/${identifier}`);

        const response = await axios.delete<Client>(url.href, {
            validateStatus: () => true,
            headers: {'Authorization': `Bearer ${token}`},
        });
    
        return response.status === 200;
    });
}
