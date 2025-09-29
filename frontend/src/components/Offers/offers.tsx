"use client";

import { getSecretGuestAPI } from "@/api/api";
import { DocsLocationResponse } from "@/api/model";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { getSession } from "next-auth/react";
import { redirect } from "next/navigation";
import { useEffect, useState } from "react";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { Button } from "../ui/button";
import { Command, CommandInput, CommandItem, CommandList } from "../ui/command";
import { CommandGroup } from "cmdk";
import { Search } from "lucide-react";

const { getLocation } = getSecretGuestAPI();

export default function Offers() {
  const [locations, setLocations] = useState<DocsLocationResponse[]>([
    { id: "1", name: "Москва" },
    { id: "2", name: "Владимир" },
    { id: "3", name: "Черноголовка" },
  ]);
  const [location, setLocation] = useState<DocsLocationResponse | null>(null);

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");

      const resp = await getLocation({ headers: withAuthHeader(session) });

      setLocations(resp.data.locations ?? []);
    })();
  }, []);

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
    </div>
  );
}
