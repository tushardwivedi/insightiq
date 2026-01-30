import { LucideIcon } from "lucide-react";

interface EmptyStateProps {
  icon: LucideIcon;
  message: string;
  children?: React.ReactNode;
}

export function EmptyState({ icon: Icon, message, children }: EmptyStateProps) {
  return (
    <div className="text-center py-12">
      <Icon className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
      <p className="text-muted-foreground mb-4">{message}</p>
      {children}
    </div>
  );
}
