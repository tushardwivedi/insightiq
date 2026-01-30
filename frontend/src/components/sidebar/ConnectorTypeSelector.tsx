"use client";

import { ChevronLeft } from "lucide-react";

interface ConnectorType {
  type: string;
  name: string;
  description: string;
  icon: string;
}

const connectorTypes: ConnectorType[] = [
  {
    type: "file_upload",
    name: "File Upload",
    description: "Upload CSV or Excel files to query with SQL",
    icon: "\u{1F4C1}",
  },
  {
    type: "superset",
    name: "Apache Superset",
    description: "Connect to Apache Superset for analytics dashboards",
    icon: "\u{1F4CA}",
  },
  {
    type: "postgres",
    name: "PostgreSQL",
    description: "Connect to PostgreSQL database",
    icon: "\u{1F418}",
  },
  {
    type: "mysql",
    name: "MySQL",
    description: "Connect to MySQL database",
    icon: "\u{1F42C}",
  },
  {
    type: "api",
    name: "REST API",
    description: "Connect to external REST API",
    icon: "\u{1F310}",
  },
];

interface ConnectorTypeSelectorProps {
  onSelect: (type: string) => void;
  onBack: () => void;
}

export function ConnectorTypeSelector({ onSelect, onBack }: ConnectorTypeSelectorProps) {
  return (
    <>
      <button
        onClick={onBack}
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
            onClick={() => onSelect(type.type)}
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
  );
}

export { connectorTypes };
