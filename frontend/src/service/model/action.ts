import axios from "axios";
import { baseUrl } from "service/common/api"
import { paginate, Paginated } from "service/common/paginate";
import { catchError } from "service/error/catch";
import { APIError } from "service/error/model";
import { Action } from "./model";
import { FilterRequest } from "service/common/filter";
import { SortRequest } from "service/common/sort";

export const getActions = async (
    token: string,
    page?: number,
    size?: number,
    filter?: FilterRequest,
    sort?: SortRequest,
): Promise<Paginated<Action> | APIError> => {
    return await catchError<Paginated<Action>>(async () => {
        return await paginate({
            url: baseUrl() + '/actions',
            token: token,
            page: page,
            size: size,
            filter: filter,
            sort: sort,
        });    
    });
}

export const getAction = async (
    token: string,
    identifier: string,
): Promise<Action | APIError> => {
    return await catchError<Action>(async () => {
        const url = new URL(baseUrl() + `/actions/${identifier}`);

        const response = await axios.get<Action>(url.href, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}

export const deleteAction = async (token: string, identifier: string): Promise<boolean | APIError> => {
    return await catchError<boolean>(async () => {
        const url = new URL(baseUrl() + `/actions/${identifier}`);

        const response = await axios.delete<Action>(url.href, {
            validateStatus: () => true,
            headers: {'Authorization': `Bearer ${token}`},
        });
    
        return response.status === 200;
    });
}
