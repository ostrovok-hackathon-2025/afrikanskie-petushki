import { getSecretGuestAPI } from "@/api/api";
import { DocsApplicationResponse } from "@/api/model";
import ApplicationCard from "@/components/common/ApplicationCard/application-card";
import { Button } from "@/components/ui/button";
import {
  Pagination,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { formatDateTime } from "@/lib/helpers";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { cn } from "@/lib/utils";
import { getSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";

const { getApplication } = getSecretGuestAPI();

export default function Applications() {
  const [pageNum, setPageNum] = useState(0);
  const [pagesCount, setPagesCount] = useState(0);
  const [apps, setApps] = useState<DocsApplicationResponse[]>([]);
  const router = useRouter();

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return router.replace("/log-in");

      const resp = await getApplication(
        {
          pageNum: pageNum,
          pageSize: 10,
        },
        {
          headers: withAuthHeader(session),
        }
      );

      setPagesCount(resp.data.pages_count ?? 0);
      setApps(resp.data.applications ?? []);
    })();
  }, [pageNum]);

  return (
    <div className="w-full h-full">
      <div className="font-gain text-2xl mb-5">Заявки</div>

      <div className="flex flex-col gap-4 mb-5">
        {apps.map((e, i) => (
          <ApplicationCard key={i} {...e} />
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
