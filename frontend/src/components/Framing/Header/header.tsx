"use client";
import { getSecretGuestAPI } from "@/api/api";
import { Button } from "@/components/ui/button";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { Search } from "lucide-react";
import { getSession, signOut } from "next-auth/react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

const { getUser } = getSecretGuestAPI();

export default function Header() {
  const [username, setUsername] = useState("");
  const router = useRouter();

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) {
        setUsername("$empty");
        return;
      }

      const resp = await getUser({ headers: withAuthHeader(session) });

      setUsername(resp.data.ostrovok_login ?? "");
    })();
  }, []);

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
    <div className="w-full h-16 fixed top-0 left-0 flex gap-3 items-center justify-end bg-[#F0F0F0] px-6 z-[100]">
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
  );
}
