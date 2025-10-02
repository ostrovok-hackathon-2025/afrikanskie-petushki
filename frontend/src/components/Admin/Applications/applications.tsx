import {
  DocsApplicationResponse,
  DocsHotelResponse,
  DocsLocationResponse,
  DocsRoomResponse,
  GetApplicationSearchParams,
} from "@/api/model";
import { useEffect, useState } from "react";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { cn } from "@/lib/utils";
import { getSession } from "next-auth/react";
import { getSecretGuestAPI } from "@/api/api";
import { redirect, useRouter } from "next/navigation";
import { Popover, PopoverContent, PopoverTrigger } from "../../ui/popover";
import {
  Command,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "../../ui/command";
import { Button } from "@/components/ui/button";
import ApplicationCard, {
  STATUS_MAP,
} from "@/components/common/ApplicationCard/application-card";
import {
  Pagination,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

const { getHotel, getRoom, getLocation, getApplicationSearch } =
  getSecretGuestAPI();

export default function Applications() {
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

  const [rooms, setRooms] = useState<DocsRoomResponse[]>([]);
  const [room, setRoom] = useState<DocsRoomResponse | null>(null);

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");

      const resp = await getRoom({ headers: withAuthHeader(session) });

      setRooms(resp.data.rooms ?? []);
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

  const [pageNum, setPageNum] = useState(0);
  const [pagesCount, setPagesCount] = useState(0);
  const [apps, setApps] = useState<DocsApplicationResponse[]>([]);
  const router = useRouter();

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return router.replace("/log-in");

      const filters: GetApplicationSearchParams = {
        pageNum: pageNum,
        pageSize: 3,
      };

      if (location) filters.cityId = location.id;
      if (hotel) filters.hotelId = hotel.id;
      if (room) filters.roomId = room.id;
      if (status) filters.status = status;

      const resp = await getApplicationSearch(filters, {
        headers: withAuthHeader(session),
      });

      setPagesCount(resp.data.pages_count ?? 0);
      setApps(resp.data.applications ?? []);
    })();
  }, [pageNum, location, hotel, room, status]);

  return (
    <div className="w-full">
      <div className="flex w-full mb-4 gap-4">
        <div className="w-2/3 rounded-lg bg-accent box-border p-2 flex items-center justify-between">
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

        <div className="w-1/3 rounded-lg bg-accent box-border p-2 flex items-center justify-between">
          <div className="font-gain">Номер: {room?.name ?? "не выбран"}</div>

          <Popover>
            <PopoverTrigger asChild className="flex-shrink-0">
              <Button onClick={() => setRoom(null)}>выбрать</Button>
            </PopoverTrigger>

            <PopoverContent>
              <Command>
                <CommandInput placeholder="введите тип номера..."></CommandInput>

                <CommandList>
                  <CommandGroup>
                    {rooms.map((e, i) => (
                      <CommandItem
                        key={e.id}
                        value={e.name}
                        onSelect={() => setRoom(e)}
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
        {apps.map((e, i) => (
          <ApplicationCard showViewReport={false} key={i} {...e} />
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
