import { useState } from "react";

export interface ReactTableSort {
  desc: boolean;
  id: string;
}

export interface TableSort {
  order: string;
  column: string;
}

export const useTableSort = (initialSort: TableSort[] = []) => {
  const [sortParams, setSort] = useState<TableSort[]>(initialSort);

  const setSortParams = (params: ReactTableSort[]) => {
    setSort(
      params.map((param) => ({
        column: param.id,
        order: param.desc ? "DESC" : "ASC",
      }))
    );
  };

  return {
    sortParams,
    setSortParams,
  };
};
