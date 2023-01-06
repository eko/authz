import axios from "axios";
import { baseUrl } from "service/common/api"
import { catchError } from "service/error/catch";
import { APIError } from "service/error/model";
import { CheckFormData } from "form/check";

export type Check = {
    is_allowed: boolean
}

export type CheckResponse = {
    checks: Check[]
}

export const check = async (
    token: string,
    data: CheckFormData,
): Promise<CheckResponse | APIError> => {
    return await catchError<CheckResponse>(async () => {
        const url = new URL(baseUrl() + '/check');

        const response = await axios.post<CheckResponse | APIError>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        return response.data;
    });
}