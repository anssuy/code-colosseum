"use client";

import { History, ListChecks, Search, Trophy } from "lucide-react";
import Link from "next/link";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";

const navigation = [
  {
    label: "Play",
    items: [{ title: "Find Match", url: "/play", icon: Search }],
  },
  {
    label: "Compete",
    items: [
      { title: "Match History", url: "/history", icon: History },
      { title: "Leaderboard", url: "/leaderboard", icon: Trophy },
    ],
  },
  {
    label: "Problems",
    items: [{ title: "All Problems", url: "/problems", icon: ListChecks }],
  },
];

export default function AppSidebar() {
  return (
    <Sidebar className="border-zinc-200">
      <SidebarContent className="px-2 py-3">
        {navigation.map((group) => (
          <SidebarGroup className="py-2" key={group.label}>
            <SidebarGroupLabel className="px-2 text-xs text-zinc-400 uppercase tracking-wider">
              {group.label}
            </SidebarGroupLabel>

            <SidebarGroupContent>
              <SidebarMenu className="space-y-2">
                {group.items.map((item) => (
                  <Link href={item.url} key={item.url}>
                    <SidebarMenuItem>
                      <SidebarMenuButton className="flex cursor-pointer items-center gap-2">
                        {item.icon && <item.icon className="size-4" />}
                        <span>{item.title}</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  </Link>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>
      <SidebarFooter />
    </Sidebar>
  );
}
