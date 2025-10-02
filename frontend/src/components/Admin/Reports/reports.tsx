import { getSecretGuestAPI } from "@/api/api";
import {
  DocsHotelResponse,
  DocsLocationResponse,
  DocsOfferResponse,
  DocsReportResponse,
  GetReportSearchParams,
} from "@/api/model";
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
import { Popover, PopoverContent, PopoverTrigger } from "../../ui/popover";
import {
  Command,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "../../ui/command";

const { getReportSearch, patchReportIdConfirm, getHotel, getLocation } =
  getSecretGuestAPI();

const STATUS_MAP = new Map<string, string>([
  ["created", "создан"],
  ["filled", "заполнен"],
  ["accepted", "принят"],
  ["declined", "отклонен"],
]);

function ReportCard({
  id,
  text,
  status,
  images,
  check_in_at,
  check_out_at,
  room_name,
  hotel_name,
  task,
}: DocsReportResponse) {
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
                <div
                  key={i}
                  className="h-[200px] relative overflow-hidden rounded-lg"
                >
                  <img
                    src={e.link}
                    alt="image"
                    className="absolute top-1/2 left-1/2 cover rounded-lg -translate-x-1/2 -translate-y-1/2"
                  />
                </div>
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
  const [hotels, setHotels] = useState<DocsHotelResponse[]>([]);
  const [hotel, setHotel] = useState<DocsHotelResponse | null>(null);

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");

      const resp = await getHotel({ headers: withAuthHeader(session) });

      setHotels(resp.data.hotels ?? []);
    })();
  }, []);
  const [locations, setLocations] = useState<DocsLocationResponse[]>([]);
  const [location, setLocation] = useState<DocsLocationResponse | null>(null);

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return;

      const resp = await getLocation({ headers: withAuthHeader(session) });

      setLocations(resp.data.locations ?? []);
    })();
  }, []);
  const [status, setStatus] = useState<string | null>(null);

  const [reports, setReports] = useState<DocsReportResponse[]>([]);
  const [pageNum, setPageNum] = useState(0);
  const [pagesCount, setPagesCount] = useState(0);

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");

      const filters: GetReportSearchParams = {
        pageNum: pageNum,
        pageSize: 3,
      };

      if (location) filters.cityId = location.id;
      if (hotel) filters.hotelId = hotel.id;
      if (status) filters.status = status;

      console.log(filters);

      const resp = await getReportSearch(filters, {
        headers: withAuthHeader(session),
      });

      setPagesCount(resp.data.pages_count ?? 0);
      setReports(resp.data.reports ?? []);
    })();
  }, [pageNum, location, hotel, status]);

  return (
    <div className="w-full h-full">
      <div className="flex w-full mb-4 gap-4">
        <div className="w-full rounded-lg bg-accent box-border p-2 flex items-center justify-between">
          <div className="font-gain">Отель: {hotel?.name ?? "не выбран"}</div>

          <Popover>
            <PopoverTrigger asChild className="flex-shrink-0">
              <Button onClick={() => setHotel(null)}>выбрать</Button>
            </PopoverTrigger>

            <PopoverContent>
              <Command>
                <CommandInput placeholder="введите название отеля..."></CommandInput>

                <CommandList>
                  <CommandGroup>
                    {hotels.map((e, i) => (
                      <CommandItem
                        key={e.id}
                        value={e.name}
                        onSelect={() => setHotel(e)}
                      >
                        {e.name}
                      </CommandItem>
                    ))}
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>
        </div>
      </div>

      <div className="flex w-full mb-12 gap-4">
        <div className="w-1/2 rounded-lg bg-accent box-border p-2 flex items-center justify-between">
          <div className="font-gain">
            Город: {location?.name ?? "не выбран"}
          </div>

          <Popover>
            <PopoverTrigger asChild className="flex-shrink-0">
              <Button onClick={() => setLocation(null)}>выбрать</Button>
            </PopoverTrigger>

            <PopoverContent>
              <Command>
                <CommandInput placeholder="введите город..."></CommandInput>

                <CommandList>
                  <CommandGroup>
                    {locations.map((e, i) => (
                      <CommandItem
                        key={e.id}
                        value={e.name}
                        onSelect={() => setLocation(e)}
                      >
                        {e.name}
                      </CommandItem>
                    ))}
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>
        </div>

        <div className="w-1/2 rounded-lg bg-accent box-border p-2 flex items-center justify-between">
          <div className="font-gain">
            Статус: {STATUS_MAP.get(status ?? "") ?? "не выбран"}
          </div>

          <Popover>
            <PopoverTrigger asChild className="flex-shrink-0">
              <Button onClick={() => setStatus(null)}>выбрать</Button>
            </PopoverTrigger>

            <PopoverContent>
              <Command>
                <CommandInput placeholder="введите город..."></CommandInput>

                <CommandList>
                  <CommandGroup>
                    {Array.from(
                      STATUS_MAP.keys().map((e, i) => (
                        <CommandItem
                          key={e}
                          value={e}
                          onSelect={() => setStatus(e)}
                        >
                          {STATUS_MAP.get(e)}
                        </CommandItem>
                      ))
                    )}
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>
        </div>
      </div>

      <div className="flex flex-col gap-4 mb-5">
        {reports.map((e, i) => (
          <ReportCard key={e.id} {...e} />
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
