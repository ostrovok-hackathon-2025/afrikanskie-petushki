import { getSecretGuestAPI } from "@/api/api";
import { DocsReportResponse } from "@/api/model";
import { Button } from "@/components/ui/button";
import { formatDateTime } from "@/lib/helpers";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { cn } from "@/lib/utils";
import { Eye } from "lucide-react";
import { getSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import {
  Pagination,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

const { getReportMy } = getSecretGuestAPI();

const STATUS_MAP = new Map<string, string>([
  ["created", "создан"],
  ["filled", "заполнен"],
  ["accepted", "принят"],
  ["declined", "отклонен"],
]);

function ReportCard({ status, expiration_at, id }: DocsReportResponse) {
  const statusCol = useMemo(() => {
    if (status === "accepted") return "text-green-500";
    if (status === "declined") return "text-destructive";
    return "";
  }, [status]);

  const router = useRouter();

  return (
    <div className="w-full box-border rounded-lg border p-4">
      <div className="w-full flex justify-between items-center mb-4">
        <div className="font-gain text-lg font-bold">
          {formatDateTime(expiration_at ?? "", 0)}
        </div>
        <div className={cn("font-gain text-lg font-bold", statusCol)}>
          {STATUS_MAP.get(status ?? "")}
        </div>
      </div>
      <Button
        type="button"
        onClick={() => router.replace(`/report/${id}/view`)}
      >
        Смотреть <Eye />
      </Button>
    </div>
  );
}

export default function Reports() {
  const [pageNum, setPageNum] = useState(0);
  const [pagesCount, setPagesCount] = useState(0);
  const [reports, setReports] = useState<DocsReportResponse[]>([]);

  const router = useRouter();

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return router.replace("/log-in");

      const resp = await getReportMy(
        {
          pageNum: pageNum,
          pageSize: 10,
        },
        {
          headers: withAuthHeader(session),
        }
      );

      setPagesCount(resp.data.pages_count ?? 0);
      setReports(resp.data.reports ?? []);
    })();
  }, [pageNum]);

  return (
    <div className="w-full">
      <div className="font-gain text-2xl mb-10">Отчеты</div>

      <div className="flex flex-col gap-6 mb-8">
        {reports.map((e, i) => (
          <ReportCard {...e} key={e.id} />
        ))}
      </div>

      <Pagination>
        {pageNum > 0 && (
          <PaginationItem>
            <PaginationPrevious onClick={() => setPageNum(pageNum - 1)} />
          </PaginationItem>
        )}

        {Array.from({ length: pagesCount }).map((_, i) => (
          <PaginationItem key={i}>
            <PaginationLink
              onClick={() => setPageNum(i)}
              className={cn(
                "hover:text-primary",
                pageNum === i && "text-primary bg-accent"
              )}
            >
              {i + 1}
            </PaginationLink>
          </PaginationItem>
        ))}

        {pageNum < pagesCount - 1 && (
          <PaginationItem>
            <PaginationNext onClick={() => setPageNum(pageNum + 1)} />
          </PaginationItem>
        )}
      </Pagination>
    </div>
  );
}
