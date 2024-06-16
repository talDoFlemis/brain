import authService from "@/services/auth-service";
import { NEXT_AUTH_SECRET } from "@/utils/constants";
import {
  handleSignIn,
  handleTokenRefreshment,
  isTokenExpired,
} from "@/utils/token";
import { isAxiosError } from "axios";
import NextAuth, { NextAuthConfig } from "next-auth";
import Credentials from "next-auth/providers/credentials";

export const authConfig: NextAuthConfig = {
  providers: [
    Credentials({
      credentials: {
        username: {},
        password: {},
      },
      async authorize(credentials) {
        if (!credentials) return null;

        try {
          const response = await authService.signIn({
            username: credentials.username as string,
            password: credentials.password as string,
          });

          return response;
        } catch (error) {
          if (isAxiosError(error)) {
            return Promise.reject(new Error(error.response?.data));
          }
          return Promise.reject(new Error(`${error}`));
        }
      },
    }),
  ],
  session: { strategy: "jwt" },
  pages: {
    signIn: "/sign-in",
    signOut: "/sign-in",
    newUser: "/sign-up",
  },
  secret: NEXT_AUTH_SECRET,
  callbacks: {
    async jwt({ token, user }) {
      // User just signed in
      if (user) {
        return await handleSignIn(token, user);
      }
      // Refresh access token
      if (token.expire_at && isTokenExpired(token.expire_at)) {
        return await handleTokenRefreshment(token);
      }

      return token;
    },
    async session({ session, token }) {
      // Set session data using token data
      session.access_token = token.access_token;
      session.error = token.error;
      session.user.id = token.id;
      session.user.username = token.username;
      session.user.email = token.email;
      return session;
    },
  },
};

export const { auth, handlers, signIn, signOut } = NextAuth(authConfig);
