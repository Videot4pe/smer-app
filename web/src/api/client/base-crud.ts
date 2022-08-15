import { stringify } from "qs";

import type { QueryParams } from "../../hooks/use-table-data";

import { ApiClient } from "./api-client";
import type { MetaData } from "./json-api-document";

const BaseCrud = <TTable, T>(url: string) => {
  const client = new ApiClient(url);

  return {
    list: (queryParams: QueryParams) =>
      client.get<MetaData<TTable[]>>("", {
        params: queryParams,
        paramsSerializer: (params: QueryParams) => {
          return stringify(params);
        },
      }),
    view: (id: number) => client.get<T>(`/${id}`),
    create: (data: T) => client.post<number>("", data),
    update: (id: number, data: T) => client.patch<number>(`/${id}`, data),
    remove: (id: number) => client.delete<number>(`/${id}`),
  };
};

export default BaseCrud;
