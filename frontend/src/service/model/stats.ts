import axios from "axios";
import { baseUrl } from "service/common/api"
import { catchError } from "service/error/catch";
import { APIError } from "service/error/model";
import { Stats } from "./model";

export const getStats = async (
    token: string,
): Promise<Stats | APIError> => {
    return await catchError<Stats>(async () => {
        const url = new URL(baseUrl() + `/stats`);

        const response = await axios.get<Stats>(url.href, {
            validateStatus: () => true,
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}
