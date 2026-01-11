"use client";

import { useState, useEffect, useCallback } from "react";
import { DataConnector, ConnectorTestResult } from "@/types";
import { apiClient } from "@/lib/api";
import ConnectorCard from "./ConnectorCard";
import SupersetConnectorForm from "./connectors/SupersetConnectorForm";
import FileUploadModal from "./FileUploadModal";
import UploadedFilesList from "./UploadedFilesList";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Database,
  Plus,
  FileUp,
  BarChart3,
  Settings,
  Database as DatabaseIcon,
  FileText,
  ChevronLeft,
  Upload,
  LayoutDashboard,
  Sparkles,
  Zap,
  FolderOpen,
  History,
} from "lucide-react";

export function AppSidebar() {
  const { state } = useSidebar()
  const isCollapsed = state === 'collapsed'
  const pathname = usePathname()

  const [connectors, setConnectors] = useState<DataConnector[]>([]);
  const [loading, setLoading] = useState(false);
  const [showAddForm, setShowAddForm] = useState(false);
  const [selectedConnectorType, setSelectedConnectorType] =
    useState<string>("");
  const [showFileUploadModal, setShowFileUploadModal] = useState(false);
  const [fileListRefreshTrigger, setFileListRefreshTrigger] = useState(0);
  const [uploadedFilesCount, setUploadedFilesCount] = useState(0);

  const loadConnectors = useCallback(async () => {
    try {
      setLoading(true);
      const response = await apiClient.getConnectors();
      setConnectors(response.data || []);
    } catch (error) {
      console.error("Failed to load connectors:", error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadConnectors();
  }, [loadConnectors]);

  const handleConnectorAdded = useCallback((connector: DataConnector) => {
    setConnectors((prev) => [...prev, connector]);
    setShowAddForm(false);
    setSelectedConnectorType("");
  }, []);

  const handleConnectorUpdated = useCallback(
    (updatedConnector: DataConnector) => {
      setConnectors((prev) =>
        prev.map((c) => (c.id === updatedConnector.id ? updatedConnector : c))
      );
    },
    []
  );

  const handleConnectorDeleted = useCallback((connectorId: string) => {
    setConnectors((prev) => prev.filter((c) => c.id !== connectorId));
  }, []);

  const handleCancelForm = useCallback(() => {
    setSelectedConnectorType("");
  }, []);

  const handleBackToConnectors = useCallback(() => {
    setShowAddForm(false);
    setSelectedConnectorType("");
  }, []);

  const handleTestConnection = useCallback(async (connector: DataConnector) => {
    try {
      setConnectors((prev) =>
        prev.map((c) =>
          c.id === connector.id ? { ...c, status: "testing" } : c
        )
      );

      const result = await apiClient.testConnector(connector.id);

      setConnectors((prev) =>
        prev.map((c) =>
          c.id === connector.id
            ? {
                ...c,
                status: result.success ? "connected" : "error",
                last_tested: new Date().toISOString(),
              }
            : c
        )
      );
    } catch (error) {
      console.error("Connection test failed:", error);
      setConnectors((prev) =>
        prev.map((c) => (c.id === connector.id ? { ...c, status: "error" } : c))
      );
    }
  }, []);

  const connectorTypes = [
    {
      type: "file_upload",
      name: "File Upload",
      description: "Upload CSV or Excel files to query with SQL",
      icon: "📁",
    },
    {
      type: "superset",
      name: "Apache Superset",
      description: "Connect to Apache Superset for analytics dashboards",
      icon: "📊",
    },
    {
      type: "postgres",
      name: "PostgreSQL",
      description: "Connect to PostgreSQL database",
      icon: "🐘",
    },
    {
      type: "mysql",
      name: "MySQL",
      description: "Connect to MySQL database",
      icon: "🐬",
    },
    {
      type: "api",
      name: "REST API",
      description: "Connect to external REST API",
      icon: "🌐",
    },
  ];

  const mainNavItems = [
    { href: "/app", label: "Dashboard", icon: LayoutDashboard },
    { href: "/editor", label: "Query Editor", icon: Zap },
  ];

  const secondaryNavItems = [
    { href: "/history", label: "Query History", icon: History },
    { href: "/connectors", label: "Data Sources", icon: Database },
    { href: "/files", label: "Uploaded Files", icon: FolderOpen },
    { href: "/queries", label: "Saved Queries", icon: FileText },
    { href: "/analytics", label: "Analytics", icon: BarChart3 },
    { href: "/reports", label: "Reports", icon: FileText },
    { href: "/team", label: "Team", icon: Settings },
    { href: "/settings", label: "Settings", icon: Settings },
  ];

  const renderContent = () => {
    if (isCollapsed) {
      return (
        <div className="flex flex-col items-center space-y-4 py-4">
          {mainNavItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`w-8 h-8 rounded-lg flex items-center justify-center transition-colors ${
                pathname === item.href
                  ? "bg-primary text-primary-foreground"
                  : "hover:bg-muted text-muted-foreground"
              }`}
              title={item.label}
            >
              <item.icon className="w-4 h-4" />
            </Link>
          ))}

          <div className="w-6 h-px bg-border my-2" />

          <button
            onClick={() => setShowAddForm(true)}
            className="w-8 h-8 rounded-lg transition-colors flex items-center justify-center"
            style={{ background: 'var(--accent-color)', color: 'var(--primary-background)' }}
            title="Add Data Source"
          >
            <Plus className="w-4 h-4" />
          </button>

          {connectors.slice(0, 3).map((connector) => (
            <div
              key={connector.id}
              className="w-8 h-8 rounded-lg flex items-center justify-center text-xs font-medium"
              style={{
                backgroundColor: connector.status === 'connected' ? '#10b981' :
                                connector.status === 'testing' ? '#f59e0b' : '#ef4444',
                color: 'white'
              }}
              title={`${connector.name} - ${connector.status}`}
            >
              {connector.name.charAt(0).toUpperCase()}
            </div>
          ))}

          {connectors.length > 3 && (
            <div className="text-xs text-center" style={{ color: 'var(--text-secondary)' }}>
              +{connectors.length - 3}
            </div>
          )}
        </div>
      )
    }

    if (showAddForm) {
      return (
        <>
          <button
            onClick={handleBackToConnectors}
            className="mb-4 flex items-center gap-2 transition-colors"
            style={{ color: "var(--text-secondary)" }}
            onMouseEnter={(e) =>
              (e.currentTarget.style.color = "var(--text-primary)")
            }
            onMouseLeave={(e) =>
              (e.currentTarget.style.color = "var(--text-secondary)")
            }
          >
            <ChevronLeft className="w-4 h-4" />
            Back to connectors
          </button>

          {!selectedConnectorType ? (
            <>
              <h3
                className="text-md font-medium mb-4"
                style={{ color: "var(--text-primary)" }}
              >
                Choose Connector Type
              </h3>
              <div className="space-y-3">
                {connectorTypes.map((type) => (
                  <button
                    key={type.type}
                    onClick={() => setSelectedConnectorType(type.type)}
                    className="w-full p-4 text-left border rounded-lg transition-colors"
                    style={{ borderColor: "var(--border-color)" }}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.borderColor = "var(--accent-color)";
                      e.currentTarget.style.background = "var(--hover-surface)";
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.borderColor = "var(--border-color)";
                      e.currentTarget.style.background = "transparent";
                    }}
                  >
                    <div className="flex items-start gap-3">
                      <span className="text-2xl">{type.icon}</span>
                      <div className="flex-1">
                        <h4
                          className="font-medium"
                          style={{ color: "var(--text-primary)" }}
                        >
                          {type.name}
                        </h4>
                        <p
                          className="text-sm mt-1"
                          style={{ color: "var(--text-secondary)" }}
                        >
                          {type.description}
                        </p>
                      </div>
                    </div>
                  </button>
                ))}
              </div>
            </>
          ) : selectedConnectorType === "superset" ? (
            <SupersetConnectorForm
              onCancel={handleCancelForm}
              onSuccess={handleConnectorAdded}
            />
          ) : selectedConnectorType === "file_upload" ? (
            <div>
              <div className="flex items-center justify-between mb-4">
                <h3
                  className="text-md font-medium"
                  style={{ color: "var(--text-primary)" }}
                >
                  File Upload
                </h3>
                {uploadedFilesCount > 0 && (
                  <span
                    className="text-xs px-2 py-1 rounded-full font-medium"
                    style={{
                      background: "var(--hover-surface)",
                      color: "var(--text-secondary)",
                    }}
                  >
                    {uploadedFilesCount} file
                    {uploadedFilesCount !== 1 ? "s" : ""}
                  </span>
                )}
              </div>
              <p
                className="text-sm mb-4"
                style={{ color: "var(--text-secondary)" }}
              >
                Upload CSV or Excel files to query with SQL. Copy the table name
                and use it in your queries.
              </p>
              <button
                onClick={() => setShowFileUploadModal(true)}
                className="w-full px-4 py-3 rounded-lg font-medium transition-all flex items-center justify-center gap-2 mb-4"
                style={{
                  background:
                    "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
                  color: "white",
                }}
                onMouseEnter={(e) => (e.currentTarget.style.opacity = "0.9")}
                onMouseLeave={(e) => (e.currentTarget.style.opacity = "1")}
              >
                <Upload className="w-5 h-5" />
                Upload New File
              </button>
              <UploadedFilesList
                refreshTrigger={fileListRefreshTrigger}
                onCountChange={setUploadedFilesCount}
              />
            </div>
          ) : (
            <div
              className="text-center py-8"
              style={{ color: "var(--text-secondary)" }}
            >
              <p>Connector form for {selectedConnectorType} coming soon!</p>
              <button
                onClick={handleCancelForm}
                className="mt-4 px-4 py-2 rounded-md transition-colors"
                style={{
                  background: "var(--hover-surface)",
                  color: "var(--text-primary)",
                  border: "1px solid var(--border-color)",
                }}
              >
                Go Back
              </button>
            </div>
          )}
        </>
      );
    }

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

        <div className="w-full h-px bg-border my-4" />

        <SidebarGroup>
          <SidebarGroupLabel className="text-xs uppercase tracking-wider text-muted-foreground mb-2">
            Data Sources
          </SidebarGroupLabel>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton
                onClick={() => setShowAddForm(true)}
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
                onTest={() => handleTestConnection(connector)}
                onEdit={(connector) => handleConnectorUpdated(connector)}
                onDelete={() => handleConnectorDeleted(connector.id)}
              />
            ))}
          </div>
        )}
      </>
    );
  };

  return (
    <>
      <Sidebar
        style={
          {
            background: "var(--surface-color)",
            "--sidebar-width": "20rem",
          } as React.CSSProperties
        }
        collapsible="icon"
      >
        <SidebarHeader
          className="flex items-center justify-center p-4 border-b"
          style={{ borderColor: "var(--border-color)" }}
        >
          {isCollapsed ? (
            <Database className="w-6 h-6" style={{ color: 'var(--accent-color)' }} />
          ) : (
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}>
                <DatabaseIcon className="w-5 h-5 text-white" />
              </div>
              <div>
                <h2
                  className="text-base font-semibold"
                  style={{ color: "var(--text-primary)" }}
                >
                  InsightIQ
                </h2>
                <p className="text-xs" style={{ color: "var(--text-secondary)" }}>
                  Analytics Platform
                </p>
              </div>
            </div>
          )}
        </SidebarHeader>

        <SidebarContent className="flex-1 overflow-y-auto p-4">
          {renderContent()}
        </SidebarContent>

        <SidebarFooter
          className="p-4 border-t"
          style={{ borderColor: "var(--border-color)" }}
        >
          <div
            className="text-xs text-center"
            style={{ color: "var(--text-secondary)" }}
          >
            {isCollapsed ? connectors.length : `${connectors.length} connector${connectors.length !== 1 ? "s" : ""} configured`}
          </div>
        </SidebarFooter>
      </Sidebar>

      <FileUploadModal
        isOpen={showFileUploadModal}
        onClose={() => setShowFileUploadModal(false)}
        onUploadComplete={() => {
          setFileListRefreshTrigger((prev) => prev + 1);
        }}
      />
    </>
  );
}
