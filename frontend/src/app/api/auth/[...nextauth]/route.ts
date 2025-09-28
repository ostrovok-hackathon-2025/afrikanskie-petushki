import { authConfig } from "@/lib/next-auth/auth.config";
import NextAuth from "next-auth";

export const runtime = "nodejs";

const handler = NextAuth(authConfig);
export { handler as GET, handler as POST };
