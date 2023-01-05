export type AuthResponse = {
    access_token: string
    expires_in: number
    token_type: string
    user: AuthResponseUser
}

export type AuthResponseUser = {
    username: string
    created_at: Date
    updated_at: Date
}