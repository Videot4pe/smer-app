import { useEffect, useState } from "react";

import type { MetaData } from "../api/client/json-api-document";
import { useErrorHandler } from "../utils/handle-get-error";

import type { TableFilter } from "./use-table-filter";
import type { TablePagination } from "./use-table-pagination";
import type { TableSort } from "./use-table-sort";

export interface QueryParams extends TablePagination {
  sort?: TableSort[];
  filter?: TableFilter[];
}

export const useTableData = <T>(
  repository: (queryParams: QueryParams) => Promise<MetaData<T[]>>,
  queryParams: TablePagination = {},
  filterParams: TableFilter[] = [],
  sortParams: TableSort[] = []
) => {
  const errorHandler = useErrorHandler();
  const [isLoading, setIsLoading] = useState(false);

  const [data, setData] = useState<T[]>([]);
  const [meta, setMeta] = useState({
    totalItems: 0,
    totalPages: 1,
  });

  const fetch = () => {
    // TODO useLoading hook
    setIsLoading(true);
    repository({ ...queryParams, filter: filterParams, sort: sortParams })
      .then((payload) => {
        setData(payload.data);
        setMeta(payload.meta);
      })
      .catch(errorHandler)
      .finally(() => setIsLoading(false));
  };

  useEffect(fetch, [queryParams, filterParams, sortParams]);

  return {
    data,
    meta,
    isLoading,
    fetch,
  };
};
