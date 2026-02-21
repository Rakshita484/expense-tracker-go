import React, { useState } from 'react';
import { NavLink, Outlet, useLocation } from 'react-router-dom';
import { LayoutDashboard, Users, FolderKanban, Receipt, ArrowRightLeft, Menu, X } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';

export default function Layout() {
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const location = useLocation();

    const navigation = [
        { name: 'Dashboard', href: '/', icon: LayoutDashboard },
        { name: 'Users', href: '/users', icon: Users },
        { name: 'Groups', href: '/groups', icon: FolderKanban },
        { name: 'Expenses', href: '/expenses', icon: Receipt },
        { name: 'Settlements', href: '/settlements', icon: ArrowRightLeft },
    ];

    return (
        <div className="flex h-screen overflow-hidden bg-slate-50">

            {/* Mobile sidebar backdrop */}
            {sidebarOpen && (
                <div
                    className="fixed inset-0 z-40 bg-black/50 lg:hidden"
                    onClick={() => setSidebarOpen(false)}
                />
            )}

            {/* Sidebar */}
            <div
                className={cn(
                    "fixed inset-y-0 left-0 z-50 w-64 transform bg-white border-r shadow-sm transition-transform duration-200 ease-in-out lg:static lg:translate-x-0",
                    sidebarOpen ? "translate-x-0" : "-translate-x-full"
                )}
            >
                <div className="flex h-16 shrink-0 items-center px-6 border-b">
                    <Receipt className="h-6 w-6 text-primary mr-2" />
                    <span className="text-xl font-bold tracking-tight text-slate-900">Splitwise<span className="text-primary">Clone</span></span>
                </div>

                <nav className="flex flex-1 flex-col overflow-y-auto px-4 py-4">
                    <ul className="space-y-1">
                        {navigation.map((item) => {
                            const isActive = location.pathname === item.href;
                            return (
                                <li key={item.name}>
                                    <NavLink
                                        to={item.href}
                                        onClick={() => setSidebarOpen(false)}
                                        className={cn(
                                            "group flex gap-x-3 rounded-md p-2.5 text-sm leading-6 font-semibold transition-all",
                                            isActive
                                                ? "bg-primary/10 text-primary"
                                                : "text-slate-700 hover:text-primary hover:bg-slate-50"
                                        )}
                                    >
                                        <item.icon
                                            className={cn(
                                                "h-5 w-5 shrink-0",
                                                isActive ? "text-primary" : "text-slate-400 group-hover:text-primary"
                                            )}
                                        />
                                        {item.name}
                                    </NavLink>
                                </li>
                            )
                        })}
                    </ul>
                </nav>
            </div>

            {/* Main content area */}
            <div className="flex flex-1 flex-col w-full h-full overflow-hidden">

                {/* Header */}
                <header className="flex h-16 shrink-0 items-center gap-x-4 border-b bg-white px-4 shadow-sm sm:gap-x-6 sm:px-6 lg:px-8">
                    <Button
                        variant="ghost"
                        size="icon"
                        className="-m-2.5 p-2.5 text-slate-700 lg:hidden"
                        onClick={() => setSidebarOpen(true)}
                    >
                        <span className="sr-only">Open sidebar</span>
                        <Menu className="h-6 w-6" aria-hidden="true" />
                    </Button>

                    <div className="flex flex-1 gap-x-4 self-stretch lg:gap-x-6 items-center">
                        <h1 className="text-lg font-semibold leading-6 text-slate-900 ml-auto lg:ml-0">
                            {navigation.find(n => n.href === location.pathname)?.name || 'App'}
                        </h1>
                    </div>
                </header>

                {/* Main routing area */}
                <main className="flex-1 overflow-y-auto bg-slate-50/50 p-4 sm:p-6 lg:p-8">
                    <div className="mx-auto max-w-5xl h-full pb-10">
                        <Outlet />
                    </div>
                </main>
            </div>
        </div>
    );
}
