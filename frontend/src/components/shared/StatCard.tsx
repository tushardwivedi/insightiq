import { Card, CardContent } from "@/components/ui/card";
import { LucideIcon } from "lucide-react";

interface StatCardProps {
  icon: LucideIcon;
  iconClassName?: string;
  iconBgClassName?: string;
  value: string | number;
  label: string;
}

export function StatCard({
  icon: Icon,
  iconClassName = "text-primary",
  iconBgClassName = "bg-blue-500/10",
  value,
  label,
}: StatCardProps) {
  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex items-center gap-3">
          <div className={`w-10 h-10 rounded-lg ${iconBgClassName} flex items-center justify-center`}>
            <Icon className={`w-5 h-5 ${iconClassName}`} />
          </div>
          <div>
            <p className="text-sm text-muted-foreground">{label}</p>
            <p className="text-2xl font-bold">{value}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
