import NextAuth, { DefaultSession } from "next-auth";

declare module "next-auth" {
  interface Session extends DefaultSession {
    accessToken?: string;
    refreshToken?: string;
    error: boolean;
  }

  interface User {
    accessToken: string;
    refreshToken: string;
    accessTTL: number;
    refreshTTL: number;
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    accessToken?: string;
    refreshToken?: string;
    accessExpires?: number;
    refreshExpires?: number;
    error?: boolean;
  }
}
