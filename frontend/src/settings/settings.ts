class Settings {
    API_SERVER_URL: string

    constructor() {
        this.API_SERVER_URL = import.meta.env.VITE_API_SERVER_URL || '';
    }
}

export const settings = new Settings();
