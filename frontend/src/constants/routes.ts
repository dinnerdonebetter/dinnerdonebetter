export const enum backendRoutes {
    USER_REGISTRATION = "/users/",
    LOGIN = "/users/login",
    LOGOUT = "/users/logout",
    VERIFY_2FA_SECRET = "/users/totp_secret/verify",

    USERS = "/api/v1/users",
    USERS_SEARCH = "/api/v1/users/search",

    RECIPES = "/api/v1/recipes",
    RECIPES_SEARCH = "/api/v1/recipes/search",
    RECIPE = "/api/v1/recipes/{}",
    VALID_INGREDIENTS = "/api/v1/valid_ingredients",
    VALID_INGREDIENTS_SEARCH = "/api/v1/valid_ingredients/search",
    VALID_INGREDIENT = "/api/v1/valid_ingredients/{}",
    VALID_INSTRUMENTS = "/api/v1/valid_instruments",
    VALID_INSTRUMENTS_SEARCH = "/api/v1/valid_instruments/search",
    VALID_INSTRUMENT = "/api/v1/valid_instruments/{}",
    VALID_PREPARATIONS = "/api/v1/valid_preparations",
    VALID_PREPARATIONS_SEARCH = "/api/v1/valid_preparations/search",
    VALID_PREPARATION = "/api/v1/valid_preparations/{}",
}

export const enum frontendRoutes {
    ADMIN_VALID_INGREDIENT = "/admin/valid_ingredients/{}",
}

