"use client";

import { motion } from "framer-motion";
import { LucideIcon } from "lucide-react";

interface PageHeaderProps {
  icon: LucideIcon;
  title: string;
  description: string;
  actions?: React.ReactNode;
}

export function PageHeader({ icon: Icon, title, description, actions }: PageHeaderProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: -20 }}
      animate={{ opacity: 1, y: 0 }}
      className="flex items-center justify-between"
    >
      <div>
        <h1 className="text-3xl font-bold flex items-center gap-3">
          <Icon className="w-8 h-8 text-primary" />
          {title}
        </h1>
        <p className="text-muted-foreground mt-1">{description}</p>
      </div>
      {actions}
    </motion.div>
  );
}
