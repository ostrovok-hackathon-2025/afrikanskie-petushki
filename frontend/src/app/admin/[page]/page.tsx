import Admin from "@/components/Admin/admin";
import { isAdmin } from "@/lib/helpers";
import { authConfig } from "@/lib/next-auth/auth.config";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";

interface HomePageProps {
  params: Promise<{ page: string }>;
}

export default async function AdminPage({ params }: HomePageProps) {
  const session = await getServerSession(authConfig);
  if (!isAdmin(session)) return redirect("/log-in");

  const { page } = await params;

  return <Admin page={page} />;
}
