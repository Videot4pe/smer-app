export interface JsonErrorResponse {
  status: number;
  statusText: string;
  data: {
    error: ApiError;
  };
}

export interface ApiError {
  status: number;
  title: string;
}
