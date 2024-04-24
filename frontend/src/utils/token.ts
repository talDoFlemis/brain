import { JWT } from "next-auth/jwt";
import { User } from "next-auth";
import authService from "@/services/auth-service";

export const REFRESH_TOKEN_ERROR = "RefreshAccessTokenError";
export const USER_INFO_ERROR = "CouldNotGerUserInfo";

export function isTokenExpired(expire_at: string): boolean {
  const parsed = new Date(expire_at.split("m=")[0]);

  if (Date.now() < parsed.getTime()) return false;

  return true;
}

export async function handleTokenRefreshment(token: JWT): Promise<JWT> {
  try {
    const response = await authService.refreshToken({
      refresh_token: token.refresh_token,
    });

    const userInfo = await authService.getUserInfo(response.access_token);

    token.id = userInfo.id;
    token.username = userInfo.username;
    token.email = userInfo.email;
    token.refresh_token = response.refresh_token;
    token.access_token = response.access_token;
    token.expire_at = response.expire_at;

    return token;
  } catch (error) {
    // Pass error to token in order to it be caught in the client side
    console.log(`[Error] Could not refresh access token: ${error}`);
    token.error = REFRESH_TOKEN_ERROR;
    return token;
  }
}

export async function handleSignIn(token: JWT, user: User): Promise<JWT> {
  try {
    const userInfo = await authService.getUserInfo(user.access_token);

    token.id = userInfo.id;
    token.username = userInfo.username;
    token.email = userInfo.email;
    token.refresh_token = user.refresh_token;
    token.access_token = user.access_token;
    token.expire_at = user.expire_at;

    return token;
  } catch (error) {
    // Pass error to token in order to it be caught in the client side
    console.log(`[Error] Could not get user info: ${error}`);
    token.error = USER_INFO_ERROR;
    return token;
  }
}
