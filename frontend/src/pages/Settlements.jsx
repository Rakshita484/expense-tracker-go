import React, { useState, useEffect } from 'react';
import { groupService } from '@/services/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { ArrowRightLeft, TrendingDown, HandCoins, CheckCircle2 } from 'lucide-react';
import toast from 'react-hot-toast';

export default function Settlements() {
    const [settlements, setSettlements] = useState([]);
    const [loading, setLoading] = useState(true);
    const groupId = 1; // Assuming default seeded group

    useEffect(() => {
        fetchSettlements();
    }, []);

    const fetchSettlements = async () => {
        try {
            setLoading(true);
            const data = await groupService.getSettlements(groupId);
            setSettlements(data || []);
        } catch (error) {
            toast.error('Failed to calculate settlements');
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="max-w-4xl mx-auto space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">

            <div className="text-center space-y-2 mb-10">
                <div className="mx-auto w-16 h-16 bg-purple-100 text-purple-600 rounded-full flex items-center justify-center mb-4 shadow-sm">
                    <ArrowRightLeft size={32} />
                </div>
                <h2 className="text-3xl font-bold tracking-tight text-slate-900">Optimal Settlements</h2>
                <p className="text-slate-500 max-w-xl mx-auto">
                    Our greedy algorithm has calculated the minimum number of transactions needed for everyone to get paid back exactly what they are owed.
                </p>
            </div>

            {loading ? (
                <div className="space-y-4">
                    {[1, 2].map(i => <Card key={i} className="animate-pulse h-24 bg-slate-50"></Card>)}
                </div>
            ) : settlements.length === 0 ? (
                <Card className="border-dashed bg-slate-50">
                    <CardContent className="flex flex-col items-center justify-center py-16 text-center">
                        <CheckCircle2 size={48} className="text-emerald-500 mb-4" />
                        <h3 className="text-xl font-semibold text-slate-900">You're all settled up!</h3>
                        <p className="text-slate-500 mt-2">No one owes anything in this group. Awesome.</p>
                    </CardContent>
                </Card>
            ) : (
                <div className="space-y-4">
                    <div className="flex items-center justify-between text-sm font-medium text-slate-500 px-4 mb-2">
                        <span>WHO PAYS</span>
                        <span>TRANSACTION</span>
                        <span>WHO RECEIVES</span>
                    </div>

                    {settlements.map((tx, idx) => (
                        <Card key={idx} className="overflow-hidden hover:shadow-md transition-all duration-300 border-l-4 border-l-purple-500">
                            <CardContent className="p-0">
                                <div className="flex flex-col sm:flex-row items-center p-6 gap-6 relative">

                                    {/* FROM */}
                                    <div className="flex-1 flex flex-col items-center sm:items-end text-center sm:text-right w-full">
                                        <span className="text-xs text-red-500 font-bold tracking-wider mb-1 flex items-center gap-1 uppercase">
                                            <TrendingDown size={14} /> Sender
                                        </span>
                                        <span className="text-lg font-semibold text-slate-900">{tx.from_user_name}</span>
                                    </div>

                                    {/* ARROW & AMOUNT */}
                                    <div className="flex flex-col items-center px-4 shrink-0">
                                        <div className="bg-purple-100 text-purple-700 font-mono text-xl font-bold py-2 px-6 rounded-full border border-purple-200 shadow-sm z-10">
                                            ${tx.amount}
                                        </div>
                                        {/* Visual connecting line (desktop only) */}
                                        <div className="hidden sm:block absolute top-1/2 left-1/4 right-1/4 h-0.5 bg-gradient-to-r from-red-200 via-purple-300 to-emerald-200 -z-10 -translate-y-1/2 rounded-full opacity-50"></div>
                                    </div>

                                    {/* TO */}
                                    <div className="flex-1 flex flex-col items-center sm:items-start text-center sm:text-left w-full">
                                        <span className="text-xs text-emerald-500 font-bold tracking-wider mb-1 flex items-center gap-1 uppercase">
                                            Receiver <HandCoins size={14} />
                                        </span>
                                        <span className="text-lg font-semibold text-slate-900">{tx.to_user_name}</span>
                                    </div>

                                </div>
                            </CardContent>
                        </Card>
                    ))}

                    <div className="flex justify-between items-center mt-8 p-4 bg-slate-50 rounded-xl border">
                        <div>
                            <p className="text-sm font-medium text-slate-900">Total Transactions Required: <span className="text-purple-600 font-bold">{settlements.length}</span></p>
                            <p className="text-xs text-slate-500 mt-1">Settled automatically using greedy algorithm</p>
                        </div>
                        <Button className="bg-purple-600 hover:bg-purple-700">
                            Mark All as Paid
                        </Button>
                    </div>
                </div>
            )}
        </div>
    );
}
