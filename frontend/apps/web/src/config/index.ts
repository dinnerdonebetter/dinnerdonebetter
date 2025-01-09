import { loadFile } from "@dinnerdonebetter/configloader";

interface Config {
    APIEndpoint: string;
    SegmentAPIToken: string;
    CookieEncryptionKey: string;
    CookieEncryptionIV: string;
    APIOAuth2ClientID: string;
    APIOAuth2ClientSecret: string;
}

const config= loadFile<Config>("routingcfg.json");

export default config;
