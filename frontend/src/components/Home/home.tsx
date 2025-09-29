"use client";
import { Tabs } from "../ui/tabs";
import { TabsContent, TabsList, TabsTrigger } from "@radix-ui/react-tabs";
import Profile from "./Profile/profile";
import { useState } from "react";
import { cn } from "@/lib/utils";
import Applications from "./Applications/applications";

const tabsStyle = `transition-all duration-150 text-foreground hover:text-primary cursor-pointer py-2 border-b-transparent 
border-b-2 hover:border-b-primary`;
const tabsActiveStyle = "border-b-primary text-primary";

export default function Home() {
  const [active, setActive] = useState("profile");

  return (
    <div className="w-full h-full box-border px-2">
      <div className="max-w-1/2 mx-auto box-border pt-20">
        <Tabs defaultValue="profile" value={active} onValueChange={setActive}>
          <TabsList className="flex gap-3 border-b-2 mb-6">
            <TabsTrigger
              value="profile"
              className={cn(tabsStyle, active == "profile" && tabsActiveStyle)}
            >
              Профиль
            </TabsTrigger>

            <TabsTrigger
              value="applications"
              className={cn(
                tabsStyle,
                active == "applications" && tabsActiveStyle
              )}
            >
              Заявки
            </TabsTrigger>
          </TabsList>
          <TabsContent value="profile">
            <Profile />
          </TabsContent>
          <TabsContent value="applications">
            <Applications />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
