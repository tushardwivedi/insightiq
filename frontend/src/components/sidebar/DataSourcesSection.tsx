"use client";

import { Plus, Database } from "lucide-react";
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import ConnectorCard from "@/components/ConnectorCard";
import { DataConnector } from "@/types";

interface DataSourcesSectionProps {
  connectors: DataConnector[];
  loading: boolean;
  onAddConnector: () => void;
  onTestConnector: (connector: DataConnector) => void;
  onEditConnector: (connector: DataConnector) => void;
  onDeleteConnector: (connectorId: string) => void;
}

export function DataSourcesSection({
  connectors,
  loading,
  onAddConnector,
  onTestConnector,
  onEditConnector,
  onDeleteConnector,
}: DataSourcesSectionProps) {
  return (
    <>
      <SidebarGroup>
        <SidebarGroupLabel className="text-xs uppercase tracking-wider text-muted-foreground mb-2">
          Data Sources
        </SidebarGroupLabel>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              onClick={onAddConnector}
              className="w-full"
            >
              <Plus className="w-5 h-5" />
              <span>Add Data Source</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarGroup>

      {loading ? (
        <div className="flex items-center justify-center py-8">
          <div
            className="animate-spin rounded-full h-8 w-8 border-b-2"
            style={{ borderColor: "var(--accent-color)" }}
          ></div>
        </div>
      ) : connectors.length === 0 ? (
        <div
          className="text-center py-8"
          style={{ color: "var(--text-secondary)" }}
        >
          <Database
            className="w-12 h-12 mx-auto mb-4"
            style={{ color: "var(--border-color)" }}
          />
          <p className="text-sm">No data sources connected</p>
          <p
            className="text-xs mt-1"
            style={{ color: "var(--text-secondary)" }}
          >
            Add your first connector to get started
          </p>
        </div>
      ) : (
        <div className="space-y-3 mt-3">
          {connectors.map((connector) => (
            <ConnectorCard
              key={connector.id}
              connector={connector}
              onTest={() => onTestConnector(connector)}
              onEdit={(connector) => onEditConnector(connector)}
              onDelete={() => onDeleteConnector(connector.id)}
            />
          ))}
        </div>
      )}
    </>
  );
}
