import { AxiosError } from "axios";
import { APIError } from "./model";

export async function catchError<T>(
    callFunc: () => Promise<T | APIError>,
): Promise<T | APIError> {
    try {
        return await callFunc();
    } catch (error) {
        if (error instanceof AxiosError) {
            return {
                error: true,
                message: error.message,
            };
        }

        return {
            error: true,
            message: (error as any).message,
        };
    }
 }