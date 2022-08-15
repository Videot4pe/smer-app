export interface Meta {
  totalItems: number;
  totalPages: number;
}

export interface MetaData<T> {
  data: T;
  meta: Meta;
}

export interface JsonApiDocument<T> {
  data: T;
  status?: number;
}

export interface JsonApiDocumentError {
  errors: JsonApiError[];
}
export interface JsonApiError {
  id?: string;
  status?: string;
  code?: string;
  title?: string;
  detail?: string;
}
