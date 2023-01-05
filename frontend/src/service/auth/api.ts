import axios from "axios"
import { SigninFormData } from "form/signin";
import { baseUrl } from "service/common/api"
import { catchError } from "service/error/catch";
import { APIError } from "service/error/model";
import { AuthResponse } from "service/auth/model";

export const signIn = async (
    data: SigninFormData,
): Promise<AuthResponse | APIError> => {
    return await catchError<AuthResponse>(async () => {
        const url = new URL(baseUrl() + '/auth');

        const response = await axios.post<AuthResponse>(url.href, data, {
            validateStatus: () => true,
            headers: {
                'Content-Type': 'application/json',
            },
        });
    
        return response.data;
    });
}
