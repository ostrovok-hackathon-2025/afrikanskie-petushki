"use client";

import { useEffect, useMemo, useState } from "react";
import { Bouncy } from "ldrs/react";
import "ldrs/react/Bouncy.css";
import { report } from "process";
import Image from "next/image";

interface ReportInfo {
  hotelName: string;
  locationName: string;
  task: string;

  text: string;
  images: {
    id: string;
    url: string;
  }[];

  status: string;
}

interface Grade {
  result: string;
  title: string;
  description: string;
  raiting: string;
}

const GRADE_MAP = new Map<string, Grade>([
  [
    "ok",
    {
      result: "принят",
      title: "Так держать! ✨",
      description:
        "Вы отлично справились с заданием! Ваши усилия помогут нам сделать опыт других людей лучше. Благодарим за сотрудничество.",
      raiting: "+20",
    },
  ],
  [
    "rejected",
    {
      result: "отклонён",
      title: "Ой-ооой! 🤔 Вам есть над чем поработать",
      description:
        "К сожалению, в этот раз работа не соответствует требуемым критериям. Мы верим, что вы учтёте замечания и покажете лучший результат в следующий раз.",
      raiting: "-5",
    },
  ],
]);

interface ReportProps {
  id: string;
}

export default function Report({ id }: ReportProps) {
  const [reportInfo, setReportInfo] = useState<ReportInfo | null>(null);

  useEffect(() => {
    setReportInfo({
      hotelName: "Superhotel",
      locationName: "Москва",
      task: "сделать траляля\nсфотографировать труляля",

      images: [
        {
          id: "1",
          url: "https://www.gentinghotel.co.uk/_next/image?url=https%3A%2F%2Fs3.eu-west-2.amazonaws.com%2Fstaticgh.gentinghotel.co.uk%2Fuploads%2Fhero%2FSuiteNov2022_0008_1920.jpg&w=3840&q=75",
        },
        {
          id: "3",
          url: "https://www.potawatomi.com/application/files/3517/4560/6138/Signature-2-Queen_body.webp",
        },
        {
          id: "4",
          url: "https://images.squarespace-cdn.com/content/v1/60a23657a6164d69e38ddad0/34e78be0-e3f9-4483-9a50-750b62a7b746/PM_FOOD_WEB_3.jpg",
        },
        {
          id: "5",
          url: "https://media.istockphoto.com/id/625006196/photo/sunrise-on-a-tropical-island-palm-trees-on-sandy-beach.jpg?s=612x612&w=0&k=20&c=qGNG4XX4d3SNPDgLgM0GpdEtcPhyldWzQTd38KoC1X8=",
        },
      ],
      text: "Все понравилось\n\nОтличное место\nПравдивые отзывы\nВкусная еда",

      status: "rejected",
    });
  }, []);

  const statusCol = useMemo(() => {
    if (reportInfo?.status === "ok") return "text-green-500";
    if (reportInfo?.status === "rejected") return "text-destructive";
    return "";
  }, [reportInfo?.status]);

  const grade = GRADE_MAP.get(reportInfo?.status ?? "");

  if (!reportInfo)
    return (
      <div className="w-full h-full flex items-center justify-center gap-8">
        <div className="font-gain text-xl">Загружаем отчет</div>
        <Bouncy size="45" speed="1.75" color="black" />
      </div>
    );

  return (
    <div className="w-full h-full box-border px-2">
      <div className="max-w-1/2 mx-auto pt-24">
        <div className="font-gain text-muted-foreground text-base mb-8">
          Ваш отчет по заявке
        </div>
        <div className="font-gain font-bold text-3xl mb-2">
          {reportInfo.hotelName}
        </div>
        <div className="font-gain text-base mb-6">
          {reportInfo.locationName}
        </div>

        <div className="font-gain font-medium text-xl mb-2">Задание</div>

        <div className="w-full whitespace-pre-wrap mb-8">{reportInfo.task}</div>

        <hr className="mb-12"></hr>

        <div className="font-gain font-medium text-xl mb-2">
          Ваши впечатления
        </div>

        <div className="w-full rounded-lg border box-border p-2 whitespace-pre-wrap mb-6">
          {reportInfo.text}
        </div>

        <div className="font-gain font-medium text-xl mb-2">
          Ваши фотографии
        </div>

        <div className="w-full grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-4 mb-8">
          {reportInfo.images.map((e, i) => (
            <img
              key={i}
              src={e.url}
              alt="image"
              className="cover w-full h-[175px] rounded-lg"
            />
          ))}
        </div>

        <hr className="mb-12"></hr>

        <div className="font-gain text-lg font-medium mb-6">
          Вердикт: <span className={statusCol}>{grade?.result + "!"}</span>
        </div>

        <div className="font-gain font-bold text-5xl mb-4">{grade?.title}</div>

        <div className="font-gain text-base text-muted-foreground mb-8">
          {grade?.description}
        </div>

        <div className="font-gain text-lg font-medium mb-2">
          Вы получается <span className={statusCol}>{grade?.raiting}</span> к
          рейтингу{" "}
          {reportInfo.status === "ok" &&
            "и промокод на бронирование в Островке"}
        </div>

        {reportInfo.status === "ok" && (
          <div className="cursor-pointer rounded-sm text-2xl text-primary-foreground font-gain text-bold bg-primary w-fit font-bold box-border px-4 py-2">
            C47-D89
          </div>
        )}

        <div className="h-40"></div>
      </div>
    </div>
  );
}
