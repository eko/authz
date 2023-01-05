export type APIError = {
    error: boolean
    message: string
}

export function isAPIError(toBeDetermined: any): toBeDetermined is APIError {
    return typeof toBeDetermined === 'object'
        && toBeDetermined.hasOwnProperty('error')
        && toBeDetermined.hasOwnProperty('message');
}