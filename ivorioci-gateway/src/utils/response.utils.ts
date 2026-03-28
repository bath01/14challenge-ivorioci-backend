export interface ApiSuccess<T> {
  success: true;
  data: T;
  timestamp: string;
}

export interface ApiError {
  success: false;
  error: {
    code: string;
    message: string;
  };
  timestamp: string;
}

export type ApiResponse<T> = ApiSuccess<T> | ApiError;

export function successResponse<T>(data: T): ApiSuccess<T> {
  return { success: true, data, timestamp: new Date().toISOString() };
}

export function errorResponse(code: string, message: string): ApiError {
  return { success: false, error: { code, message }, timestamp: new Date().toISOString() };
}
