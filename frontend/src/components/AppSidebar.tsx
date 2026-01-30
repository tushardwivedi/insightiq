"use client";

import { useState, useEffect, useCallback } from "react";
import { DataConnector } from "@/types";
import { apiClient } from "@/lib/api";
import ConnectorCard from "./ConnectorCard";
import SupersetConnectorForm from "./connectors/SupersetConnectorForm";
import FileUploadModal from "./FileUploadModal";
import UploadedFilesList from "./UploadedFilesList";
import { usePathname } from "next/navigation";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  useSidebar,
} from "@/components/ui/sidebar";
import { Upload } from "lucide-react";
import { SidebarBrand } from "@/components/sidebar/SidebarBrand";
import { SidebarNavigation } from "@/components/sidebar/SidebarNavigation";
import { ConnectorTypeSelector } from "@/components/sidebar/ConnectorTypeSelector";
import { CollapsedSidebar } from "@/components/sidebar/CollapsedSidebar";
import { DataSourcesSection } from "@/components/sidebar/DataSourcesSection";

export function AppSidebar() {
  const { state } = useSidebar();
  const isCollapsed = state === "collapsed";
  const pathname = usePathname();

  const [connectors, setConnectors] = useState<DataConnector[]>([]);
  const [loading, setLoading] = useState(false);
  const [showAddForm, setShowAddForm] = useState(false);
  const [selectedConnectorType, setSelectedConnectorType] = useState<string>("");
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

  const handleConnectorUpdated = useCallback((updatedConnector: DataConnector) => {
    setConnectors((prev) =>
      prev.map((c) => (c.id === updatedConnector.id ? updatedConnector : c))
    );
  }, []);

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

  const renderContent = () => {
    if (isCollapsed) {
      return (
        <CollapsedSidebar
          pathname={pathname}
          connectors={connectors}
          onAddConnector={() => setShowAddForm(true)}
        />
      );
    }

    if (showAddForm) {
      return (
        <>
          {!selectedConnectorType ? (
            <ConnectorTypeSelector
              onSelect={setSelectedConnectorType}
              onBack={handleBackToConnectors}
            />
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
                  background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
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
        <SidebarNavigation pathname={pathname} />

        <div className="w-full h-px bg-border my-4" />

        <DataSourcesSection
          connectors={connectors}
          loading={loading}
          onAddConnector={() => setShowAddForm(true)}
          onTestConnector={handleTestConnection}
          onEditConnector={handleConnectorUpdated}
          onDeleteConnector={handleConnectorDeleted}
        />
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
          <SidebarBrand isCollapsed={isCollapsed} />
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
            {isCollapsed
              ? connectors.length
              : `${connectors.length} connector${connectors.length !== 1 ? "s" : ""} configured`}
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
