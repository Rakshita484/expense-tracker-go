import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { Users, FolderKanban, Receipt, ArrowRightLeft, TrendingUp, DollarSign, Clock } from 'lucide-react';
import { dashboardService, groupService } from '@/services/api';

export default function Dashboard() {
    const [stats, setStats] = useState(null);
    const [recentExpenses, setRecentExpenses] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchDashboardData();
    }, []);

    const fetchDashboardData = async () => {
        try {
            setLoading(true);
            const [statsData, expensesData] = await Promise.all([
                dashboardService.getStats(),
                groupService.getExpenses(1).catch(() => [])
            ]);
            setStats(statsData);
            setRecentExpenses((expensesData || []).slice(0, 5));
        } catch (error) {
            console.error('Failed to load dashboard data', error);
        } finally {
            setLoading(false);
        }
    };

    const statCards = stats ? [
        {
            title: "Total Users",
            value: stats.total_users,
            icon: Users,
            color: "text-blue-600",
            bg: "bg-blue-100/80",
            borderColor: "border-l-blue-500"
        },
        {
            title: "Total Groups",
            value: stats.total_groups,
            icon: FolderKanban,
            color: "text-emerald-600",
            bg: "bg-emerald-100/80",
            borderColor: "border-l-emerald-500"
        },
        {
            title: "Total Expenses",
            value: stats.total_expenses,
            icon: Receipt,
            color: "text-amber-600",
            bg: "bg-amber-100/80",
            borderColor: "border-l-amber-500"
        },
        {
            title: "Total Spent",
            value: `$${parseFloat(stats.total_spent).toFixed(2)}`,
            icon: DollarSign,
            color: "text-purple-600",
            bg: "bg-purple-100/80",
            borderColor: "border-l-purple-500"
        }
    ] : [];

    const formatTimeAgo = (dateStr) => {
        const date = new Date(dateStr);
        const now = new Date();
        const diff = Math.floor((now - date) / 1000);
        if (diff < 60) return 'just now';
        if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
        if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
        return `${Math.floor(diff / 86400)}d ago`;
    };

    return (
        <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <div>
                <h2 className="text-2xl font-bold tracking-tight text-slate-900">Dashboard</h2>
                <p className="text-muted-foreground mt-1">
                    Overview of your shared expenses at a glance.
                </p>
            </div>

            {/* Stats Cards */}
            {loading ? (
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                    {[1, 2, 3, 4].map(i => (
                        <Card key={i} className="animate-pulse">
                            <CardContent className="p-6">
                                <div className="h-16 bg-slate-100 rounded" />
                            </CardContent>
                        </Card>
                    ))}
                </div>
            ) : (
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                    {statCards.map((card) => (
                        <Card key={card.title} className={`border-l-4 ${card.borderColor} hover:shadow-md transition-shadow`}>
                            <CardContent className="p-6">
                                <div className="flex items-center justify-between">
                                    <div>
                                        <p className="text-xs font-medium text-slate-500 uppercase tracking-wider">{card.title}</p>
                                        <p className="text-2xl font-bold text-slate-900 mt-1">{card.value}</p>
                                    </div>
                                    <div className={`p-3 rounded-xl ${card.bg}`}>
                                        <card.icon className={`h-5 w-5 ${card.color}`} />
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}

            {/* Recent Activity and Quick Actions */}
            <div className="grid gap-6 lg:grid-cols-3">

                {/* Recent Activity Feed */}
                <Card className="lg:col-span-2 border-t-4 border-t-slate-800">
                    <CardHeader className="pb-3">
                        <div className="flex items-center gap-2">
                            <div className="p-2 bg-slate-100 rounded-lg">
                                <Clock size={18} className="text-slate-600" />
                            </div>
                            <div>
                                <CardTitle className="text-lg">Recent Activity</CardTitle>
                                <CardDescription>Latest expenses from your groups</CardDescription>
                            </div>
                        </div>
                    </CardHeader>
                    <CardContent>
                        {recentExpenses.length === 0 ? (
                            <div className="text-center py-8 text-slate-400 bg-slate-50 border border-dashed rounded-xl">
                                <Receipt className="h-8 w-8 mx-auto mb-2 text-slate-300" />
                                <p className="text-sm font-medium">No expenses yet</p>
                                <p className="text-xs mt-1">Add your first expense to see activity here</p>
                            </div>
                        ) : (
                            <div className="space-y-3">
                                {recentExpenses.map(expense => (
                                    <div key={expense.id} className="flex items-center gap-4 p-3 rounded-xl bg-slate-50/80 border border-slate-100 hover:bg-slate-100/80 transition-colors">
                                        <div className="h-10 w-10 rounded-full bg-amber-100 flex items-center justify-center text-amber-700 font-bold text-sm shrink-0">
                                            {expense.paid_by?.name?.charAt(0) || '?'}
                                        </div>
                                        <div className="flex-1 min-w-0">
                                            <p className="text-sm font-semibold text-slate-900 truncate">{expense.description}</p>
                                            <p className="text-xs text-slate-500">
                                                <span className="font-medium text-amber-700">{expense.paid_by?.name || 'Unknown'}</span> paid • split among {expense.splits?.length || 0} people
                                            </p>
                                        </div>
                                        <div className="text-right shrink-0">
                                            <p className="font-mono text-sm font-bold text-slate-900">${parseFloat(expense.amount).toFixed(2)}</p>
                                            <p className="text-xs text-slate-400">{formatTimeAgo(expense.created_at)}</p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </CardContent>
                </Card>

                {/* Quick Actions */}
                <Card className="border-t-4 border-t-blue-500">
                    <CardHeader className="pb-3">
                        <CardTitle className="text-lg">Quick Actions</CardTitle>
                        <CardDescription>Jump to common tasks</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-3">
                        <Button asChild variant="outline" className="w-full justify-start gap-3 h-12 text-sm">
                            <Link to="/expenses">
                                <div className="p-1.5 bg-amber-100 rounded-md">
                                    <Receipt size={14} className="text-amber-600" />
                                </div>
                                Add New Expense
                            </Link>
                        </Button>
                        <Button asChild variant="outline" className="w-full justify-start gap-3 h-12 text-sm">
                            <Link to="/users">
                                <div className="p-1.5 bg-blue-100 rounded-md">
                                    <Users size={14} className="text-blue-600" />
                                </div>
                                Manage Users
                            </Link>
                        </Button>
                        <Button asChild variant="outline" className="w-full justify-start gap-3 h-12 text-sm">
                            <Link to="/groups">
                                <div className="p-1.5 bg-emerald-100 rounded-md">
                                    <FolderKanban size={14} className="text-emerald-600" />
                                </div>
                                View Groups
                            </Link>
                        </Button>
                        <Button asChild variant="outline" className="w-full justify-start gap-3 h-12 text-sm">
                            <Link to="/settlements">
                                <div className="p-1.5 bg-purple-100 rounded-md">
                                    <ArrowRightLeft size={14} className="text-purple-600" />
                                </div>
                                Calculate Settlements
                            </Link>
                        </Button>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
