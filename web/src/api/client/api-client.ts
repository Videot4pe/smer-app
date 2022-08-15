import type { AxiosInstance, AxiosRequestConfig } from "axios";
import axios from "axios";

import AuthApi from "../auth-api";

import type { JsonApiDocument } from "./json-api-document";

export const getClient = (url: string) => {
  const baseUrl = process.env.SERVER_URL;
  const options: AxiosRequestConfig = {
    // WithCredentials: true,
    baseURL: `${baseUrl}/api/${url}`,
  };

  const client = axios.create(options);

  client.interceptors.request.use(
    (config) => {
      if (config.headers) {
        const token = localStorage.getItem("jwt-token");
        // eslint-disable-next-line no-param-reassign
        config.headers.Authorization = `Bearer ${token}`;
      }

      return config;
    },
    (requestError) => Promise.reject(requestError)
  );

  client.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config;
      if (error.response.status === 403 && !originalRequest.retry) {
        localStorage.removeItem("jwt-token");
        originalRequest.retry = true;
        const refreshToken = localStorage.getItem("refresh-token");
        if (refreshToken) {
          await AuthApi.refresh(refreshToken)
            .then((payload) => {
              localStorage.setItem("jwt-token", payload.token);
              localStorage.setItem("refresh-token", payload.refreshToken);
            })
            .catch(() => {
              window.location.href = "localhost:3000/signin";
            });
          return client(originalRequest);
        }
      }
      return Promise.reject(error.response);
    }
  );

  return client;
};

class ApiClient {
  private client: AxiosInstance;

  constructor(baseUrl: string) {
    this.client = getClient(baseUrl);
  }

  get<T>(url: string, conf = {}) {
    return this.client
      .get<JsonApiDocument<T>>(this.client.defaults.baseURL + url, conf)
      .then((response) => Promise.resolve(response.data.data))
      .catch((error) => Promise.reject(error));
  }

  delete<T>(url: string, conf = {}) {
    return this.client
      .delete<JsonApiDocument<T>>(this.client.defaults.baseURL + url, conf)
      .then((response) => Promise.resolve(response.data.data))
      .catch((error) => Promise.reject(error));
  }

  head<T>(url: string, conf = {}) {
    return this.client
      .head<JsonApiDocument<T>>(this.client.defaults.baseURL + url, conf)
      .then((response) => Promise.resolve(response.data.data))
      .catch((error) => Promise.reject(error));
  }

  options<T>(url: string, conf = {}) {
    return this.client
      .options<JsonApiDocument<T>>(this.client.defaults.baseURL + url, conf)
      .then((response) => Promise.resolve(response.data.data))
      .catch((error) => Promise.reject(error));
  }

  post<T>(url: string, data = {}, conf = {}) {
    return this.client
      .post<JsonApiDocument<T>>(this.client.defaults.baseURL + url, data, conf)
      .then((response) => Promise.resolve(response.data.data))
      .catch((error) => Promise.reject(error));
  }

  put<T>(url: string, data = {}, conf = {}) {
    return this.client
      .put<T>(this.client.defaults.baseURL + url, data, conf)
      .then((response) => Promise.resolve(response.data))
      .catch((error) => Promise.reject(error));
  }

  patch<T>(url: string, data = {}, conf = {}) {
    return this.client
      .patch<T>(this.client.defaults.baseURL + url, data, conf)
      .then((response) => Promise.resolve(response.data))
      .catch((error) => Promise.reject(error));
  }
}

export { ApiClient };
