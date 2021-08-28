export class authenticationState {
    accountStatus: string;
    isAuthenticated: boolean;
    userIsServiceAdmin: boolean;
    reputationExplanation: string;
    expiresOn: number;

    constructor(
        accountStatus: string = "",
        isAuthenticated: boolean = false,
        userIsServiceAdmin: boolean = false,
        reputationExplanation: string = "",
        expiresOn: number = 0, //parseInt(Date.now() / 1000),
    ) {
        this.accountStatus = accountStatus;
        this.isAuthenticated = isAuthenticated;
        this.userIsServiceAdmin = userIsServiceAdmin;
        this.reputationExplanation = reputationExplanation;
        this.expiresOn = expiresOn;
    }
}

const invalidAuthenticationState = new authenticationState();

const localStorageAuthKey = "prixfixe_authentication";

export function isAuthenticated(): boolean {
    return fetchAuthState().isAuthenticated
}

export function fetchAuthState(): authenticationState {
    return JSON.parse(localStorage.getItem(localStorageAuthKey) || JSON.stringify(invalidAuthenticationState)) as authenticationState;
}

export function setAuthState(state: authenticationState): void {
    localStorage.setItem(localStorageAuthKey, JSON.stringify(state));
}

export function clearAuthState(): void {
    localStorage.removeItem(localStorageAuthKey);
}