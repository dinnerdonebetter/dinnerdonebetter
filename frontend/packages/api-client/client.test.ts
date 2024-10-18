import axios from "axios";
import AxiosMockAdapter from "axios-mock-adapter";

import { APIResponse, ValidIngredient, ValidInstrument } from "@dinnerdonebetter/models";

import { DinnerDoneBetterAPIClient } from "./client.gen";

const mock = new AxiosMockAdapter(axios);
const baseURL = "http://things.stuff";
const client = new DinnerDoneBetterAPIClient(baseURL, 'test-token');

beforeEach(() => mock.reset());

describe('basic', () => {
    it("should return a valid ingredient", () => {
        const exampleResponse = new ValidIngredient();
        mock.onGet(`${baseURL}/api/v1/valid_ingredients/test`).reply(200, exampleResponse);        

        client.getValidIngredient('test').then((response: APIResponse<ValidIngredient>) => {
            expect(response).toEqual(exampleResponse);
        });
    });

    it("should return a valid instrument", () => {
        const exampleResponse = new ValidInstrument();
        mock.onGet(`${baseURL}/api/v1/valid_instruments/test`).reply(200, exampleResponse);        

        client.getValidInstrument('test').then((response: APIResponse<ValidInstrument>) => {
            expect(response).toEqual(exampleResponse);
        });
    });
});