import axios from "axios";
import { baseUrl } from "service/common/api"
import { paginate, Paginated } from "service/common/paginate";
import { catchError } from "service/error/catch";
import { APIError } from "service/error/model";
import { User } from "./model";
import { FilterRequest } from "service/common/filter";
import { SortRequest } from "service/common/sort";
import { UserFormData } from "form/user";

export const getUsers = async (
    token: string,
    page?: number,
    size?: number,
    filter?: FilterRequest,
    sort?: SortRequest,
): Promise<Paginated<User> | APIError> => {
    return await catchError<Paginated<User>>(async () => {
        return await paginate({
            url: baseUrl() + '/users',
            token: token,
            page: page,
            size: size,
            filter: filter,
            sort: sort,
        });    
    });
}

export const createUser = async (
    token: string,
    data: UserFormData,
): Promise<User | APIError> => {
    return await catchError<User>(async () => {
        const url = new URL(baseUrl() + '/users');

        const response = await axios.post<User | APIError>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        return response.data;
    });
}

export const getUser = async (
    token: string,
    identifier: string,
): Promise<User | APIError> => {
    return await catchError<User>(async () => {
        const url = new URL(baseUrl() + `/users/${identifier}`);

        const response = await axios.get<User>(url.href, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}

export const deleteUser = async (token: string, identifier: string): Promise<boolean | APIError> => {
    return await catchError<boolean>(async () => {
        const url = new URL(baseUrl() + `/users/${identifier}`);

        const response = await axios.delete<User>(url.href, {
            validateStatus: () => true,
            headers: {'Authorization': `Bearer ${token}`},
        });
    
        return response.status === 200;
    });
}
