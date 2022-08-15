import type { ReactNode } from "react";

export interface TableColumn {
  Header: string;
  accessor: string;
  Cell?: (data: any) => ReactNode;
  isNumeric: boolean;
}
