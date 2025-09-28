import { getSecretGuestAPI } from "@/api/api";
import { DocsApplicationResponse } from "@/api/model";
import { Pagination, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from "@/components/ui/pagination";
import { formatDateTime } from "@/lib/helpers";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { getSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react"

const { getApplication } = getSecretGuestAPI();

const STATUS_MAP = new Map<string, string>([
    ["__app_created", "создана"],
    ["__app_accepted", "принята"],
    ["__app_declined", "отклонена"],
]);

function Application({ hotel_name, expiration_at, status }: DocsApplicationResponse) {
    const statusName = STATUS_MAP.get(status ?? "");

    return <div className="w-full box-border rounded-lg border p-4">
        <div className="font-gain font-bold text-3xl mb-5">{hotel_name}</div>
        <div className="w-full flex justify-between">
            <div className="text-foreground-muted">Дата розыгрыша: {formatDateTime(expiration_at ?? "", 0)}</div>
            <div className="font-bold">{statusName}</div>
        </div>
    </div>
}

export default function Applications() {
    const [pageNum, setPageNum] = useState(0);
    const [pagesCount, setPagesCount] = useState(0);
    const [apps, setApps] = useState<DocsApplicationResponse[]>([]);
    const router = useRouter();

    useEffect(() => {
        (async () => {
            const session = await getSession();

            if (!session) return router.replace("/log-in");
            
            const resp = await getApplication({
                pageNum: pageNum,
                pageSize: 10
            }, {
                headers: withAuthHeader(session)
            });

            setPagesCount(resp.data.pages_count ?? 0);
            setApps(resp.data.applications ?? []);
        })();
    }, [pageNum]);

    return <div className="w-full h-full">
        <div className="font-gain text-2xl mb-5">Заявки</div>

        <div className="flex flex-col gap-4 mb-5">
            {apps.map((e, i) => <Application key={i} {...e}/>)}
        </div>

        <Pagination>
            {pageNum > 0 && <PaginationItem>
                <PaginationPrevious onClick={() => setPageNum(pageNum - 1)}/>
            </PaginationItem>}

        {Array.from({ length: pagesCount }).map((_, i) => <PaginationItem key={i}>
            <PaginationLink onClick={() => setPageNum(i)}>{i+1}</PaginationLink>
        </PaginationItem>)}

            {pageNum < pagesCount - 1 && <PaginationItem>
                <PaginationNext onClick={() => setPageNum(pageNum + 1)}/>
            </PaginationItem>}
        </Pagination>
    </div>
}