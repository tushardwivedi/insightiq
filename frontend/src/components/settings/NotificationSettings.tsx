"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Mail, Smartphone, Bell } from "lucide-react";
import { LucideIcon } from "lucide-react";

interface NotificationState {
  email: boolean;
  push: boolean;
  queryComplete: boolean;
  weeklyDigest: boolean;
}

interface NotificationSettingsProps {
  notifications: NotificationState;
  onNotificationsChange: (notifications: NotificationState) => void;
}

interface NotificationRowProps {
  icon: LucideIcon;
  title: string;
  description: string;
  enabled: boolean;
  onToggle: () => void;
}

function NotificationRow({ icon: Icon, title, description, enabled, onToggle }: NotificationRowProps) {
  return (
    <div className="flex items-center justify-between p-4 rounded-lg border">
      <div className="flex items-center gap-3">
        <Icon className="w-5 h-5 text-muted-foreground" />
        <div>
          <p className="font-medium">{title}</p>
          <p className="text-sm text-muted-foreground">{description}</p>
        </div>
      </div>
      <Button
        variant={enabled ? "default" : "outline"}
        size="sm"
        onClick={onToggle}
      >
        {enabled ? "Enabled" : "Disabled"}
      </Button>
    </div>
  );
}

export function NotificationSettings({ notifications, onNotificationsChange }: NotificationSettingsProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Notification Preferences</CardTitle>
        <CardDescription>Choose how you want to be notified</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <NotificationRow
          icon={Mail}
          title="Email Notifications"
          description="Receive updates via email"
          enabled={notifications.email}
          onToggle={() => onNotificationsChange({ ...notifications, email: !notifications.email })}
        />
        <NotificationRow
          icon={Smartphone}
          title="Push Notifications"
          description="Receive push notifications in browser"
          enabled={notifications.push}
          onToggle={() => onNotificationsChange({ ...notifications, push: !notifications.push })}
        />
        <NotificationRow
          icon={Bell}
          title="Query Complete"
          description="Get notified when queries finish processing"
          enabled={notifications.queryComplete}
          onToggle={() => onNotificationsChange({ ...notifications, queryComplete: !notifications.queryComplete })}
        />
        <NotificationRow
          icon={Mail}
          title="Weekly Digest"
          description="Receive a weekly summary of your analytics"
          enabled={notifications.weeklyDigest}
          onToggle={() => onNotificationsChange({ ...notifications, weeklyDigest: !notifications.weeklyDigest })}
        />
      </CardContent>
    </Card>
  );
}
