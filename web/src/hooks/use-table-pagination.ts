import { useState } from "react";

export interface TablePagination {
  page?: number;
  limit?: number;
}

export const useTablePagination = (initialPage = 1, initialLimit = 10) => {
  const [queryParams, setQueryParams] = useState<TablePagination>({
    page: initialPage,
    limit: initialLimit,
  });

  const setPage = (page: number) => {
    setQueryParams({ ...queryParams, page });
  };

  const setLimit = (limit: number) => {
    setQueryParams({ ...queryParams, limit });
  };

  return {
    queryParams,
    setPage,
    setLimit,
  };
};
