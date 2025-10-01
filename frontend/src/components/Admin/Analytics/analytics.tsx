import { getSecretGuestAPI } from "@/api/api";
import { DocsAnalyticsResponse } from "@/api/model";
import CountUp from "@/components/CountUp";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { getSession } from "next-auth/react";
import { redirect } from "next/navigation";
import { useEffect, useState } from "react";

const { getAnalytics } = getSecretGuestAPI();

interface AnalyticsCardProps {
  text: string;
  value: number;
}

function AnalyticsCard({ text, value }: AnalyticsCardProps) {
  const avr = value / new Date().getDate();

  return (
    <div className="w-full h-[400px] box-border rounded-lg border p-4 flex flex-col justify-center items-center">
      <div className="font-medium text-xl mb-8 text-center min-h-14">
        {text}
      </div>
      <CountUp
        from={0}
        to={value}
        direction="up"
        duration={1}
        className="text-4xl font-bold mb-12"
      />

      <div className="font-sm text-center">за последний месяц</div>
      <div className="font-sm text-center">
        в сред. <span className="font-medium">{avr.toFixed(2)}/день</span>
      </div>
    </div>
  );
}

export default function Analytics() {
  const [analytics, setAnalytics] = useState<DocsAnalyticsResponse | null>();

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");

      const resp = await getAnalytics({
        headers: withAuthHeader(session),
      });

      setAnalytics(resp.data);
    })();
  }, []);

  return (
    <div className="w-full">
      <div className="w-full font-bold text-3xl text-center mb-16 mt-20">
        Общая аналитика
      </div>
      <div className="flex w-full gap-8">
        <AnalyticsCard
          text="Cобрано заявок "
          value={analytics?.applications_received ?? 0}
        />
        <AnalyticsCard
          text="Завершено розыгрышей"
          value={analytics?.completed_offers ?? 0}
        />
        <AnalyticsCard
          text="Принято отчетов"
          value={analytics?.accepted_reports ?? 0}
        />
      </div>
    </div>
  );
}
