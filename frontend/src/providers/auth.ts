import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";
import authService from "@/services/auth-service";
import { handleSignIn, handleTokenRefreshment, isTokenExpired } from "@/utils/token";

export const handler = NextAuth({
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
            username: credentials.username,
            password: credentials.password,
          });
          
          return response;
        } catch (error) {
          console.log(`[Error] Could not sign in: ${error}`);
          return null;
        }
      },
    }),
  ],
  session: { strategy: "jwt" },
  pages: {
    signIn: "/sign-in",
    signOut: "/sign-in",
  },
  secret: process.env.NEXTAUTH_SECRET,
  callbacks: {
    async jwt({ token, user }) {
      // User just signed in
      if (user) {
        return await handleSignIn(token, user);
      }
      // Refresh access token
      if (token.expire_at && isTokenExpired(token.expire_at)) {
        return await handleTokenRefreshment(token)
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
});
