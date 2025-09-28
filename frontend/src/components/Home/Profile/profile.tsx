import { getSecretGuestAPI } from "@/api/api";
import { Input } from "@/components/ui/input";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { getSession } from "next-auth/react";
import { redirect } from "next/navigation";
import { useEffect, useState } from "react";

const { getUser } = getSecretGuestAPI();

export default function Profile() {
    const [username, setUsername] = useState("");

    useEffect(() => {
        (async () => {
            const session = await getSession();

            if (!session) return redirect("log-in");

            const resp = await getUser({ headers: withAuthHeader(session) });

            setUsername(resp.data.ostrovok_login ?? "");
        })();
    }, []);

    return <div className="w-full">
        <div className="font-gain text-2xl mb-5">Профиль</div>

        <div className="text-base mb-1">
            Пользователь
        </div>
        <Input value={username} className="mb-3 max-w-1/3" readOnly></Input>

        <div className="text-base">Raiting: {" "}
            <div className="text-primary text-3xl font-bold animate-bounce inline-block">100</div>
        </div>
    </div>
}