'use client'

import { motion } from 'framer-motion'
import { LucideIcon } from 'lucide-react'

interface ChartCardProps {
  icon: LucideIcon;
  iconClassName: string;
  title: string;
  children: React.ReactNode;
}

const itemVariants = {
  hidden: { opacity: 0, scale: 0.95 },
  visible: {
    opacity: 1,
    scale: 1,
    transition: { duration: 0.5 },
  },
};

export function ChartCard({ icon: Icon, iconClassName, title, children }: ChartCardProps) {
  return (
    <motion.div variants={itemVariants} className="card p-6  ">
      <div className="flex items-center gap-2 mb-4">
        <Icon className={`w-5 h-5 ${iconClassName}`} />
        <h4 className="font-semibold" style={{ color: 'var(--text-primary)' }}>{title}</h4>
      </div>
      <div className="h-64">
        {children}
      </div>
    </motion.div>
  );
}
