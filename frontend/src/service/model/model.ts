export type Action = {
    id: string
    created_at: Date
    updated_at: Date
}

export type Client = {
    client_id: string
    client_secret: string
    name: string
    created_at: Date
    updated_at: Date
}

export type Policy = {
    id: string
    actions: Action[]
    resources: Resource[]
    attribute_rules?: string[]
    created_at: Date
    updated_at: Date
}

export type PrincipalAttribute = {
    key: string
    value: string
}

export type Principal = {
    id: string
    name: string
    roles: Role[]
    attributes?: PrincipalAttribute[]
    created_at: Date
    updated_at: Date
}

export type ResourceAttribute = {
    key: string
    value: string
}

export type Resource = {
    id: string
    kind: string
    value: string
    attributes?: ResourceAttribute[]
    is_locked: boolean
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

export type StatsDay = {
    id: string
    date: string
    checks_allowed_number: number
    checks_denied_number: number
}

export type Stats = StatsDay[]

export type User = {
    username: string
    password?: string
    created_at: Date
    updated_at: Date
}