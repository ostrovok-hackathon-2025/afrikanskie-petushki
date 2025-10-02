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
import { cn } from "@/lib/utils";
import OfferCard from "@/components/common/OfferCard/offer-card";

const { getOffer } = getSecretGuestAPI();

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
          <OfferCard key={i} {...e} />
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
