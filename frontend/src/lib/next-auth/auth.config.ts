import NextAuth, { NextAuthOptions } from "next-auth";
import Credentials from "next-auth/providers/credentials";
import type { RequestInternal } from "next-auth";
import { resolve } from "path";
import { JWT } from "next-auth/jwt";
import type { NextMiddlewareWithAuth } from "next-auth/middleware";
import { error } from "console";
import { getSecretGuestAPI } from "@/api/api";

const { postUserRefresh } = getSecretGuestAPI();

const refresh = async (old: JWT): Promise<JWT> => {
  try {
    const response = await postUserRefresh({ refresh_token: old.refreshToken ?? "" });

    return {
      ...old,
      accessToken: response.data.access_token,
      refreshToken: response.data.refresh_token,
      accessExpires: Date.now() + (response.data.access_ttl ?? 0) * 1000,
      refreshExpires:
        Date.now() + (response.data.refresh_ttl ?? 0) * 1000,
      error: false,
    };
  } catch (err) {
    return {
      ...old,
      error: true,
    };
  }
};

const providers = [
  Credentials({
    id: "credentials",
    name: "Credentials",
    credentials: {
      accessToken: { label: "accessToken", type: "text" },
      refreshToken: { label: "refreshToken", type: "text" },
      accessTTL: { label: "accessTTL", type: "text" },
      refreshTTL: { label: "refreshTTL", type: "text" },
    },
    async authorize(
      credentials:
        | Record<
            "accessToken" | "refreshToken" | "accessTTL" | "refreshTTL",
            string
          >
        | undefined,
      _req: Pick<RequestInternal, "body" | "query" | "headers" | "method">
    ) {
      if (!credentials) return null;
      return {
        id: "",
        accessToken: credentials.accessToken,
        refreshToken: credentials.refreshToken,
        accessTTL: parseInt(credentials.accessTTL),
        refreshTTL: parseInt(credentials.refreshTTL),
      };
    },
  }),
];

export const authConfig = {
  debug: process.env.NODE_ENV === "development",
  secret: process.env.NEXTAUTH_SECRET,
  session: { strategy: "jwt" },
  providers: providers,

  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.id = user.id;
        token.accessToken = user.accessToken;
        token.refreshToken = user.refreshToken;
        token.accessExpires = Date.now() + user.accessTTL * 1000;
        token.refreshExpires = Date.now() + user.refreshTTL * 1000;
        token.error = false;
        return token;
      }

      if (Date.now() < (token.accessExpires as number) - 5000) return token;
      if (Date.now() < (token.refreshExpires as number) - 5000)
        return await refresh(token);
      return { ...token, error: true };
    },
    async session({ session, token }) {
      session.accessToken = token.accessToken;
      session.refreshToken = token.refreshToken;
      session.error = token.error ?? true;

      return session;
    },
  },
} satisfies NextAuthOptions;

const { auth: middlewareAuth } = NextAuth(authConfig);
export const auth = middlewareAuth as NextMiddlewareWithAuth;
