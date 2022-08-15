import { debounce } from "lodash";
import { useMemo, useState } from "react";

export interface TableFilter {
  column: string;
  value: string;
}

export type InternalTableFilter = Record<string, string>;

export const useTableFilter = (initialFilter: InternalTableFilter = {}) => {
  const [filterParams, setFilter] =
    useState<InternalTableFilter>(initialFilter);

  const setFilterParams = debounce((column: string, value: string) => {
    setFilter({ ...filterParams, [column]: value });
  }, 700);

  const arrayFilterParams = useMemo<TableFilter[]>(() => {
    // TODO type
    // TODO clearable
    return Object.keys(filterParams)
      .map((key) => ({
        column: key,
        value: filterParams[key],
      }))
      .filter((param) => param.value !== "");
  }, [filterParams]);

  return {
    filterParams,
    arrayFilterParams,
    setFilterParams,
  };
};
