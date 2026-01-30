import { Database, Database as DatabaseIcon } from "lucide-react";

interface SidebarBrandProps {
  isCollapsed: boolean;
}

export function SidebarBrand({ isCollapsed }: SidebarBrandProps) {
  if (isCollapsed) {
    return <Database className="w-6 h-6" style={{ color: 'var(--accent-color)' }} />;
  }

  return (
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
  );
}
