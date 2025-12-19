export type SidebarRouteRule = {
    hide?: boolean
    forceClosed?: boolean
}

export const SIDEBAR_ROUTE_RULES: Record<string, SidebarRouteRule> = {
    PDF: {
        hide: false,
        forceClosed: true,
    },
    Auth: {
        hide: true,
        forceClosed: true,
    },
}
