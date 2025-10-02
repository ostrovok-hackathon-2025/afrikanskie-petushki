import { getSecretGuestAPI } from "@/api/api";
import { DocsAchievementResponse, DocsUserResponse } from "@/api/model";
import CountUp from "@/components/CountUp";
import Silk from "@/components/Silk";
import { Input } from "@/components/ui/input";
import { withAuthHeader } from "@/lib/next-auth/with-auth-header";
import { cn } from "@/lib/utils";
import { getSession } from "next-auth/react";
import Image from "next/image";
import { redirect } from "next/navigation";
import { useEffect, useState } from "react";

const { getUser } = getSecretGuestAPI();

function AchievementCard({ name, raiting_limit }: DocsAchievementResponse) {
  return (
    <div className="box-border rounded-lg border p-1 w-full flex items-center gap-2">
      <Image
        src={"/profile/achievement.jpg"}
        width={80}
        height={80}
        alt="achievement"
        className="aspect-square"
      />
      <div>
        <div className="text-lg mb-1 font-medium">{name}</div>
        <div className="text-sm">Получить рейтинг не менее {raiting_limit}</div>
      </div>
    </div>
  );
}

export default function Profile() {
  const [username, setUsername] = useState("загружаем...");
  const [email, setEmail] = useState("загружаем...");
  const [raiting, setRaiting] = useState(0);
  const [achievements, setAchievements] = useState<DocsAchievementResponse[]>(
    []
  );

  useEffect(() => {
    (async () => {
      const session = await getSession();

      if (!session) return redirect("log-in");

      const resp = await getUser({ headers: withAuthHeader(session) });

      setUsername(resp.data.ostrovok_login ?? "");
      setEmail(resp.data.email ?? "");
      setRaiting(resp.data.rating ?? 0);
      setAchievements(resp.data.achievements ?? []);
    })();
  }, []);

  return (
    <div className="w-full">
      <div className="w-full font-bold text-3xl text-center mb-8 mt-12">
        Профиль
      </div>

      <div className="flex gap-4 mb-12">
        <div className="w-1/3 box-border rounded-lg border p-4">
          <div className="text-lg mb-1 font-medium">Пользователь</div>
          <Input value={username} className="mb-3 w-full" readOnly></Input>

          <div className="text-lg mb-1 font-medium">Почта</div>
          <Input value={email} className="mb-3 w-full" readOnly></Input>
        </div>
        <div className="relative w-2/3 h-[400px] rounded-lg overflow-hidden">
          <Silk
            speed={5}
            scale={1}
            color="#93389d"
            noiseIntensity={1.5}
            rotation={0}
          />
          <div
            className={cn(
              "text-primary-foreground font-bold text-6xl",
              " absolute top-0 left-0 w-full h-full flex justify-center items-center"
            )}
          >
            Рейтинг:&nbsp;{" "}
            <CountUp
              from={0}
              to={raiting}
              direction="up"
              duration={1 + raiting / 25}
            />
          </div>
        </div>
      </div>

      <div className="w-full font-bold text-3xl text-center mb-8">
        Достижения
      </div>

      <div className="flex w-full gap-6 flex-col items-center">
        {achievements.map((e, i) => (
          <AchievementCard {...e} key={i} />
        ))}
        {achievements.length === 0 && (
          <div className="text-sm text-center">
            вы еще не заработали достижений, копите рейтинг для разблокировки
          </div>
        )}
      </div>
      <div className="h-20"></div>
    </div>
  );
}
