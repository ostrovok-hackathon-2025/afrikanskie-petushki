"use client";
import { getSecretGuestAPI } from "@/api/api";
import { Button } from "@/components/ui/button";
import { isAdmin } from "@/lib/helpers";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { Search } from "lucide-react";
import { getSession, signOut } from "next-auth/react";
import Image from "next/image";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

const { getUser } = getSecretGuestAPI();

export default function Header() {
  const [username, setUsername] = useState("");
  const router = useRouter();

  const [isAdminAccount, setIsAdmin] = useState(false);

  const pathname = usePathname();

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) {
        setUsername("$empty");
        return;
      }

      setIsAdmin(isAdmin(session));

      const resp = await getUser({ headers: withAuthHeader(session) });

      setUsername(resp.data.ostrovok_login ?? "");
    })();
  }, [pathname]);

  const handleRedirect = useCallback(() => {
    if (username === "$empty") {
      router.replace("/log-in");
    } else {
      router.replace("/home/profile");
    }
  }, [username]);

  const handleExit = useCallback(async () => {
    await signOut();
    router.replace("/log-in");
  }, []);

  return (
    <div className="w-full h-16 fixed top-0 left-0 flex gap-3 items-center justify-between bg-[#F0F0F0] px-6 z-[100]">
      <Link href={"/"} className="flex gap-3 items-center">
        <Image
          className="h-full aspect-square"
          width={60}
          height={60}
          src={"/logo.png"}
          alt="logo"
        ></Image>
        <div className="text-2xl font-bold">Инкогнито</div>
        {isAdminAccount && (
          <div className="text-2xl font-bold text-primary">Админ</div>
        )}
      </Link>

      <div className="flex gap-3 items-center">
        <Button asChild>
          <Link href="/offers">
            <Search />
          </Link>
        </Button>
        {username && username !== "$empty" && (
          <Button onClick={handleExit}>Выйти</Button>
        )}
        {username && (
          <Button onClick={handleRedirect}>
            {username == "$empty" ? "Войти" : username}
          </Button>
        )}
      </div>
    </div>
  );
}
