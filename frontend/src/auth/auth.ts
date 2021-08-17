export class authenticationState {
    accountStatus: string;
    isAuthenticated: boolean;
    userIsServiceAdmin: boolean;
    reputationExplanation: string;

    constructor(
        accountStatus: string = "",
        isAuthenticated: boolean = false,
        userIsServiceAdmin: boolean = false,
        reputationExplanation: string = "",
    ) {
        this.accountStatus = accountStatus;
        this.isAuthenticated = isAuthenticated;
        this.userIsServiceAdmin = userIsServiceAdmin;
        this.reputationExplanation = reputationExplanation;
    }
}

const localStorageAuthKey = "prixfixe_authentication";

export function fetchAuthState(): authenticationState {
    return JSON.parse(localStorage.getItem(localStorageAuthKey) || "{}") as authenticationState;
}

export function setAuthState(state: authenticationState): void {
    localStorage.setItem(localStorageAuthKey, JSON.stringify(state));
}

export function clearAuthState(): void {
    localStorage.removeItem(localStorageAuthKey);
}