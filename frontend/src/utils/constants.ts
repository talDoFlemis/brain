export const BACKEND_API_URL =
  process.env.NEXT_PUBLIC_BACKEND_API_URL || "http://localhost:42069";
export const SIGN_OUT_CALLBACK_URL = "/sign-in";
export const SIGN_IN_CALLBACK_URL = "/dashboard";
export const NEXT_AUTH_SECRET = process.env.NEXTAUTH_SECRET || "top-secret";
