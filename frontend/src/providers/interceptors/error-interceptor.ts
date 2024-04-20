import { AxiosError } from "axios";

const handleNetworkError = () => {
  return Promise.reject(new Error("Erro de conexão"));
};

const handleBadRequest = (error: AxiosError) => {
  const data = error.response?.data as { error: string };

  if (data.error) return Promise.reject(new Error(data.error));

  return data;
};

const handleUnauthorized = (error: AxiosError) => {
  const data = error.response?.data as { error: string };

  if (data.error) return Promise.reject(new Error(data.error));

  return Promise.reject(new Error("Usuário não autorizado"));
};

const handleNotFound = (error: AxiosError) => {
  const data = error.response?.data as { error: string };

  if (data.error) return Promise.reject(new Error(data.error));

  return data;
};

const handleUnprocessableContent = (error: AxiosError) => {
  const data = error.response?.data as { error: string };

  if (data.error) return Promise.reject(new Error(data.error));

  return data;
};

const handleInternalServerError = (error: AxiosError) => {
  const data = error.response?.data as { error: string };

  return Promise.reject(
    new Error(data.error || "Erro interno. Tente novamente mais tarde"),
  );
};

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

  if (error.message === "Request failed with status code 422") {
    return handleUnprocessableContent(error);
  }

  if (error.message === "Request failed with status code 500") {
    return handleInternalServerError(error);
  }

  return Promise.reject(error);
};
