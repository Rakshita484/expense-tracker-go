import React, { useState, useEffect } from 'react';
import { groupService, userService } from '@/services/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select } from '@/components/ui/select';
import { Users, FolderPlus, UserPlus, Users2, Loader2 } from 'lucide-react';
import toast from 'react-hot-toast';

export default function Groups() {
    const [groups, setGroups] = useState([]);
    const [users, setUsers] = useState([]);
    const [groupMembers, setGroupMembers] = useState({}); // { groupId: [members] }
    const [loading, setLoading] = useState(true);
    const [isAdding, setIsAdding] = useState(false);

    // Create Group Form
    const [name, setName] = useState('');
    const [submitting, setSubmitting] = useState(false);

    // Add Member Form
    const [selectedGroup, setSelectedGroup] = useState(null);
    const [selectedUser, setSelectedUser] = useState('');

    useEffect(() => {
        fetchInitialData();
    }, []);

    const fetchInitialData = async () => {
        try {
            setLoading(true);
            const [fetchedUsers, fetchedGroups] = await Promise.all([
                userService.getUsers(),
                groupService.getGroups()
            ]);
            setUsers(fetchedUsers || []);
            setGroups(fetchedGroups || []);

            // Fetch members for each group
            const membersMap = {};
            for (const group of (fetchedGroups || [])) {
                try {
                    const members = await groupService.getMembers(group.id);
                    membersMap[group.id] = members || [];
                } catch {
                    membersMap[group.id] = [];
                }
            }
            setGroupMembers(membersMap);
        } catch (error) {
            console.error(error);
            toast.error('Failed to load groups');
        } finally {
            setLoading(false);
        }
    };

    const handleCreateGroup = async (e) => {
        e.preventDefault();
        if (!name) return toast.error("Group name is required");

        setSubmitting(true);
        try {
            await groupService.createGroup({ name });
            toast.success('Group created successfully');
            setName('');
            setIsAdding(false);
            fetchInitialData(); // Refresh to get updated list
        } catch (error) {
            toast.error('Failed to create group');
        } finally {
            setSubmitting(false);
        }
    };

    const handleAddMember = async (e, groupId) => {
        e.preventDefault();
        if (!selectedUser) return toast.error("Please select a user");

        try {
            await groupService.addMember(groupId, parseInt(selectedUser));
            toast.success('Member added successfully');
            setSelectedGroup(null);
            setSelectedUser('');
            fetchInitialData(); // Refresh to show new member
        } catch (error) {
            toast.error(error.response?.data?.error || 'Failed to add member. They might already be in the group.');
        }
    };

    return (
        <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-2xl font-bold tracking-tight text-slate-900">Groups</h2>
                    <p className="text-muted-foreground mt-1">Create groups to organize shared expenses.</p>
                </div>
                <Button onClick={() => setIsAdding(!isAdding)} className="gap-2 bg-emerald-600 hover:bg-emerald-700">
                    {isAdding ? 'Cancel' : <><FolderPlus size={16} /> Create Group</>}
                </Button>
            </div>

            {isAdding && (
                <Card className="border-emerald-500/20 bg-emerald-500/5">
                    <CardHeader>
                        <CardTitle className="text-lg">Create New Group</CardTitle>
                        <CardDescription>Give your group a name (e.g. "Apartment 4B", "Road Trip").</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleCreateGroup} className="flex gap-4 max-w-md items-end">
                            <div className="space-y-2 flex-1">
                                <Label htmlFor="name">Group Name</Label>
                                <Input
                                    id="name"
                                    placeholder="e.g. Weekend Trip"
                                    value={name}
                                    onChange={e => setName(e.target.value)}
                                    required
                                />
                            </div>
                            <Button type="submit" disabled={submitting} className="bg-emerald-600 hover:bg-emerald-700">
                                {submitting ? 'Creating...' : 'Create'}
                            </Button>
                        </form>
                    </CardContent>
                </Card>
            )}

            {loading ? (
                <div className="flex items-center justify-center py-20 text-slate-400">
                    <Loader2 className="animate-spin mr-2" size={20} />
                    Loading groups...
                </div>
            ) : groups.length === 0 ? (
                <div className="text-center py-16 text-slate-400 bg-slate-50 border border-dashed rounded-xl">
                    <FolderPlus className="h-10 w-10 mx-auto mb-3 text-slate-300" />
                    <p className="text-sm font-medium">No groups yet</p>
                    <p className="text-xs mt-1">Create your first group to start tracking shared expenses</p>
                </div>
            ) : (
                <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                    {groups.map(group => {
                        const members = groupMembers[group.id] || [];
                        return (
                            <Card key={group.id} className="relative overflow-hidden hover:shadow-lg transition-all duration-300 border-t-4 border-t-emerald-500">
                                <CardHeader className="pb-3">
                                    <div className="flex justify-between items-start">
                                        <div>
                                            <CardTitle className="text-xl">{group.name}</CardTitle>
                                            <CardDescription className="mt-1 flex items-center gap-1">
                                                <Users2 size={14} /> Group #{group.id}
                                            </CardDescription>
                                        </div>
                                        <div className="bg-emerald-100 text-emerald-700 px-2.5 py-1 rounded-full text-xs font-semibold flex items-center gap-1">
                                            <Users size={12} />
                                            {group.member_count || members.length} members
                                        </div>
                                    </div>
                                </CardHeader>

                                <CardContent>
                                    {/* Member Avatars */}
                                    {members.length > 0 && (
                                        <div className="mb-4">
                                            <p className="text-xs font-medium text-slate-500 uppercase tracking-wider mb-2">Members</p>
                                            <div className="flex flex-wrap gap-2">
                                                {members.map(member => (
                                                    <div
                                                        key={member.id}
                                                        className="flex items-center gap-1.5 bg-slate-50 border px-2.5 py-1.5 rounded-full text-xs"
                                                        title={member.email}
                                                    >
                                                        <div className="h-5 w-5 rounded-full bg-emerald-100 flex items-center justify-center text-emerald-700 font-bold text-[10px]">
                                                            {member.name?.charAt(0)}
                                                        </div>
                                                        <span className="font-medium text-slate-700">{member.name}</span>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    )}

                                    {selectedGroup === group.id ? (
                                        <form onSubmit={(e) => handleAddMember(e, group.id)} className="space-y-3 mt-2">
                                            <div className="space-y-2">
                                                <Label className="text-xs text-slate-500 uppercase tracking-wider">Add existing user to group</Label>
                                                <Select value={selectedUser} onChange={e => setSelectedUser(e.target.value)} required>
                                                    <option value="" disabled>Select a user...</option>
                                                    {users.map(u => (
                                                        <option key={u.id} value={u.id}>{u.name} ({u.email})</option>
                                                    ))}
                                                </Select>
                                            </div>
                                            <div className="flex gap-2">
                                                <Button type="submit" size="sm" className="w-full">Add Member</Button>
                                                <Button type="button" variant="outline" size="sm" onClick={() => setSelectedGroup(null)}>Cancel</Button>
                                            </div>
                                        </form>
                                    ) : (
                                        <div className="pt-4 border-t flex justify-between items-center">
                                            <span className="text-sm font-medium text-slate-500">Manage Members</span>
                                            <Button variant="ghost" size="sm" className="h-8 gap-1 text-emerald-600 hover:text-emerald-700 hover:bg-emerald-50" onClick={() => setSelectedGroup(group.id)}>
                                                <UserPlus size={14} /> Add User
                                            </Button>
                                        </div>
                                    )}
                                </CardContent>
                            </Card>
                        );
                    })}
                </div>
            )}
        </div>
    );
}
