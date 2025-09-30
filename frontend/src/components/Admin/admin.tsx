"use client";

import { cn } from "@/lib/utils";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@radix-ui/react-tabs";
import { useRouter } from "next/navigation";
import AllOffers from "./AllOffers/all-offers";
import CreateOffer from "./CreateOffer/create-offer";
import Reports from "./Reports/reports";
import Applications from "./Applications/applications";

const tabsStyle = `transition-all duration-150 text-foreground hover:text-primary cursor-pointer py-2 border-b-transparent 
border-b-2 hover:border-b-primary`;
const tabsActiveStyle = "border-b-primary text-primary";

interface AdminProps {
  page: string;
}

export default function Admin({ page }: AdminProps) {
  const router = useRouter();

  return (
    <div className="w-full h-full box-border px-2">
      <div className="max-w-1/2 mx-auto box-border pt-20">
        <Tabs
          defaultValue="offers"
          value={page}
          onValueChange={(value) => router.replace(`/admin/${value}`)}
        >
          <TabsList className="flex gap-3 border-b-2 mb-6">
            <TabsTrigger
              value="offers"
              className={cn(tabsStyle, page == "offers" && tabsActiveStyle)}
            >
              Розыгрыши
            </TabsTrigger>

            <TabsTrigger
              value="create-offer"
              className={cn(
                tabsStyle,
                page == "create-offer" && tabsActiveStyle
              )}
            >
              Начать розыгрыш
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
          <TabsContent value="offers">
            <AllOffers />
          </TabsContent>
          <TabsContent value="create-offer">
            <CreateOffer />
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
