import { getSecretGuestAPI } from "@/api/api";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { getSession } from "next-auth/react";

const { getReportMyId } = getSecretGuestAPI();

export interface ReportInfo {
  hotelName: string;
  locationName: string;
  task: string;

  text: string;
  images: {
    id: string;
    url: string;
  }[];

  status: string;
  promocode: string;
}

export const loadReport = async (id: string): Promise<ReportInfo | null> => {
  const session = await getSession();

  if (!session) return null;

  const resp = await getReportMyId(id, {
    headers: withAuthHeader(session),
  });

  return {
    hotelName: resp.data.hotel_name ?? "",
    locationName: resp.data.location_name ?? "",
    task: resp.data.task ?? "",

    images:
      resp.data.images?.map((e) => ({ id: e.id ?? "", url: e.link ?? "" })) ??
      [],
    text: resp.data.text ?? "",

    status: resp.data.status ?? "",
    promocode: resp.data.promocode ?? "",
  };
};
