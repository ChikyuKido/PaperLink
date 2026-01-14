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
    D4S: {
        hide: false,
        forceClosed: false,
    },
    Admin: {
        hide: false,
        forceClosed: false,
    },
    AdminSettings: {
        hide: false,
        forceClosed: false,
    },
    AdminIntegrations: {
        hide: false,
        forceClosed: false,
    },
    AdminStatistics: {
        hide: false,
        forceClosed: false,
    },
}
