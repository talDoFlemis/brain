import { AxiosError } from "axios";

const baseErrorHandler = (error: AxiosError, customMsg?: string) => {
  const data = error.response?.data as { errors: string[] };

  if (data.errors) return Promise.reject(new Error(data.errors.join(",")));
  
  if (customMsg) return Promise.reject(new Error(customMsg))
  
  return Promise.reject(new Error(`${data}`));
}

const handleNetworkError = () => {
  return Promise.reject(new Error("Erro de conexão"));
};

const handleBadRequest = (error: AxiosError) => baseErrorHandler(error);

const handleUnauthorized = (error: AxiosError) => baseErrorHandler(error, "Usuário não autorizado") 

const handleNotFound = (error: AxiosError) => baseErrorHandler(error)

const handleConflict = (error: AxiosError) => baseErrorHandler(error)

const handleUnprocessableContent = (error: AxiosError) => baseErrorHandler(error)

const handleInternalServerError = (error: AxiosError) => baseErrorHandler(error)

export const errorInterceptor = (error: AxiosError) => {
  if (error.message === "Network Error") {
    return handleNetworkError();
  }

  if (error.message === "Request failed with status code 400") {
    return handleBadRequest(error);
  }

  if (error.message === "Request failed with status code 401") {
    return handleUnauthorized(error);
  }

  if (error.message === "Request failed with status code 404") {
    return handleNotFound(error);
  }

  if (error.message === "Request failed with status code 409") {
    return handleConflict(error);
  }

  if (error.message === "Request failed with status code 422") {
    return handleUnprocessableContent(error);
  }

  if (error.message === "Request failed with status code 500") {
    return handleInternalServerError(error);
  }

  return Promise.reject(error);
};
