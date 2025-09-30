import Application from "@/components/Application/application";
import Offer from "@/components/Offer/offer";
import { isAdmin } from "@/lib/helpers";
import { authConfig } from "@/lib/next-auth/auth.config";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";

interface OfferPageProps {
  params: Promise<{ id: string }>;
}

export default async function OfferPage({ params }: OfferPageProps) {
  const session = await getServerSession(authConfig);
  if (!isAdmin(session)) return redirect("/log-in");

  const { id } = await params;
  return <Offer id={id} />;
}
