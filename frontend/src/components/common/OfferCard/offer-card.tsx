import { DocsOfferResponse } from "@/api/model";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { formatDateTime } from "@/lib/helpers";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function OfferCard({
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
