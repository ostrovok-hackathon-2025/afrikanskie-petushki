"use client";

import { useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@radix-ui/react-tabs";
import { cn } from "@/lib/utils";
import ReportView from "./View/view";
import ReportEdit from "./Edit/edit";
import { useRouter } from "next/navigation";

const tabsStyle = `transition-all duration-150 text-foreground hover:text-primary cursor-pointer py-2 border-b-transparent 
border-b-2 hover:border-b-primary`;
const tabsActiveStyle = "border-b-primary text-primary";

interface ReportProps {
  id: string;
  page: string;
}

export default function Report({ id, page }: ReportProps) {
  const [activeTab, setActiveTab] = useState(page);
  const router = useRouter();

  return (
    <div className="w-full h-full box-border px-2">
      <div className="max-w-1/2 mx-auto pt-24">
        <Tabs
          value={activeTab}
          onValueChange={(value) => router.replace(`/report/${id}/${value}`)}
        >
          <TabsList className="flex gap-3 border-b-2 mb-6">
            <TabsTrigger
              value="view"
              className={cn(tabsStyle, activeTab == "view" && tabsActiveStyle)}
            >
              просмотр
            </TabsTrigger>
            <TabsTrigger
              value="edit"
              className={cn(tabsStyle, activeTab == "edit" && tabsActiveStyle)}
            >
              редактирование
            </TabsTrigger>
          </TabsList>

          <TabsContent value="view">
            <ReportView id={id} goToEdit={() => setActiveTab("edit")} />
          </TabsContent>
          <TabsContent value="edit">
            <ReportEdit id={id} />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
