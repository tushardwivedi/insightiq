"use client";

import { motion } from "framer-motion";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Shield, Trash2, Crown, User, CheckCircle, XCircle } from "lucide-react";

interface TeamMember {
  id: string;
  email: string;
  name?: string;
  role: string;
  status: "active" | "pending" | "inactive";
  created_at: string;
  avatar?: string;
}

interface TeamMemberRowProps {
  member: TeamMember;
  index: number;
  onRemove: (id: string) => void;
}

function getRoleIcon(role: string) {
  switch (role) {
    case "admin": return <Crown className="w-4 h-4 text-yellow-500" />;
    case "editor": return <Shield className="w-4 h-4 text-blue-500" />;
    default: return <User className="w-4 h-4 text-gray-500" />;
  }
}

function getRoleBadgeColor(role: string) {
  switch (role) {
    case "admin": return "bg-yellow-500/10 text-yellow-500 border-yellow-500/20";
    case "editor": return "bg-blue-500/10 text-blue-500 border-blue-500/20";
    default: return "bg-gray-500/10 text-gray-500 border-gray-500/20";
  }
}

function formatDate(date?: string) {
  if (!date) return "";
  return new Date(date).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export function TeamMemberRow({ member, index, onRemove }: TeamMemberRowProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.03 }}
      className="p-4 rounded-lg border hover:bg-muted/30 transition-colors"
    >
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <div className="w-12 h-12 rounded-full bg-gradient-to-br from-primary/20 to-purple-500/20 flex items-center justify-center text-lg font-semibold">
            {member.email.charAt(0).toUpperCase()}
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h3 className="font-semibold">{member.name || member.email}</h3>
              {getRoleIcon(member.role)}
            </div>
            <p className="text-sm text-muted-foreground">{member.email}</p>
            <div className="flex items-center gap-3 mt-1">
              <Badge
                variant="outline"
                className={`text-xs ${getRoleBadgeColor(member.role)}`}
              >
                {member.role}
              </Badge>
              {member.status === "active" ? (
                <span className="flex items-center gap-1 text-xs text-green-500">
                  <CheckCircle className="w-3 h-3" />
                  Active
                </span>
              ) : member.status === "pending" ? (
                <span className="flex items-center gap-1 text-xs text-orange-500">
                  <XCircle className="w-3 h-3" />
                  Pending
                </span>
              ) : null}
              <span className="text-xs text-muted-foreground">
                Joined {formatDate(member.created_at)}
              </span>
            </div>
          </div>
        </div>
        <div className="flex items-center gap-1">
          <Button variant="ghost" size="icon" title="Edit">
            <Shield className="w-4 h-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onRemove(member.id)}
            title="Remove"
          >
            <Trash2 className="w-4 h-4 text-destructive" />
          </Button>
        </div>
      </div>
    </motion.div>
  );
}
