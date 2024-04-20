import axios, { AxiosRequestConfig } from "axios";
import { BACKEND_API_URL } from "@/utils/constants";
import { errorInterceptor } from "./interceptors/error-interceptor";

const provider = axios.create({
  baseURL: BACKEND_API_URL,
});

provider.interceptors.response.use(
  (response) => response,
  (error) => errorInterceptor(error),
);

const useGet = async <T>(
  path: string,
  options?: AxiosRequestConfig,
): Promise<T> => {
  try {
    const response = await provider.get<T>(path, options);

    return response.data;
  } catch (error) {
    throw error;
  }
};

const usePost = async <T, D>(
  path: string,
  data?: D,
  options?: AxiosRequestConfig,
): Promise<T> => {
  try {
    const response = await provider.post<T>(path, data, options);

    return response.data;
  } catch (error) {
    throw error;
  }
};

const usePut = async <T, D>(
  path: string,
  data?: D,
  options?: AxiosRequestConfig,
): Promise<T> => {
  try {
    const response = await provider.put<T>(path, data, options);

    return response.data;
  } catch (error) {
    throw error;
  }
};

const useDelete = async <T>(
  path: string,
  options?: AxiosRequestConfig,
): Promise<T> => {
  try {
    const response = await provider.delete<T>(path, options);

    return response.data;
  } catch (error) {
    throw error;
  }
};

const apiProvider = {
  provider,
  useGet,
  usePost,
  usePut,
  useDelete,
};

export default apiProvider;
