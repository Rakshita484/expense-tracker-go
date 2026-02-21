import React, { useState, useEffect } from 'react';
import { userService } from '@/services/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { UserPlus, UserRound, Mail } from 'lucide-react';
import toast from 'react-hot-toast';

export default function Users() {
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [isAdding, setIsAdding] = useState(false);

    // Form state
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        fetchUsers();
    }, []);

    const fetchUsers = async () => {
        try {
            setLoading(true);
            const data = await userService.getUsers();
            setUsers(data || []);
        } catch (error) {
            toast.error('Failed to load users');
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleCreateUser = async (e) => {
        e.preventDefault();
        if (!name || !email) return toast.error("Please fill all fields");

        setSubmitting(true);
        try {
            const newUser = await userService.createUser({ name, email });
            setUsers(prev => [...prev, newUser]);
            toast.success('User created successfully');
            setName('');
            setEmail('');
            setIsAdding(false);
        } catch (error) {
            toast.error('Failed to create user. Email might be duplicate.');
        } finally {
            setSubmitting(false);
        }
    };

    return (
        <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-2xl font-bold tracking-tight text-slate-900">Users</h2>
                    <p className="text-muted-foreground mt-1">Manage people participating in the expense tracker.</p>
                </div>
                <Button onClick={() => setIsAdding(!isAdding)} className="gap-2">
                    {isAdding ? 'Cancel' : <><UserPlus size={16} /> Add User</>}
                </Button>
            </div>

            {isAdding && (
                <Card className="border-primary/20 bg-primary/5">
                    <CardHeader>
                        <CardTitle className="text-lg">Create New User</CardTitle>
                        <CardDescription>Add someone to start splitting bills with them.</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleCreateUser} className="space-y-4 max-w-md">
                            <div className="space-y-2">
                                <Label htmlFor="name">Full Name</Label>
                                <Input
                                    id="name"
                                    placeholder="e.g. Alice Johnson"
                                    value={name}
                                    onChange={e => setName(e.target.value)}
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="email">Email Address</Label>
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="e.g. alice@example.com"
                                    value={email}
                                    onChange={e => setEmail(e.target.value)}
                                    required
                                />
                            </div>
                            <Button type="submit" disabled={submitting}>
                                {submitting ? 'Creating...' : 'Create User'}
                            </Button>
                        </form>
                    </CardContent>
                </Card>
            )}

            {loading ? (
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                    {[1, 2, 3].map(i => (
                        <Card key={i} className="animate-pulse shadow-sm h-32"></Card>
                    ))}
                </div>
            ) : users.length === 0 ? (
                <Card className="border-dashed h-40 flex items-center justify-center text-slate-500">
                    No users found. Create one to get started!
                </Card>
            ) : (
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                    {users.map(user => (
                        <Card key={user.id} className="overflow-hidden hover:shadow-md transition-shadow duration-200 border-l-4 border-l-primary/60">
                            <CardContent className="p-0">
                                <div className="flex items-center gap-4 p-5">
                                    <div className="h-12 w-12 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 flex-shrink-0">
                                        <UserRound size={24} />
                                    </div>
                                    <div className="min-w-0 flex-1">
                                        <p className="text-sm font-semibold text-slate-900 truncate">
                                            {user.name}
                                        </p>
                                        <div className="flex items-center gap-1.5 mt-1 text-slate-500">
                                            <Mail size={12} />
                                            <p className="text-xs truncate">{user.email}</p>
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
}
