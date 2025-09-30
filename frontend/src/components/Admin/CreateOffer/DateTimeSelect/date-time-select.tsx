"use client";

import { useState, useEffect, useMemo } from "react";
import { Calendar } from "@/components/ui/calendar";
import { format, addDays, setDate } from "date-fns";
import { ru } from "date-fns/locale";
import { Clock } from "lucide-react";

interface DateTimeSelectProps {
  onDateTimeChange: (timestamp: number) => void;
  minDate: Date;
}

export function DateTimeSelect({
  onDateTimeChange,
  minDate,
}: DateTimeSelectProps) {
  const [selectedDate, setSelectedDate] = useState<Date>(addDays(minDate, 1));
  const [selectedTime, setSelectedTime] = useState<string>("12:00");

  // Обработчик изменения даты
  const handleDateChange = (date: Date | undefined) => {
    if (date) {
      setSelectedDate(date);
    }
  };

  // Обработчик изменения времени
  const handleTimeChange = (time: string) => {
    setSelectedTime(time);
  };

  // Формируем timestamp и передаем наружу
  useEffect(() => {
    if (selectedDate && selectedTime) {
      const [hours, minutes] = selectedTime.split(":").map(Number);
      const dateTime = new Date(selectedDate);
      dateTime.setHours(hours, minutes, 0, 0);

      // Преобразуем в timestamp (миллисекунды)
      const timestamp = dateTime.getTime();

      if (timestamp >= minDate.getTime()) {
        onDateTimeChange(timestamp);
        return;
      }

      setSelectedDate(minDate);
      setSelectedTime("12:00");
      onDateTimeChange(minDate.getTime());
    }
  }, [selectedDate, selectedTime, minDate, onDateTimeChange]);

  // Функция для блокировки прошедших дат
  const isDateDisabled = (date: Date) => {
    const tomorrow = addDays(minDate, 1);
    tomorrow.setHours(0, 0, 0, 0);
    return date < tomorrow;
  };

  // Генерация доступных временных слотов
  const timeSlots = useMemo(() => {
    const slots = [];
    const startHour = 0;
    const endHour = 23;

    for (let hour = startHour; hour <= endHour; hour++) {
      const timeString = `${hour.toString().padStart(2, "0")}:00`;
      slots.push(timeString);
    }

    return slots;
  }, []);

  return (
    <div className="flex gap-4 p-4">
      {/* Календарь */}
      <div>
        <h3 className="text-sm font-medium mb-2">Выберите дату</h3>
        <Calendar
          mode="single"
          selected={selectedDate}
          onSelect={handleDateChange}
          disabled={isDateDisabled}
          locale={ru}
          className="rounded-md border"
        />
      </div>

      {/* Выбор времени */}
      <div>
        <h3 className="text-sm font-medium mb-2 flex items-center gap-2">
          <Clock className="w-4 h-4" />
          Выберите время (не ранее 12:00)
        </h3>
        <div className="grid grid-cols-3 gap-2 max-h-full overflow-y-auto">
          {timeSlots.map((time) => (
            <button
              key={time}
              type="button"
              onClick={() => handleTimeChange(time)}
              className={`p-2 text-sm border rounded-md transition-colors ${
                selectedTime === time
                  ? "bg-primary text-white border-none"
                  : "bg-white text-gray-700 border-gray-300 hover:bg-gray-50"
              }`}
            >
              {time}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
