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

const { getLocation, getOfferSearch } = getSecretGuestAPI();

function Offer({ hotel_id }: DocsOfferResponse) {
  return (
    <div className="w-full box-border rounded-lg border p-4">{hotel_id}</div>
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

  useEffect(() => {
    if (location === null) return;

    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");
      console.log(location);
      const resp = await getOfferSearch(
        { pageNum: pageNum, pageSize: 10, cityId: location.id ?? "" },
        { headers: withAuthHeader(session) }
      );

      setPagesCount(resp.data.pages_count ?? 0);
      setOffers(resp.data.offers ?? []);
    })();
  }, [pageNum, location]);

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

          <Button>
            <Search />
          </Button>
        </div>
      </div>

      <div className="max-w-1/2 mx-auto box-border pt-20">
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
