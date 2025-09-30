import { getSecretGuestAPI } from "@/api/api";
import { DocsOfferResponse } from "@/api/model";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { getSession } from "next-auth/react";
import { redirect, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import {
  Pagination,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { cn } from "@/lib/utils";
import { formatDateTime } from "@/lib/helpers";
import { Button } from "@/components/ui/button";
import Link from "next/link";

const { getOffer } = getSecretGuestAPI();

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
  const isOpen = (participants_count ?? 0) < (participants_limit ?? 0);

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

      <div className="w-full flex items-center justify-between mb-10">
        <div className="font-gain">
          {isOpen ? (
            <>
              Места: {participants_count}/{participants_limit}
            </>
          ) : (
            <>Места на розыгрыш закончились</>
          )}
        </div>
      </div>

      <Button asChild>
        <Link href={`/edit-offer/${id}`}>Редактировать</Link>
      </Button>
    </div>
  );
}

export default function AllOffers() {
  const [offers, setOffers] = useState<DocsOfferResponse[]>([]);
  const [pageNum, setPageNum] = useState(0);
  const [pagesCount, setPagesCount] = useState(0);

  const search = () => {
    if (location === null) return;

    (async () => {
      const session = await getSession();

      if (!session) return redirect("/log-in");

      const resp = await getOffer(
        { pageNum: pageNum, pageSize: 5 },
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
    <div className="w-full h-full">
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
  );
}
