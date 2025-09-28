import Home from "@/components/Home/home";
import { authConfig } from "@/lib/next-auth/auth.config";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";

export default async function HomePage() {
  const session = await getServerSession(authConfig);

  if (!session) return redirect("log-in");

  return <Home />;
}
