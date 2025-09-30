"use client";

import { getSecretGuestAPI } from "@/api/api";
import { DocsLocationResponse, DocsOfferResponse } from "@/api/model";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { getSession } from "next-auth/react";
import { redirect } from "next/navigation";
import { useEffect, useState } from "react";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { Button } from "../ui/button";
import { Command, CommandInput, CommandItem, CommandList } from "../ui/command";
import { CommandGroup } from "cmdk";
import { Search } from "lucide-react";
import {
  Pagination,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { cn } from "@/lib/utils";
import { formatDateTime } from "@/lib/helpers";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { useRouter } from "next/navigation";

const { getLocation, getOfferSearch, postApplication } = getSecretGuestAPI();

const lerp = (a: number, b: number, c: number) => a + (b - a) * c;

function Offer({
  id,
  hotel_name,
  room_name,
  check_in_at,
  check_out_at,
  task,
  participants_count,
  participants_limit,
}: DocsOfferResponse) {
  const router = useRouter();

  const participateHandler = async () => {
    const session = await getSession();

    if (!session) return redirect("/log-in");

    const resp = await postApplication(
      {
        offer_id: id ?? "",
      },
      {
        headers: withAuthHeader(session),
      }
    );

    const applicationId = resp.data.application_id;

    router.replace(`/application/${applicationId}`);
  };

  const isOpen = (participants_count ?? 0) < (participants_limit ?? 0);
  const rest = (participants_limit ?? 0) - (participants_count ?? 0);
  const fill = rest / (participants_limit ?? 1);

  const color = `rgb(${[
    lerp(233, 0, fill),
    lerp(0, 201, fill),
    lerp(11, 80, fill),
  ]
    .map(Math.floor)
    .join(", ")})`;

  return (
    <div className="w-full box-border rounded-lg border p-4">
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
          <AccordionTrigger className="font-gain font-medium text-xl items-center mb-2 no-underline">
            Задание
          </AccordionTrigger>
          <AccordionContent className="font-gain whitespace-pre-wrap">
            {task}
          </AccordionContent>
        </AccordionItem>
      </Accordion>

      <div className="w-full flex items-center justify-between">
        <div className="font-gain">
          {isOpen ? (
            <>
              Осталось мест на розыгрыш{" "}
              <span className="font-bold" style={{ color }}>
                {rest}
              </span>{" "}
            </>
          ) : (
            <>Места на розыгрыш закончились</>
          )}
        </div>

        {isOpen ? (
          <Button onClick={() => participateHandler()}>участвовать</Button>
        ) : (
          <div className="rounded-sm bg-accent box-border p-2">
            набор закрыт
          </div>
        )}
      </div>
    </div>
  );
}

export default function Offers() {
  const [locations, setLocations] = useState<DocsLocationResponse[]>([]);
  const [location, setLocation] = useState<DocsLocationResponse | null>(null);

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");

      const resp = await getLocation({ headers: withAuthHeader(session) });

      setLocations(resp.data.locations ?? []);
    })();
  }, []);

  const [offers, setOffers] = useState<DocsOfferResponse[]>([]);
  const [pageNum, setPageNum] = useState(0);
  const [pagesCount, setPagesCount] = useState(0);

  const search = () => {
    if (location === null) return;

    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");
      console.log(location);
      const resp = await getOfferSearch(
        { pageNum: pageNum, pageSize: 5, cityId: location.id ?? "" },
        { headers: withAuthHeader(session) }
      );

      setPagesCount(resp.data.pages_count ?? 0);
      setOffers(resp.data.offers ?? []);
    })();
  };

  useEffect(() => {
    search();
  }, [pageNum]);

  return (
    <div className="w-full h-full box-border px-2">
      <div className="max-w-1/2 w-full box-border fixed top-20 left-1/2 -translate-x-1/2">
        <div className="w-full flex gap-2 items-center">
          <Popover>
            <PopoverTrigger asChild className="flex-shrink-0">
              <Button>выберите город</Button>
            </PopoverTrigger>

            <PopoverContent>
              <Command>
                <CommandInput placeholder="введите название города..."></CommandInput>

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

          <div className="bg-accent h-9 rounded-sm flex-grow flex items-center box-border px-2">
            Город: {location?.name ?? "не выбран"}
          </div>

          <Button onClick={search}>
            <Search />
          </Button>
        </div>
      </div>

      <div className="max-w-1/2 mx-auto box-border pt-34">
        <div className="flex flex-col gap-4 mb-5">
          {offers.map((e, i) => (
            <Offer key={i} {...e} />
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
    </div>
  );
}
