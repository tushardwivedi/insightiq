import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import {
  Database,
  FileText,
  BarChart3,
  Settings,
  FolderOpen,
  History,
  LayoutDashboard,
  Zap,
  Sparkles,
} from "lucide-react";
import { LucideIcon } from "lucide-react";

interface NavItem {
  href: string;
  label: string;
  icon: LucideIcon;
}

const mainNavItems: NavItem[] = [
  { href: "/app", label: "Dashboard", icon: LayoutDashboard },
  { href: "/editor", label: "Query Editor", icon: Zap },
];

const secondaryNavItems: NavItem[] = [
  { href: "/history", label: "Query History", icon: History },
  { href: "/connectors", label: "Data Sources", icon: Database },
  { href: "/files", label: "Uploaded Files", icon: FolderOpen },
  { href: "/queries", label: "Saved Queries", icon: FileText },
  { href: "/analytics", label: "Analytics", icon: BarChart3 },
  { href: "/reports", label: "Reports", icon: FileText },
  { href: "/team", label: "Team", icon: Settings },
  { href: "/settings", label: "Settings", icon: Settings },
];

interface SidebarNavigationProps {
  pathname: string;
}

export function SidebarNavigation({ pathname }: SidebarNavigationProps) {
  return (
    <>
      <SidebarGroup>
        <SidebarGroupLabel className="text-xs uppercase tracking-wider text-muted-foreground mb-2">
          Navigation
        </SidebarGroupLabel>
        <SidebarMenu>
          {mainNavItems.map((item) => (
            <SidebarMenuItem key={item.href}>
              <SidebarMenuButton asChild isActive={pathname === item.href}>
                <Link href={item.href} className="flex items-center gap-3">
                  <item.icon className="w-5 h-5" />
                  <span>{item.label}</span>
                  {item.href === '/editor' && (
                    <Badge variant="secondary" className="ml-auto text-xs">
                      <Sparkles className="w-3 h-3 mr-1" />
                      AI
                    </Badge>
                  )}
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarGroup>

      <SidebarGroup>
        <SidebarGroupLabel className="text-xs uppercase tracking-wider text-muted-foreground mb-2 mt-4">
          Quick Access
        </SidebarGroupLabel>
        <SidebarMenu>
          {secondaryNavItems.map((item) => (
            <SidebarMenuItem key={item.href}>
              <SidebarMenuButton asChild isActive={pathname === item.href}>
                <Link href={item.href} className="flex items-center gap-3">
                  <item.icon className="w-5 h-5" />
                  <span>{item.label}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarGroup>
    </>
  );
}

export { mainNavItems, secondaryNavItems };
