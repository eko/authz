import axios from "axios";
import { baseUrl } from "service/common/api"
import { paginate, Paginated } from "service/common/paginate";
import { catchError } from "service/error/catch";
import { APIError } from "service/error/model";
import { Principal } from "./model";
import { FilterRequest } from "service/common/filter";
import { SortRequest } from "service/common/sort";
import { PrincipalFormData } from "form/principal";

export const getPrincipals = async (
    token: string,
    page?: number,
    size?: number,
    filter?: FilterRequest,
    sort?: SortRequest,
): Promise<Paginated<Principal> | APIError> => {
    return await catchError<Paginated<Principal>>(async () => {
        return await paginate({
            url: baseUrl() + '/principals',
            token: token,
            page: page,
            size: size,
            filter: filter,
            sort: sort,
        });    
    });
}

export const createPrincipal = async (
    token: string,
    data: PrincipalFormData,
): Promise<Principal | APIError> => {
    return await catchError<Principal>(async () => {
        const url = new URL(baseUrl() + '/principals');

        const response = await axios.post<Principal | APIError>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        return response.data;
    });
}

export const updatePrincipal = async (
    token: string,
    identifier: string,
    data: PrincipalFormData,
): Promise<Principal | APIError> => {
    return await catchError<Principal>(async () => {
        const url = new URL(baseUrl() + `/principals/${identifier}`);

        const response = await axios.put<Principal>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}

export const getPrincipal = async (
    token: string,
    identifier: string,
): Promise<Principal | APIError> => {
    return await catchError<Principal>(async () => {
        const url = new URL(baseUrl() + `/principals/${identifier}`);

        const response = await axios.get<Principal>(url.href, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}

export const deletePrincipal = async (token: string, identifier: string): Promise<boolean | APIError> => {
    return await catchError<boolean>(async () => {
        const url = new URL(baseUrl() + `/principals/${identifier}`);

        const response = await axios.delete<Principal>(url.href, {
            validateStatus: () => true,
            headers: {'Authorization': `Bearer ${token}`},
        });
    
        return response.status === 200;
    });
}
