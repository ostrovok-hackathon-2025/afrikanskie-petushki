import { Bouncy } from "ldrs/react";
import "ldrs/react/Bouncy.css";
import { ChangeEvent, useEffect, useMemo, useRef, useState } from "react";
import { loadReport, ReportInfo } from "../model";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { PlusCircle, Trash, Trash2 } from "lucide-react";
import {
  Accordion,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { AccordionContent } from "@radix-ui/react-accordion";

interface ReportEditProps {
  id: string;
}

export default function ReportEdit({ id }: ReportEditProps) {
  const [reportText, setReportText] = useState("");
  const [reportInfo, setReportInfo] = useState<ReportInfo | null>(null);

  useEffect(() => {
    (async () => {
      const res = await loadReport(id);

      setReportInfo(res);
      setReportText(res.text);
    })();
  }, []);

  const editable = useMemo(
    () => reportInfo?.status === "created" || reportInfo?.status === "filled",
    [reportInfo?.status]
  );

  const textareaRef = useRef<HTMLTextAreaElement | null>(null);

  const autoResize = () => {
    if (!textareaRef.current) return;
    const textarea = textareaRef.current;
    textarea.style.height = "auto";
    textarea.style.height = textarea.scrollHeight + 10 + "px";
  };

  useEffect(() => {
    autoResize();
  }, [reportText]);

  const fileInputRef = useRef<HTMLInputElement>(null);
  const [previewUrls, setPreviewUrls] = useState<string[]>([]);
  const [selectedImages, setSelectedImages] = useState<File[]>([]);

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;

    if (!files || files.length === 0) return;

    const validFiles: File[] = [];
    const newPreviewUrls: string[] = [];

    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      if (file.type === "image/png" || file.type === "image/jpeg") {
        validFiles.push(file);
        newPreviewUrls.push(URL.createObjectURL(file));
      }
    }

    previewUrls.forEach((url) => URL.revokeObjectURL(url));

    setPreviewUrls(newPreviewUrls);
    setSelectedImages(validFiles);
  };

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  useEffect(() => {
    return () => previewUrls.forEach((url) => URL.revokeObjectURL(url));
  }, []);

  const removeImage = (index: number) => {
    const newPreviewUrls = [...previewUrls];
    const removedUrl = newPreviewUrls.splice(index, 1)[0];

    // Освобождаем память
    URL.revokeObjectURL(removedUrl);

    setPreviewUrls(newPreviewUrls);

    const images = [...selectedImages].splice(index, 1);
    setSelectedImages(images);
  };

  const handleSubmit = () => {};

  if (!reportInfo)
    return (
      <div className="w-full h-full flex items-center justify-center gap-8">
        <div className="font-gain text-xl">Загружаем отчет</div>
        <Bouncy size="45" speed="1.75" color="black" />
      </div>
    );

  return (
    <div className="w-full">
      {!editable && (
        <div className="font-gain text-xl font-bold mb-16">
          Этот отчет уже получил оценку. Его нельзя редактировать
        </div>
      )}

      <Accordion type="single" collapsible className="mb-8">
        <AccordionItem value="item-1">
          <AccordionTrigger className="font-gain font-medium text-xl mb-2 ">
            Задание
          </AccordionTrigger>
          <AccordionContent className="font-gain whitespace-pre-wrap">
            {reportInfo.task}
          </AccordionContent>
        </AccordionItem>
      </Accordion>

      <hr className="mb-12"></hr>

      <form>
        <div className="font-gain font-medium text-xl mb-6">
          Ваши впечатления
        </div>
        <textarea
          value={reportText}
          onChange={(e) => setReportText(e.target.value)}
          disabled={!editable}
          className={cn(
            "outline-none border w-full h-fit rounded-lg border box-border p-2 resize-none mb-12",
            !editable && "text-muted-foreground"
          )}
          ref={textareaRef}
          name="text"
        ></textarea>

        <div className="font-gain font-medium text-xl mb-6">
          Ваши фотографии
        </div>

        {editable ? (
          <div className="mb-16">
            <input
              type="file"
              ref={fileInputRef}
              onChange={handleFileChange}
              accept="image/png, image/jpeg, image/jpg"
              multiple
              style={{ display: "none" }}
            />
            <Button
              onClick={handleClick}
              className="px-6 py-3 mb-4"
              type="button"
            >
              <PlusCircle color="#FFF" className="w-8 h-8 mx-6" />
            </Button>

            <div className="w-full flex flex-wrap gap-4">
              {previewUrls.map((e, i) => (
                <div key={e} className="w-[240px] h-[135px] relative">
                  <img
                    src={e}
                    key={e}
                    className="absolute top-0 left-0 rounded-sm"
                    style={{
                      width: 240,
                      height: 135,
                      objectFit: "cover",
                    }}
                  />
                  <div
                    className={cn(
                      "absolute w-full h-full bg-[#0004] opacity-0 hover:opacity-100",
                      "flex justify-center items-center transition-all cursor-pointer rounded-sm"
                    )}
                    onClick={() => removeImage(i)}
                  >
                    <Trash2 color="#FFF" />
                  </div>
                </div>
              ))}
            </div>
          </div>
        ) : (
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
        )}

        {editable && (
          <Button
            onClick={handleSubmit}
            type="button"
            size={"lg"}
            className="text-xl w-40 h-14"
          >
            Сохранить
          </Button>
        )}
      </form>
    </div>
  );
}
