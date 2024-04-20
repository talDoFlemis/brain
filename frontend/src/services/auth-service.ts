import apiProvider from "@/providers/api-provider";

type SignUpRequest = {
  username: string;
  email: string;
  password: string;
};

type SignInRequest = {
  username: string;
  password: string;
};

type RefreshRequest = {
  refresh_token: string;
};

type AccessTokens = {
  access_token: string;
  refresh_token: string;
  expire_at: string;
};

type SignUpResponse = AccessTokens;

type SignInResponse = SignUpResponse;

type RefreshResponse = SignUpResponse;

type UserInfoResponse = {
  id: string;
  username: string;
  email: string;
};

const signUp = async (data: SignUpRequest): Promise<SignUpResponse> => {
  try {
    const response = await apiProvider.usePost<SignUpResponse, SignUpRequest>(
      "/auth/register",
      data,
    );

    return response;
  } catch (error) {
    throw error;
  }
};

const signIn = async (data: SignInRequest): Promise<SignUpResponse> => {
  try {
    const response = await apiProvider.usePost<SignInResponse, SignInRequest>(
      "/auth/",
      data,
    );

    return response;
  } catch (error) {
    throw error;
  }
};

const refreshToken = async (data: RefreshRequest): Promise<RefreshResponse> => {
  try {
    const response = await apiProvider.usePost<RefreshResponse, RefreshRequest>(
      "/auth/refresh",
      data,
    );

    return response;
  } catch (error) {
    throw error;
  }
};

const getUserInfo = async (access_token: string): Promise<UserInfoResponse> => {
  try {
    const response = await apiProvider.useGet<UserInfoResponse>(
      "/auth/userinfo",
      {
        headers: {
          Authorization: `Bearer ${access_token}`,
        },
      },
    );

    return response;
  } catch (error) {
    throw error;
  }
};

const authService = {
  signIn,
  signUp,
  refreshToken,
  getUserInfo,
};

export default authService;
