"use client";

import Image from "next/image";
import ScrollVelocity from "../ScrollVelocity";
import { Button } from "../ui/button";
import Link from "next/link";
import { useRef } from "react";

export default function Landing() {
  const targetRef = useRef<HTMLDivElement>(null);

  const scrollToElement = () => {
    targetRef.current?.scrollIntoView({
      behavior: "smooth",
      block: "start",
    });
  };

  return (
    <div className="w-full relative flex flex-col justify-center items-center">
      <div className="w-full absolute top-0 left-0 overflow-hidden h-screen flex items-end pb-10 z-[-1]">
        <ScrollVelocity
          texts={["Фотография", "Рецензия", "Оценка"]}
          className="text-foreground/20"
        />
      </div>

      <div className="w-2/3 h-screen pt-40 px-auto ">
        <h1 className="text-6xl font-black mb-8">Инкогнито</h1>
        <div className="text-lg font-medium mb-12">
          Наш проект помогает сделать бронирование лучше. Здесь эксперт - вы.
        </div>

        <Button
          size={"lg"}
          className="text-xl cursor-pointer"
          onClick={scrollToElement}
        >
          Узнать больше
        </Button>
      </div>

      <div className="mb-20" ref={targetRef}></div>

      <div className="w-2/3 flex gap-4 mb-16">
        <div className="w-2/3">
          <h2 className="text-5xl font-bold mb-8">Испытайте отель инкогнито</h2>
          <p className="text-base">
            Вы — тайный рецензент, чья задача — оценить отель так, как это
            сделал бы обычный гость. Подавайте заявки на понравившиеся
            программы, и если ваша кандидатура будет выбрана, вас ждет поездка в
            указанные даты. Помните: главное правило — сохранять инкогнито.
            Персонал отеля не должен догадываться, что вы здесь с проверкой,
            чтобы ваши впечатления были максимально честными.
          </p>
        </div>
        <div className="w-1/3 aspect-square relative overflow-hidden rounded-lg">
          <Image
            src={"/landing/step1.jpg"}
            alt="step 1"
            fill
            className="absolute cover rounded-lg object-center"
          />
        </div>
      </div>

      <div className="w-2/3 flex gap-4 mb-16">
        <div className="w-1/3 aspect-square relative overflow-hidden rounded-lg">
          <Image
            src={"/landing/step2.jpg"}
            alt="step 2"
            fill
            className="absolute cover rounded-lg object-center"
          />
        </div>
        <div className="w-2/3">
          <h2 className="text-5xl font-bold mb-8">Запечатлейте все детали</h2>
          <p className="text-base">
            Во время визита вам предстоит создать фотоотчет. Внимательно изучите
            задание — в нем указано, какие именно зоны и детали отеля нужно
            снять. Старайтесь делать четкие, качественные фотографии при хорошем
            освещении. Они — главное доказательство вашего визита и основа для
            будущего отчета.
          </p>
        </div>
      </div>

      <div className="w-2/3 flex gap-4 mb-16">
        <div className="w-2/3">
          <h2 className="text-5xl font-bold mb-8">Напишите отчет</h2>
          <p className="text-base">
            После выезда из отеля вам нужно будет оформить все свои наблюдения в
            отчет. Опишите свои впечатления от сервиса, чистоты номеров, работы
            ресторана и других критериев, указанных в задании. Не забудьте
            прикрепить все сделанные фотографии. Ваш текст должен быть
            объективным, подробным и полезным для других путешественников.
          </p>
        </div>
        <div className="w-1/3 aspect-square relative overflow-hidden rounded-lg">
          <Image
            src={"/landing/step3.jpg"}
            alt="step 3"
            fill
            className="absolute cover rounded-lg object-center"
          />
        </div>
      </div>

      <div className="w-2/3 flex gap-4 mb-16">
        <div className="w-1/3 aspect-square relative overflow-hidden rounded-lg">
          <Image
            src={"/landing/step4.jpg"}
            alt="step 4"
            fill
            className="absolute cover rounded-lg object-center"
          />
        </div>
        <div className="w-2/3">
          <h2 className="text-5xl font-bold mb-8">Растите в сообществе</h2>
          <p className="text-base">
            Каждый ваш отчет оценивается нашими модераторами. Качественные и
            подробные отчеты повышают ваш личный рейтинг рецензента. Чем он
            выше, тем больше ваши шансы побеждать в конкурсах на самые
            интересные отели и получать промокоды от агрегаторов. Кроме того,
            высокая репутация открывает доступ к специальным достижениям .
            Помните: недобросовестные отчеты понижают рейтинг и уменьшают ваши
            шансы на новые миссии.
          </p>
        </div>
      </div>

      <Button asChild size={"lg"} className="text-2xl cursor-pointer h-14">
        <Link href={"/sign-up"}>Хочу участвовать</Link>
      </Button>
      <div className="h-20"></div>
    </div>
  );
}
