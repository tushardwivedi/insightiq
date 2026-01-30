"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { apiClient } from "@/lib/api";
import { motion } from "framer-motion";
import {
  Users,
  Search,
  Mail,
  UserPlus,
  CheckCircle,
} from "lucide-react";
import { TeamMemberRow } from "@/components/team/TeamMemberRow";
import { InviteMemberModal } from "@/components/team/InviteMemberModal";

interface TeamMember {
  id: string;
  email: string;
  name?: string;
  role: string;
  status: "active" | "pending" | "inactive";
  created_at: string;
  avatar?: string;
}

export default function TeamPage() {
  const [members, setMembers] = useState<TeamMember[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [showInvite, setShowInvite] = useState(false);
  const [inviteEmail, setInviteEmail] = useState("");
  const [inviteRole, setInviteRole] = useState("member");

  useEffect(() => {
    loadMembers();
  }, []);

  const loadMembers = async () => {
    setLoading(true);
    try {
      const result = await apiClient.getTeamMembers();
      setMembers(result);
    } catch (error) {
      console.error("Failed to load team members:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleInvite = async () => {
    if (!inviteEmail) return;
    try {
      await apiClient.inviteTeamMember(inviteEmail, inviteRole);
      setShowInvite(false);
      setInviteEmail("");
      loadMembers();
    } catch (error) {
      console.error("Failed to invite:", error);
    }
  };

  const handleRemove = async (id: string) => {
    if (!confirm("Are you sure you want to remove this team member?")) return;
    try {
      await apiClient.removeTeamMember(id);
      setMembers(prev => prev.filter(m => m.id !== id));
    } catch (error) {
      console.error("Failed to remove:", error);
    }
  };

  const filteredMembers = members.filter(m =>
    m.email?.toLowerCase().includes(search.toLowerCase()) ||
    m.name?.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="flex-1 overflow-auto">
      <div className="container mx-auto max-w-6xl space-y-6 p-6">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="flex items-center justify-between"
        >
          <div>
            <h1 className="text-3xl font-bold flex items-center gap-3">
              <Users className="w-8 h-8 text-primary" />
              Team
            </h1>
            <p className="text-muted-foreground mt-1">
              Manage your team members and their access
            </p>
          </div>
          <Button className="bg-gradient-to-r from-primary to-cyan-500" onClick={() => setShowInvite(true)}>
            <UserPlus className="w-4 h-4 mr-2" />
            Invite Member
          </Button>
        </motion.div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-500/20 flex items-center justify-center">
                  <Users className="w-6 h-6 text-primary" />
                </div>
                <div>
                  <p className="text-2xl font-bold">{members.length}</p>
                  <p className="text-sm text-muted-foreground">Total Members</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center">
                  <CheckCircle className="w-6 h-6 text-green-500" />
                </div>
                <div>
                  <p className="text-2xl font-bold">
                    {members.filter(m => m.status === "active").length}
                  </p>
                  <p className="text-sm text-muted-foreground">Active</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-orange-500/20 to-red-500/20 flex items-center justify-center">
                  <Mail className="w-6 h-6 text-orange-500" />
                </div>
                <div>
                  <p className="text-2xl font-bold">
                    {members.filter(m => m.status === "pending").length}
                  </p>
                  <p className="text-sm text-muted-foreground">Pending</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Team Members</CardTitle>
                <CardDescription>All members in your workspace</CardDescription>
              </div>
              <div className="relative">
                <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Search members..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-9 w-[250px]"
                />
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-20 bg-muted/30 rounded-lg animate-pulse" />
                ))}
              </div>
            ) : filteredMembers.length === 0 ? (
              <div className="text-center py-12">
                <Users className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <p className="text-muted-foreground mb-4">No team members yet</p>
                <Button onClick={() => setShowInvite(true)}>
                  <UserPlus className="w-4 h-4 mr-2" />
                  Invite Your First Member
                </Button>
              </div>
            ) : (
              <div className="space-y-3">
                {filteredMembers.map((member, index) => (
                  <TeamMemberRow
                    key={member.id}
                    member={member}
                    index={index}
                    onRemove={handleRemove}
                  />
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {showInvite && (
          <InviteMemberModal
            inviteEmail={inviteEmail}
            inviteRole={inviteRole}
            onEmailChange={setInviteEmail}
            onRoleChange={setInviteRole}
            onInvite={handleInvite}
            onClose={() => setShowInvite(false)}
          />
        )}
      </div>
    </div>
  );
}
