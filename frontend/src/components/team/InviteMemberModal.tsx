"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Mail } from "lucide-react";

interface InviteMemberModalProps {
  inviteEmail: string;
  inviteRole: string;
  onEmailChange: (email: string) => void;
  onRoleChange: (role: string) => void;
  onInvite: () => void;
  onClose: () => void;
}

export function InviteMemberModal({
  inviteEmail,
  inviteRole,
  onEmailChange,
  onRoleChange,
  onInvite,
  onClose,
}: InviteMemberModalProps) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <Card className="w-full max-w-md mx-4">
        <CardHeader>
          <CardTitle>Invite Team Member</CardTitle>
          <CardDescription>Send an invitation to join your workspace</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email">Email Address</Label>
            <Input
              id="email"
              type="email"
              placeholder="colleague@company.com"
              value={inviteEmail}
              onChange={(e) => onEmailChange(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="role">Role</Label>
            <select
              id="role"
              value={inviteRole}
              onChange={(e) => onRoleChange(e.target.value)}
              className="w-full p-2 rounded-lg border bg-background"
            >
              <option value="viewer">Viewer - Can view reports only</option>
              <option value="member">Member - Can run queries</option>
              <option value="editor">Editor - Can create reports</option>
              <option value="admin">Admin - Full access</option>
            </select>
          </div>
          <div className="flex justify-end gap-2 pt-4">
            <Button variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button onClick={onInvite} disabled={!inviteEmail}>
              <Mail className="w-4 h-4 mr-2" />
              Send Invitation
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
