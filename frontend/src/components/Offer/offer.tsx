"use client";

import { getSecretGuestAPI } from "@/api/api";
import { DocsOfferResponse } from "@/api/model";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { getSession } from "next-auth/react";
import { redirect } from "next/navigation";
import { useEffect, useState } from "react";
import Loader from "../common/Loader/loader";
import OfferCard from "../common/OfferCard/offer-card";

const { getOfferId } = getSecretGuestAPI();

interface ApplicationProps {
  id: string;
}

export default function Offer({ id }: ApplicationProps) {
  const [offer, setOffer] = useState<DocsOfferResponse | null>(null);

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");

      const resp = await getOfferId(id, {
        headers: withAuthHeader(session),
      });

      setOffer(resp.data);
    })();
  }, []);

  if (!offer) return <Loader text="Загружаем розыгрыш" />;

  return (
    <div className="w-full h-full box-border px-2">
      <div className="max-w-1/2 mx-auto box-border pt-20">
        <OfferCard {...offer} />
      </div>
    </div>
  );
}
