export type SidebarRouteRule = {
    hide?: boolean
    forceClosed?: boolean
}

export const SIDEBAR_ROUTE_RULES: Record<string, SidebarRouteRule> = {
    '/pdf': {
        hide: false,
        forceClosed: true,
    },
    '/auth': {
        hide: true,
        forceClosed: true,
    },
}
