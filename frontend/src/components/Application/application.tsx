"use client";

import { getSecretGuestAPI } from "@/api/api";
import { DocsApplicationResponse } from "@/api/model";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { getSession } from "next-auth/react";
import { redirect } from "next/navigation";
import { useEffect, useState } from "react";
import Loader from "../common/Loader/loader";
import ApplicationCard from "../common/ApplicationCard/application-card";

const { getApplicationId } = getSecretGuestAPI();

interface ApplicationProps {
  id: string;
}

export default function Application({ id }: ApplicationProps) {
  const [app, setApp] = useState<DocsApplicationResponse | null>(null);

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");

      const resp = await getApplicationId(id, {
        headers: withAuthHeader(session),
      });

      setApp(resp.data);
    })();
  }, []);

  if (!app) return <Loader text="Загружаем вашу заявку" />;

  return (
    <div className="w-full h-full box-border px-2">
      <div className="max-w-1/2 mx-auto box-border pt-20">
        <ApplicationCard {...app} />
      </div>
    </div>
  );
}
