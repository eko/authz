export type Action = {
    id: string
    created_at: Date
    updated_at: Date
}

export type Resource = {
    id: string
    kind: string
    value: string
    is_locked: boolean
    created_at: Date
    updated_at: Date
}

export type Policy = {
    id: string
    actions: Action[]
    resources: Resource[]
    created_at: Date
    updated_at: Date
}

export type Role = {
    id: string
    name: string
    policies: Policy[]
    created_at: Date
    updated_at: Date
}

export type Principal = {
    id: string
    name: string
    roles: Role[]
    created_at: Date
    updated_at: Date
}

export type Group = {
    id: string
    name: string
    created_at: Date
    updated_at: Date
}

export type User = {
    id: string
    name: string
    created_at: Date
    updated_at: Date
}