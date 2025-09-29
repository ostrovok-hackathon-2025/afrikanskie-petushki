"use client";
import { Tabs } from "../ui/tabs";
import { TabsContent, TabsList, TabsTrigger } from "@radix-ui/react-tabs";
import Profile from "./Profile/profile";
import { useState } from "react";
import { cn } from "@/lib/utils";
import Applications from "./Applications/applications";
import { useRouter } from "next/navigation";
import Reports from "./Reports/reports";

const tabsStyle = `transition-all duration-150 text-foreground hover:text-primary cursor-pointer py-2 border-b-transparent 
border-b-2 hover:border-b-primary`;
const tabsActiveStyle = "border-b-primary text-primary";

interface HomeProps {
  page: string;
}

export default function Home({ page }: HomeProps) {
  const router = useRouter();

  return (
    <div className="w-full h-full box-border px-2">
      <div className="max-w-1/2 mx-auto box-border pt-20">
        <Tabs
          defaultValue="profile"
          value={page}
          onValueChange={(value) => router.replace(`/home/${value}`)}
        >
          <TabsList className="flex gap-3 border-b-2 mb-6">
            <TabsTrigger
              value="profile"
              className={cn(tabsStyle, page == "profile" && tabsActiveStyle)}
            >
              Профиль
            </TabsTrigger>

            <TabsTrigger
              value="applications"
              className={cn(
                tabsStyle,
                page == "applications" && tabsActiveStyle
              )}
            >
              Заявки
            </TabsTrigger>

            <TabsTrigger
              value="reports"
              className={cn(tabsStyle, page == "reports" && tabsActiveStyle)}
            >
              Отчеты
            </TabsTrigger>
          </TabsList>
          <TabsContent value="profile">
            <Profile />
          </TabsContent>
          <TabsContent value="applications">
            <Applications />
          </TabsContent>
          <TabsContent value="reports">
            <Reports />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
