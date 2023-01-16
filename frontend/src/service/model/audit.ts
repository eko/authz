import { baseUrl } from "service/common/api"
import { FilterRequest } from "service/common/filter";
import { paginate, Paginated } from "service/common/paginate";
import { SortRequest } from "service/common/sort";
import { catchError } from "service/error/catch";
import { APIError } from "service/error/model";
import { Audit } from "service/model/model";

export const getAudits = async (
    token: string,
    page?: number,
    size?: number,
    filter?: FilterRequest,
    sort?: SortRequest,
): Promise<Paginated<Audit> | APIError> => {
    return await catchError<Paginated<Audit>>(async () => {
        return await paginate({
            url: baseUrl() + '/audits',
            token: token,
            page: page,
            size: size,
            filter: filter,
            sort: sort,
        });    
    });
}