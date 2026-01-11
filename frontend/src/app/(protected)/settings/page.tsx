"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { apiClient } from "@/lib/api";
import { motion } from "framer-motion";
import { useTheme } from "@/contexts/ThemeContext";
import {
  Settings,
  User,
  Bell,
  Shield,
  Palette,
  Globe,
  Key,
  Mail,
  Smartphone,
  Save,
  Check,
  Camera,
  Moon,
  Sun,
  Monitor
} from "lucide-react";

export default function SettingsPage() {
  const { theme, setTheme, accentColor, setAccentColor } = useTheme();
  const [settings, setSettings] = useState<any>(null);
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
      setSettings(result);
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
            <Card>
              <CardHeader>
                <CardTitle>Profile Information</CardTitle>
                <CardDescription>Update your personal information</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="flex items-center gap-6">
                  <div className="relative">
                    <div className="w-24 h-24 rounded-full bg-gradient-to-br from-primary/20 to-purple-500/20 flex items-center justify-center text-2xl font-semibold">
                      {profile.name ? profile.name.charAt(0).toUpperCase() : "?"}
                    </div>
                    <Button
                      variant="outline"
                      size="icon"
                      className="absolute bottom-0 right-0 rounded-full w-8 h-8"
                    >
                      <Camera className="w-4 h-4" />
                    </Button>
                  </div>
                  <div>
                    <h3 className="font-semibold">Profile Photo</h3>
                    <p className="text-sm text-muted-foreground">
                      JPG, GIF or PNG. Max size 2MB.
                    </p>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="name">Full Name</Label>
                    <Input
                      id="name"
                      value={profile.name}
                      onChange={(e) => setProfile({ ...profile, name: e.target.value })}
                      placeholder="John Doe"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="email">Email Address</Label>
                    <Input
                      id="email"
                      type="email"
                      value={profile.email}
                      onChange={(e) => setProfile({ ...profile, email: e.target.value })}
                      placeholder="john@example.com"
                    />
                  </div>
                </div>

                <div className="flex justify-end">
                  <Button onClick={handleSave} disabled={saving}>
                    {saving ? (
                      <>
                        <div className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin mr-2" />
                        Saving...
                      </>
                    ) : saved ? (
                      <>
                        <Check className="w-4 h-4 mr-2 text-green-500" />
                        Saved!
                      </>
                    ) : (
                      <>
                        <Save className="w-4 h-4 mr-2" />
                        Save Changes
                      </>
                    )}
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="notifications" className="mt-6 space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Notification Preferences</CardTitle>
                <CardDescription>Choose how you want to be notified</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between p-4 rounded-lg border">
                  <div className="flex items-center gap-3">
                    <Mail className="w-5 h-5 text-muted-foreground" />
                    <div>
                      <p className="font-medium">Email Notifications</p>
                      <p className="text-sm text-muted-foreground">
                        Receive updates via email
                      </p>
                    </div>
                  </div>
                  <Button
                    variant={notifications.email ? "default" : "outline"}
                    size="sm"
                    onClick={() => setNotifications({ ...notifications, email: !notifications.email })}
                  >
                    {notifications.email ? "Enabled" : "Disabled"}
                  </Button>
                </div>

                <div className="flex items-center justify-between p-4 rounded-lg border">
                  <div className="flex items-center gap-3">
                    <Smartphone className="w-5 h-5 text-muted-foreground" />
                    <div>
                      <p className="font-medium">Push Notifications</p>
                      <p className="text-sm text-muted-foreground">
                        Receive push notifications in browser
                      </p>
                    </div>
                  </div>
                  <Button
                    variant={notifications.push ? "default" : "outline"}
                    size="sm"
                    onClick={() => setNotifications({ ...notifications, push: !notifications.push })}
                  >
                    {notifications.push ? "Enabled" : "Disabled"}
                  </Button>
                </div>

                <div className="flex items-center justify-between p-4 rounded-lg border">
                  <div className="flex items-center gap-3">
                    <Bell className="w-5 h-5 text-muted-foreground" />
                    <div>
                      <p className="font-medium">Query Complete</p>
                      <p className="text-sm text-muted-foreground">
                        Get notified when queries finish processing
                      </p>
                    </div>
                  </div>
                  <Button
                    variant={notifications.queryComplete ? "default" : "outline"}
                    size="sm"
                    onClick={() => setNotifications({ ...notifications, queryComplete: !notifications.queryComplete })}
                  >
                    {notifications.queryComplete ? "Enabled" : "Disabled"}
                  </Button>
                </div>

                <div className="flex items-center justify-between p-4 rounded-lg border">
                  <div className="flex items-center gap-3">
                    <Mail className="w-5 h-5 text-muted-foreground" />
                    <div>
                      <p className="font-medium">Weekly Digest</p>
                      <p className="text-sm text-muted-foreground">
                        Receive a weekly summary of your analytics
                      </p>
                    </div>
                  </div>
                  <Button
                    variant={notifications.weeklyDigest ? "default" : "outline"}
                    size="sm"
                    onClick={() => setNotifications({ ...notifications, weeklyDigest: !notifications.weeklyDigest })}
                  >
                    {notifications.weeklyDigest ? "Enabled" : "Disabled"}
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="security" className="mt-6 space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Security Settings</CardTitle>
                <CardDescription>Manage your account security</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="flex items-center justify-between p-4 rounded-lg border">
                  <div className="flex items-center gap-3">
                    <Key className="w-5 h-5 text-muted-foreground" />
                    <div>
                      <p className="font-medium">Password</p>
                      <p className="text-sm text-muted-foreground">
                        Last changed 30 days ago
                      </p>
                    </div>
                  </div>
                  <Button variant="outline">Change Password</Button>
                </div>

                <div className="flex items-center justify-between p-4 rounded-lg border">
                  <div className="flex items-center gap-3">
                    <Shield className="w-5 h-5 text-muted-foreground" />
                    <div>
                      <p className="font-medium">Two-Factor Authentication</p>
                      <p className="text-sm text-muted-foreground">
                        Add an extra layer of security
                      </p>
                    </div>
                  </div>
                  <Badge variant="outline">Not Enabled</Badge>
                </div>

                <div className="flex items-center justify-between p-4 rounded-lg border">
                  <div className="flex items-center gap-3">
                    <Globe className="w-5 h-5 text-muted-foreground" />
                    <div>
                      <p className="font-medium">Active Sessions</p>
                      <p className="text-sm text-muted-foreground">
                        Manage your active sessions
                      </p>
                    </div>
                  </div>
                  <Button variant="outline">View Sessions</Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="appearance" className="mt-6 space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Appearance Settings</CardTitle>
                <CardDescription>Customize how InsightIQ looks</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="space-y-2">
                  <Label>Theme</Label>
                  <div className="grid grid-cols-3 gap-4">
                    <div
                      className={`p-4 rounded-lg border-2 cursor-pointer transition-all ${
                        theme === "dark" ? "border-primary bg-primary/5" : "border-border hover:border-primary/50"
                      }`}
                      onClick={() => setTheme("dark")}
                    >
                      <div className="h-16 rounded bg-slate-900 mb-2 flex items-center justify-center">
                        <Moon className="w-6 h-6 text-slate-400" />
                      </div>
                      <p className="font-medium text-sm text-center">Dark</p>
                    </div>
                    <div
                      className={`p-4 rounded-lg border-2 cursor-pointer transition-all ${
                        theme === "light" ? "border-primary bg-primary/5" : "border-border hover:border-primary/50"
                      }`}
                      onClick={() => setTheme("light")}
                    >
                      <div className="h-16 rounded bg-white border mb-2 flex items-center justify-center">
                        <Sun className="w-6 h-6 text-amber-500" />
                      </div>
                      <p className="font-medium text-sm text-center">Light</p>
                    </div>
                    <div
                      className={`p-4 rounded-lg border-2 cursor-pointer transition-all ${
                        theme === "system" ? "border-primary bg-primary/5" : "border-border hover:border-primary/50"
                      }`}
                      onClick={() => setTheme("system")}
                    >
                      <div className="h-16 rounded bg-gradient-to-br from-slate-900 to-white mb-2 flex items-center justify-center">
                        <Monitor className="w-6 h-6 text-slate-600" />
                      </div>
                      <p className="font-medium text-sm text-center">System</p>
                    </div>
                  </div>
                </div>

                <div className="space-y-2">
                  <Label>Accent Color</Label>
                  <div className="flex gap-3">
                    {[
                      { name: "cyan", color: "#4fd1c5", bg: "bg-teal-400" },
                      { name: "indigo", color: "#6366f1", bg: "bg-indigo-500" },
                      { name: "purple", color: "#8b5cf6", bg: "bg-purple-500" },
                      { name: "pink", color: "#ec4899", bg: "bg-pink-500" },
                      { name: "orange", color: "#f59e0b", bg: "bg-amber-500" }
                    ].map((item) => (
                      <button
                        key={item.name}
                        onClick={() => setAccentColor(item.name as any)}
                        className={`w-12 h-12 rounded-full border-2 transition-all ${
                          accentColor === item.name
                            ? "border-white scale-110 shadow-lg"
                            : "border-transparent hover:scale-105"
                        }`}
                        style={{ backgroundColor: item.color }}
                        title={item.name.charAt(0).toUpperCase() + item.name.slice(1)}
                      />
                    ))}
                  </div>
                  <p className="text-sm text-muted-foreground">
                    Current: <span className="capitalize">{accentColor}</span>
                  </p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
