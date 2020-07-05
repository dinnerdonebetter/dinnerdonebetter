import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';

const service = axios.create({
  timeout: 5000,
});

axios.defaults.validateStatus = () => { return true; };

// Request interceptors
service.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    // Add X-Access-Token header to every request, for example
    return config;
  },
  (error) => {
    Promise.reject(error);
  },
);

// Response interceptors
service.interceptors.response.use(
  (response: AxiosResponse) => {
    // Some example codes here:
    // code == 20000: success
    // code == 50001: invalid access token
    // code == 50002: already login in other place
    // code == 50003: access token expired
    // code == 50004: invalid user (user not exist)
    // code == 50005: username or password is incorrect
    // You can change this part for your own usage.

    return response;
  },
  (error) => {
    return Promise.reject(error);
  },
);

export default service;
