"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Save, Check, Camera } from "lucide-react";

interface ProfileSettingsProps {
  profile: { name: string; email: string; avatar: string };
  onProfileChange: (profile: { name: string; email: string; avatar: string }) => void;
  onSave: () => void;
  saving: boolean;
  saved: boolean;
}

export function ProfileSettings({ profile, onProfileChange, onSave, saving, saved }: ProfileSettingsProps) {
  return (
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
              onChange={(e) => onProfileChange({ ...profile, name: e.target.value })}
              placeholder="John Doe"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="email">Email Address</Label>
            <Input
              id="email"
              type="email"
              value={profile.email}
              onChange={(e) => onProfileChange({ ...profile, email: e.target.value })}
              placeholder="john@example.com"
            />
          </div>
        </div>

        <div className="flex justify-end">
          <Button onClick={onSave} disabled={saving}>
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
  );
}
