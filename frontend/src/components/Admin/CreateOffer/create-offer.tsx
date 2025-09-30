import { getSecretGuestAPI } from "@/api/api";
import { DocsHotelResponse, DocsRoomResponse } from "@/api/model";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { cn } from "@/lib/utils";
import { getSession } from "next-auth/react";
import { redirect, useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { Popover, PopoverContent, PopoverTrigger } from "../../ui/popover";
import {
  Command,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "../../ui/command";
import { Button } from "@/components/ui/button";
import { DateTimeSelect } from "./DateTimeSelect/date-time-select";
import { addDays, format } from "date-fns";
import { Input } from "@/components/ui/input";
import { CircleMinus, CirclePlus } from "lucide-react";
import { formatDateTime, toRGC3339 } from "@/lib/helpers";

const { getHotel, getRoom, postOffer } = getSecretGuestAPI();

export default function CreateOffer() {
  const [task, setTask] = useState("");

  const textareaRef = useRef<HTMLTextAreaElement | null>(null);

  const autoResize = () => {
    if (!textareaRef.current) return;
    const textarea = textareaRef.current;
    textarea.style.height = "auto";
    textarea.style.height = textarea.scrollHeight + 10 + "px";
  };

  useEffect(() => {
    autoResize();
  }, [task]);

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

  const [endDate, setEndDate] = useState<string>(
    format(addDays(new Date().setHours(0, 0, 0, 0), 1), "yyyy-MM-dd HH:mm:ss")
  );

  const [checkIn, setCheckIn] = useState<string>(
    format(addDays(new Date().setHours(0, 0, 0, 0), 2), "yyyy-MM-dd HH:mm:ss")
  );

  const [checkOut, setCheckOut] = useState<string>(
    format(addDays(new Date().setHours(0, 0, 0, 0), 3), "yyyy-MM-dd HH:mm:ss")
  );

  useEffect(() => {
    if (addDays(endDate, 1) >= new Date(checkIn)) {
      setCheckIn(format(addDays(new Date(endDate), 1), "yyyy-MM-dd HH:mm:ss"));
    }
  }, [endDate, checkIn]);

  useEffect(() => {
    if (addDays(checkIn, 1) >= new Date(checkOut)) {
      setCheckOut(format(addDays(new Date(checkIn), 1), "yyyy-MM-dd HH:mm:ss"));
    }
  }, [checkIn, checkOut]);

  const [participantsLimit, setParticipantsLimit] = useState(20);

  const router = useRouter();

  const handleSubmit = async () => {
    const session = await getSession();

    if (!session) return redirect("/log-in");

    if (hotel === null) return;
    if (room === null) return;
    if (task.length === 0) return;

    const body = {
      check_in: toRGC3339(checkIn),
      check_out: toRGC3339(checkOut),
      expiration_at: toRGC3339(endDate),
      hotel_id: hotel.id ?? "",
      room_id: room.id ?? "",
      participants_limit: participantsLimit,
      task: task,

      // legacy
      location_id: hotel.id ?? "",
    };

    console.log(body);

    const res = await postOffer(body, {
      headers: withAuthHeader(session),
    });

    router.replace(`/offer/${res.data.id}`);
  };

  return (
    <div className="w-full h-full select-none">
      <div className="text-3xl font-bold mb-20 mt-20">Создание розыгрыша</div>

      <form>
        <div className="font-gain font-medium text-xl mb-4">Дата окончания</div>

        <div className="w-2/3 rounded-lg bg-accent box-border p-2 flex items-center justify-between mb-16">
          <div className="font-gain">Дата: {formatDateTime(endDate, 0)}</div>

          <Popover>
            <PopoverTrigger asChild className="flex-shrink-0">
              <Button>выбрать</Button>
            </PopoverTrigger>

            <PopoverContent className="w-full">
              <DateTimeSelect
                onDateTimeChange={(ts) =>
                  setEndDate(format(ts, "yyyy-MM-dd HH:mm:ss"))
                }
                minDate={new Date()}
              />
            </PopoverContent>
          </Popover>
        </div>

        <div className="font-gain font-medium text-xl mb-4">
          Отель и тип номера
        </div>

        <div className="flex w-full mb-16 gap-4">
          <div className="w-2/3 rounded-lg bg-accent box-border p-2 flex items-center justify-between">
            <div className="font-gain">Отель: {hotel?.name ?? "не выбран"}</div>

            <Popover>
              <PopoverTrigger asChild className="flex-shrink-0">
                <Button>выбрать</Button>
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
                <Button>выбрать</Button>
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

        <div className="font-gain font-medium text-xl mb-4">
          Период проживания
        </div>

        <div className="flex w-full mb-16 gap-4">
          <div className="w-1/2 rounded-lg bg-accent box-border p-2 flex items-center justify-between">
            <div className="font-gain">
              Заселение: {formatDateTime(checkIn, 0)}
            </div>

            <Popover>
              <PopoverTrigger asChild className="flex-shrink-0">
                <Button>выбрать</Button>
              </PopoverTrigger>

              <PopoverContent className="w-full">
                <DateTimeSelect
                  onDateTimeChange={(ts) =>
                    setCheckIn(format(ts, "yyyy-MM-dd HH:mm:ss"))
                  }
                  minDate={new Date(endDate)}
                />
              </PopoverContent>
            </Popover>
          </div>

          <div className="w-1/2 rounded-lg bg-accent box-border p-2 flex items-center justify-between">
            <div className="font-gain">
              Выселение: {formatDateTime(checkOut, 0)}
            </div>

            <Popover>
              <PopoverTrigger asChild className="flex-shrink-0">
                <Button>выбрать</Button>
              </PopoverTrigger>

              <PopoverContent className="w-full">
                <DateTimeSelect
                  onDateTimeChange={(ts) =>
                    setCheckOut(format(ts, "yyyy-MM-dd HH:mm:ss"))
                  }
                  minDate={new Date(checkIn)}
                />
              </PopoverContent>
            </Popover>
          </div>
        </div>

        <div className="font-gain font-medium text-xl mb-4">Задание</div>
        <textarea
          value={task}
          onChange={(e) => setTask(e.target.value)}
          className={cn(
            "outline-none border w-full h-fit rounded-lg border box-border p-2 resize-none mb-16"
          )}
          ref={textareaRef}
          name="text"
        ></textarea>

        <div className="font-gain font-medium text-xl mb-4">
          Количество мест
        </div>
        <div className="flex gap-2 border box-border py-2 px-4 rounded-sm mb-20 w-fit items-center">
          <CircleMinus
            className="cursor-pointer hover:scale-[1.05] transition-all"
            onClick={() =>
              setParticipantsLimit(Math.max(10, participantsLimit - 1))
            }
          />
          <div className="font-bold text-xl select-none">
            {participantsLimit}
          </div>
          <CirclePlus
            className="cursor-pointer hover:scale-[1.05] transition-all"
            onClick={() =>
              setParticipantsLimit(Math.min(100, participantsLimit + 1))
            }
          />
        </div>

        <Button
          size={"lg"}
          className="text-xl select-none"
          type="button"
          onClick={handleSubmit}
        >
          начать розыгрыш
        </Button>
        <div className="h-20"></div>
      </form>
    </div>
  );
}
