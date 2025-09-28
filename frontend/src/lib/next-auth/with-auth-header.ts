import { RawAxiosRequestHeaders } from "axios";
import { Session } from "next-auth";

export const withAuthHeader = (session: Session): RawAxiosRequestHeaders => ({
  Authorization: `Bearer ${session.accessToken ?? ""}`,
});
