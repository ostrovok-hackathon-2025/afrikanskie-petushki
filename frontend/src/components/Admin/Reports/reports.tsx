import { getSecretGuestAPI } from "@/api/api";
import { DocsOfferResponse, DocsReportResponse } from "@/api/model";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { getSession } from "next-auth/react";
import { redirect, useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import {
  Pagination,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { cn } from "@/lib/utils";
import OfferCard from "@/components/common/OfferCard/offer-card";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { formatDateTime } from "@/lib/helpers";
import { Button } from "@/components/ui/button";

const { getReport, patchReportIdConfirm } = getSecretGuestAPI();

function ReportCard({ id, text, status, images }: DocsReportResponse) {
  const room_name = "mock_room";
  const hotel_name = "mock_hotel";
  const check_in_at = "2025-12-22 14:00:00";
  const check_out_at = "2025-12-22 14:00:00";
  const task = "mock_task";

  const [statusText, setStatusText] = useState(status);

  const waiting = useMemo(
    () => statusText === "created" || statusText === "filled",
    [statusText]
  );

  const updateStatus = async (status: string) => {
    setStatusText(status);

    const session = await getSession();

    if (!session) return redirect("/log-in");

    const resp = await patchReportIdConfirm(
      id ?? "",
      {
        status: status,
      },
      {
        headers: withAuthHeader(session),
      }
    );
  };

  return (
    <div className="w-full box-border rounded-lg border p-4">
      <div className="text-2xl font-bold mb-8">Отчет по заданию</div>
      <div className="font-gain text-muted-foreground text-base mb-2">
        Номер {room_name} в отеле
      </div>
      <div className="font-gain font-bold text-3xl mb-4">{hotel_name}</div>

      <div className="font-gain rounded-sm bg-accent box-border p-2 mb-6">
        Заезд:{" "}
        <span className="font-medium">
          {formatDateTime(check_in_at ?? "", 0)}
        </span>{" "}
        - Выезд:{" "}
        <span className="font-medium">
          {formatDateTime(check_out_at ?? "", 0)}
        </span>
      </div>

      <Accordion type="single" collapsible>
        <AccordionItem value="item-1">
          <AccordionTrigger className="font-gain font-medium text-xl items-center mb-2 no-underline py-0">
            Текст задания
          </AccordionTrigger>
          <AccordionContent className="font-gain whitespace-pre-wrap">
            {task}
          </AccordionContent>
        </AccordionItem>
      </Accordion>

      <hr className="my-6"></hr>

      {status !== "created" ? (
        <>
          <Accordion type="single" collapsible className="mb-6">
            <AccordionItem value="item-1">
              <AccordionTrigger className="font-gain font-medium text-xl items-center mb-2 no-underline py-0">
                Текст отчета
              </AccordionTrigger>
              <AccordionContent className="font-gain whitespace-pre-wrap">
                {text}
              </AccordionContent>
            </AccordionItem>
          </Accordion>

          <div className="font-gain font-medium text-xl mb-2">
            {images && images.length > 0
              ? "Изображения"
              : "Пользователь не загрузил изображений"}
          </div>

          <div className="w-full grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-4 mb-8">
            {images &&
              images.map((e, i) => (
                <img
                  key={i}
                  src={e.link}
                  alt="image"
                  className="cover w-full h-[175px] rounded-lg"
                />
              ))}
          </div>
        </>
      ) : (
        <div className="w-full flex justify-center font-bold">
          Пользователь не заполнил отчет
        </div>
      )}

      <hr className="my-6"></hr>

      {waiting && (
        <div className="flex gap-2">
          <Button
            className="bg-green-500 hover:bg-green-500/90"
            onClick={() => updateStatus("accepted")}
          >
            Принять
          </Button>
          <Button
            className="bg-destructive hover:bg-destructive/90"
            onClick={() => updateStatus("declined")}
          >
            Отклонить
          </Button>
        </div>
      )}

      {!waiting && (
        <div className="font-bold">
          Вердикт:{" "}
          <span
            className={
              statusText === "accepted" ? "text-green-500" : "text-destructive"
            }
          >
            {statusText === "accepted" ? "принят" : "отклонен"}
          </span>
        </div>
      )}
    </div>
  );
}

export default function Reports() {
  const [reports, setReports] = useState<DocsReportResponse[]>([]);
  const [pageNum, setPageNum] = useState(0);
  const [pagesCount, setPagesCount] = useState(1);

  const search = () => {
    if (location === null) return;

    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");

      const resp = await getReport(
        { pageNum: pageNum, pageSize: 5 },
        { headers: withAuthHeader(session) }
      );

      setPagesCount(resp.data.pages_count ?? 0);
      setReports(resp.data.reports ?? []);
    })();
  };

  useEffect(() => {
    search();
  }, [pageNum]);

  return (
    <div className="w-full h-full">
      <div className="flex flex-col gap-4 mb-5">
        {reports.map((e, i) => (
          <ReportCard key={i} {...e} />
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
