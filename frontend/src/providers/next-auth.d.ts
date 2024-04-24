import "next-auth";

declare module "next-auth" {
  interface Session {
    user: {
      id: string;
      username: string;
      email: string;
    };
    access_token: string;
    error?: string;
  }

  interface User {
    id?: string;
    username?: string;
    email?: string;
    access_token: string;
    refresh_token: string;
    expire_at: string;
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    id: string;
    username: string;
    email: string;
    access_token: string;
    refresh_token: string;
    expire_at: string;
    error?: string;
  }
}
