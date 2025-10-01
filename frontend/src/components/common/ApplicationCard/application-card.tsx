import { getSecretGuestAPI } from "@/api/api";
import { DocsApplicationResponse } from "@/api/model";
import { Button } from "@/components/ui/button";
import { formatDateTime } from "@/lib/helpers";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { cn } from "@/lib/utils";
import { getSession } from "next-auth/react";
import { redirect, useRouter } from "next/navigation";
import { useMemo } from "react";

const { getReportMyApplicationId } = getSecretGuestAPI();

export const STATUS_MAP = new Map<string, string>([
  ["__app_created", "создана"],
  ["__app_accepted", "принята"],
  ["__app_declined", "отклонена"],
]);

interface ApplicationCardProps extends DocsApplicationResponse {
  showViewReport?: boolean;
}

export default function ApplicationCard({
  id,
  hotel_name,
  expiration_at,
  status,
  showViewReport = true,
}: ApplicationCardProps) {
  const statusName = STATUS_MAP.get(status ?? "");

  const statusCol = useMemo(() => {
    if (status === "__app_accepted") return "text-green-500";
    if (status === "__app_declined") return "text-destructive";
    return "";
  }, [status]);

  const router = useRouter();

  const handleRedirect = async () => {
    const session = await getSession();

    if (!session) return redirect("/log-in");

    const resp = await getReportMyApplicationId(id ?? "", {
      headers: withAuthHeader(session),
    });

    const reportId = resp.data.id;

    router.replace(`/report/${reportId}/view`);
  };

  return (
    <div className="w-full box-border rounded-lg border p-4">
      <div className="font-gain font-bold text-3xl mb-5">{hotel_name}</div>
      <div
        className={cn(
          "w-full flex justify-between",
          status === "__app_accepted" && showViewReport && "mb-9"
        )}
      >
        <div className="text-foreground-muted">
          Дата розыгрыша: {formatDateTime(expiration_at ?? "", 0)}
        </div>
        <div className={cn("font-bold", statusCol)}>{statusName}</div>
      </div>
      {status === "__app_accepted" && showViewReport && (
        <div className="w-full flex justify-start">
          <Button onClick={handleRedirect}>перейти к отчету</Button>
        </div>
      )}
    </div>
  );
}
