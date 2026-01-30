"use client";

import { useState, useEffect } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { apiClient } from "@/lib/api";
import { motion } from "framer-motion";
import { useTheme } from "@/contexts/ThemeContext";
import { Settings, User, Bell, Shield, Palette } from "lucide-react";
import { ProfileSettings } from "@/components/settings/ProfileSettings";
import { NotificationSettings } from "@/components/settings/NotificationSettings";
import { SecuritySettings } from "@/components/settings/SecuritySettings";
import { AppearanceSettings } from "@/components/settings/AppearanceSettings";

export default function SettingsPage() {
  const { theme, setTheme, accentColor, setAccentColor } = useTheme();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [activeTab, setActiveTab] = useState("profile");

  const [profile, setProfile] = useState({
    name: "",
    email: "",
    avatar: "",
  });

  const [notifications, setNotifications] = useState({
    email: true,
    push: true,
    queryComplete: true,
    weeklyDigest: false,
  });

  useEffect(() => {
    loadSettings();
  }, []);

  const loadSettings = async () => {
    setLoading(true);
    try {
      const result = await apiClient.getUserSettings();
      if (result.data) {
        setProfile({
          name: result.data.name || "",
          email: result.data.email || "",
          avatar: result.data.avatar || "",
        });
        if (result.data.notifications) {
          setNotifications(result.data.notifications);
        }
      }
    } catch (error) {
      console.error("Failed to load settings:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      await apiClient.updateUserSettings({
        ...profile,
        notifications,
      });
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch (error) {
      console.error("Failed to save settings:", error);
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex-1 overflow-auto">
        <div className="container mx-auto max-w-4xl space-y-6 p-6">
          <div className="h-10 w-48 bg-muted/30 rounded animate-pulse" />
          <div className="h-96 bg-muted/30 rounded-lg animate-pulse" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 overflow-auto">
      <div className="container mx-auto max-w-4xl space-y-6 p-6">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <h1 className="text-3xl font-bold flex items-center gap-3">
            <Settings className="w-8 h-8 text-primary" />
            Settings
          </h1>
          <p className="text-muted-foreground mt-1">
            Manage your account settings and preferences
          </p>
        </motion.div>

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="profile" className="flex items-center gap-2">
              <User className="w-4 h-4" />
              Profile
            </TabsTrigger>
            <TabsTrigger value="notifications" className="flex items-center gap-2">
              <Bell className="w-4 h-4" />
              Notifications
            </TabsTrigger>
            <TabsTrigger value="security" className="flex items-center gap-2">
              <Shield className="w-4 h-4" />
              Security
            </TabsTrigger>
            <TabsTrigger value="appearance" className="flex items-center gap-2">
              <Palette className="w-4 h-4" />
              Appearance
            </TabsTrigger>
          </TabsList>

          <TabsContent value="profile" className="mt-6 space-y-6">
            <ProfileSettings
              profile={profile}
              onProfileChange={setProfile}
              onSave={handleSave}
              saving={saving}
              saved={saved}
            />
          </TabsContent>

          <TabsContent value="notifications" className="mt-6 space-y-6">
            <NotificationSettings
              notifications={notifications}
              onNotificationsChange={setNotifications}
            />
          </TabsContent>

          <TabsContent value="security" className="mt-6 space-y-6">
            <SecuritySettings />
          </TabsContent>

          <TabsContent value="appearance" className="mt-6 space-y-6">
            <AppearanceSettings
              theme={theme}
              accentColor={accentColor}
              onThemeChange={setTheme}
              onAccentColorChange={setAccentColor}
            />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
