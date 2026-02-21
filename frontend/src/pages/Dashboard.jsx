import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { Users, FolderKanban, Receipt, ArrowRightLeft } from 'lucide-react';

export default function Dashboard() {
    const cards = [
        {
            title: "Users",
            description: "Manage people in the system",
            icon: Users,
            href: "/users",
            color: "text-blue-500",
            bg: "bg-blue-100/50"
        },
        {
            title: "Groups",
            description: "Create and manage expense groups",
            icon: FolderKanban,
            href: "/groups",
            color: "text-emerald-500",
            bg: "bg-emerald-100/50"
        },
        {
            title: "Expenses",
            description: "Add new expenses and view balances",
            icon: Receipt,
            href: "/expenses",
            color: "text-amber-500",
            bg: "bg-amber-100/50"
        },
        {
            title: "Settlements",
            description: "Calculate optimal debt payments",
            icon: ArrowRightLeft,
            href: "/settlements",
            color: "text-purple-500",
            bg: "bg-purple-100/50"
        }
    ];

    return (
        <div className="space-y-6">
            <div>
                <h2 className="text-2xl font-bold tracking-tight text-slate-900">Welcome to Splitwise Clone</h2>
                <p className="text-muted-foreground mt-2">
                    Manage your shared expenses, calculate balances, and settle debts efficiently.
                </p>
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {cards.map((card) => (
                    <Card key={card.title} className="hover:shadow-md transition-shadow">
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium">
                                {card.title}
                            </CardTitle>
                            <div className={`p-2 rounded-full ${card.bg}`}>
                                <card.icon className={`h-4 w-4 ${card.color}`} />
                            </div>
                        </CardHeader>
                        <CardContent>
                            <p className="text-xs text-muted-foreground mb-4">
                                {card.description}
                            </p>
                            <Button asChild variant="outline" className="w-full h-8 text-xs">
                                <Link to={card.href}>Go to {card.title} →</Link>
                            </Button>
                        </CardContent>
                    </Card>
                ))}
            </div>

            <Card className="mt-8 border-dashed bg-slate-50 border-slate-300">
                <CardContent className="flex flex-col items-center justify-center p-12 text-center text-slate-500 h-64">
                    <div className="rounded-full bg-slate-100 p-3 mb-4">
                        <Receipt className="h-6 w-6 text-slate-400" />
                    </div>
                    <p className="text-sm font-medium">Get started by creating some users</p>
                    <p className="text-xs mt-1 text-slate-400">Head over to the Users tab to begin adding people.</p>
                    <Button asChild className="mt-6" size="sm">
                        <Link to="/users">Manage Users</Link>
                    </Button>
                </CardContent>
            </Card>
        </div>
    );
}
