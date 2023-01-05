import axios from "axios";
import { RoleFormData } from "../../form/role";
import { baseUrl } from "service/common/api"
import { paginate, Paginated } from "service/common/paginate";
import { catchError } from "service/error/catch";
import { APIError } from "service/error/model";
import { Role } from "./model";
import { FilterRequest } from "service/common/filter";
import { SortRequest } from "service/common/sort";

export const getRoles = async (
    token: string,
    page?: number,
    size?: number,
    filter?: FilterRequest,
    sort?: SortRequest,
): Promise<Paginated<Role> | APIError> => {
    return await catchError<Paginated<Role>>(async () => {
        return await paginate({
            url: baseUrl() + '/roles',
            token: token,
            page: page,
            size: size,
            filter: filter,
            sort: sort,
        });    
    });
}

export const getRole = async (
    token: string,
    identifier: string,
): Promise<Role | APIError> => {
    return await catchError<Role>(async () => {
        const url = new URL(baseUrl() + `/roles/${identifier}`);

        const response = await axios.get<Role>(url.href, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}

export const createRole = async (
    token: string,
    data: RoleFormData,
): Promise<Role | APIError> => {
    return await catchError<Role>(async () => {
        const url = new URL(baseUrl() + '/roles/');

        const response = await axios.post<Role | APIError>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        return response.data;
    });
}

export const updateRole = async (
    token: string,
    identifier: string,
    data: RoleFormData,
): Promise<Role | APIError> => {
    return await catchError<Role>(async () => {
        const url = new URL(baseUrl() + `/roles/${identifier}`);

        const response = await axios.put<Role>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}

export const deleteRole = async (token: string, identifier: string): Promise<boolean | APIError> => {
    return await catchError<boolean>(async () => {
        const url = new URL(baseUrl() + `/roles/${identifier}`);

        const response = await axios.delete<Role>(url.href, {
            validateStatus: () => true,
            headers: {'Authorization': `Bearer ${token}`},
        });
    
        return response.status === 200;
    });
}
