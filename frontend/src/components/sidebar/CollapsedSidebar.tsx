import Link from "next/link";
import { Plus } from "lucide-react";
import { DataConnector } from "@/types";
import { mainNavItems } from "./SidebarNavigation";

interface CollapsedSidebarProps {
  pathname: string;
  connectors: DataConnector[];
  onAddConnector: () => void;
}

export function CollapsedSidebar({ pathname, connectors, onAddConnector }: CollapsedSidebarProps) {
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
        onClick={onAddConnector}
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
  );
}
