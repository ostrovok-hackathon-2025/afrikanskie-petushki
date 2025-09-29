import { Bouncy } from "ldrs/react";
import "ldrs/react/Bouncy.css";
import { useEffect, useMemo, useState } from "react";
import { loadReport, ReportInfo } from "../model";
import { Button } from "@/components/ui/button";

interface Grade {
  result: string;
  title: string;
  description: string;
  raiting: string;
}

const GRADE_MAP = new Map<string, Grade>([
  [
    "accepted",
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

interface ReportViewProps {
  id: string;
  goToEdit: () => void;
}

export default function ReportView({ id, goToEdit }: ReportViewProps) {
  const [reportInfo, setReportInfo] = useState<ReportInfo | null>(null);

  useEffect(() => {
    (async () => {
      const res = await loadReport(id);

      setReportInfo(res);
    })();
  }, []);

  const statusCol = useMemo(() => {
    if (reportInfo?.status === "accepted") return "text-green-500";
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
    <div className="w-full">
      <div className="font-gain text-muted-foreground text-base mb-8">
        Ваш отчет по заявке
      </div>
      <div className="font-gain font-bold text-3xl mb-2">
        {reportInfo.hotelName}
      </div>
      <div className="font-gain text-base mb-6">{reportInfo.locationName}</div>

      <div className="font-gain font-medium text-xl mb-2">Задание</div>

      <div className="w-full whitespace-pre-wrap mb-8">{reportInfo.task}</div>

      <hr className="mb-12"></hr>

      {reportInfo.status === "created" ? (
        <>
          <div className="font-gain w-full flex flex-col items-center gap-4 pt-8 pb-20">
            Отчет еще не заполнен
            <Button onClick={goToEdit}>Перейти к заполнению</Button>
          </div>
        </>
      ) : (
        <>
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
        </>
      )}

      <hr className="mb-12"></hr>

      {reportInfo.status === "accepted" || reportInfo.status === "declined" ? (
        <>
          <div className="font-gain text-lg font-medium mb-6">
            Вердикт: <span className={statusCol}>{grade?.result + "!"}</span>
          </div>

          <div className="font-gain font-bold text-5xl mb-4">
            {grade?.title}
          </div>

          <div className="font-gain text-base text-muted-foreground mb-8">
            {grade?.description}
          </div>

          <div className="font-gain text-lg font-medium mb-2">
            Вы получаете <span className={statusCol}>{grade?.raiting}</span> к
            рейтингу{" "}
            {reportInfo.status === "accepted" &&
              "и промокод на бронирование в Островке"}
          </div>

          {reportInfo.status === "accepted" && (
            <div className="cursor-pointer rounded-sm text-2xl text-primary-foreground font-gain text-bold bg-primary w-fit font-bold box-border px-4 py-2">
              C47-D89
            </div>
          )}
        </>
      ) : (
        <>
          <div className="font-gain text-xl font-bold mb-16">
            Ожидает оценки
          </div>
        </>
      )}

      <div className="h-40"></div>
    </div>
  );
}
