import Home from "@/components/Home/home";
import { isAdmin } from "@/lib/helpers";
import { authConfig } from "@/lib/next-auth/auth.config";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";

interface HomePageProps {
  params: Promise<{ page: string }>;
}

export default async function HomePage({ params }: HomePageProps) {
  const { page } = await params;

  const session = await getServerSession(authConfig);
  if (!session) return redirect("/log-in");
  if (isAdmin(session)) return redirect("/admin/offers");

  return <Home page={page} />;
}
