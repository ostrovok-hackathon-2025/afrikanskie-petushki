import { DocsApplicationResponse } from "@/api/model";
import { Button } from "@/components/ui/button";
import { formatDateTime } from "@/lib/helpers";
import { cn } from "@/lib/utils";
import { useMemo } from "react";

export const STATUS_MAP = new Map<string, string>([
  ["__app_created", "создана"],
  ["__app_accepted", "принята"],
  ["__app_declined", "отклонена"],
]);

interface ApplicationCardProps extends DocsApplicationResponse {
  showViewReport?: boolean;
}

export default function ApplicationCard({
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
          <Button>перейти к отчету</Button>
        </div>
      )}
    </div>
  );
}
