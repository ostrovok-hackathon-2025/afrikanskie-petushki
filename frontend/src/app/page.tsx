import Home from "@/components/Home/home";
import Landing from "@/components/Landing/landing";
import { authConfig } from "@/lib/next-auth/auth.config";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";

export default async function HomePage() {
  return <Landing />;
}
