import { Session } from "next-auth";
import { jwtDecode } from "jwt-decode";

export function formatDateTime(isoString: string, timeShift: number) {
  const months = [
    "января",
    "февраля",
    "марта",
    "апреля",
    "мая",
    "июня",
    "июля",
    "августа",
    "сентября",
    "октября",
    "ноября",
    "декабря",
  ];

  const now = new Date();
  const inputDate = new Date(isoString);

  inputDate.setHours(inputDate.getHours() + timeShift);

  const isToday = now.toDateString() === inputDate.toDateString();

  let dateString;
  if (isToday) {
    dateString = "Сегодня";
  } else {
    dateString = `${inputDate.getDate()} ${months[inputDate.getMonth()]}`;
  }

  const timeString = inputDate.toLocaleTimeString("ru-RU", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });

  return `${dateString}, ${timeString}`;
}

export function isAdmin(session: Session | null) {
  if (!session || !session.accessToken) return false;

  const claims: { is_admin: boolean } = jwtDecode(session.accessToken);
  return claims.is_admin ?? false;
}

export function toRGC3339(time: string) {
  const date = new Date(time.replace(" ", "T") + "Z");
  return date.toISOString();
}
