import React, { useState, useEffect } from 'react';
import { groupService, userService } from '@/services/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select } from '@/components/ui/select';
import { Receipt, CheckCircle2, AlertCircle, Clock, User, Users2, ChevronDown, ChevronUp } from 'lucide-react';
import toast from 'react-hot-toast';

export default function Expenses() {
    const [users, setUsers] = useState([]);
    const [balances, setBalances] = useState([]);
    const [expenses, setExpenses] = useState([]);
    const [loading, setLoading] = useState(false);
    const [expandedExpense, setExpandedExpense] = useState(null);

    // Static group ID 1 for MVP (Matches backend seed script)
    const groupId = 1;

    // Form state
    const [description, setDescription] = useState('');
    const [amount, setAmount] = useState('');
    const [paidBy, setPaidBy] = useState('');
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        fetchData();
    }, [groupId]);

    const fetchData = async () => {
        setLoading(true);
        try {
            const [usersData, balancesData, expensesData] = await Promise.all([
                userService.getUsers(),
                groupService.getBalances(groupId),
                groupService.getExpenses(groupId)
            ]);
            setUsers(usersData || []);
            setBalances(balancesData || []);
            setExpenses(expensesData || []);
        } catch (error) {
            toast.error('Failed to load data from server');
        } finally {
            setLoading(false);
        }
    };

    const handleAddExpense = async (e) => {
        e.preventDefault();
        if (!description || !amount || !paidBy) {
            return toast.error("Please fill all fields");
        }

        setSubmitting(true);
        try {
            await groupService.addExpense(groupId, {
                description,
                amount,
                paid_by_id: parseInt(paidBy)
            });

            toast.success('Expense added successfully!');
            setDescription('');
            setAmount('');
            setPaidBy('');

            // Refresh all data after adding expense
            fetchData();
        } catch (error) {
            toast.error(error.response?.data?.details || 'Failed to add expense');
        } finally {
            setSubmitting(false);
        }
    };

    const formatDate = (dateStr) => {
        const date = new Date(dateStr);
        return date.toLocaleDateString('en-US', {
            month: 'short', day: 'numeric', year: 'numeric',
            hour: '2-digit', minute: '2-digit'
        });
    };

    return (
        <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <div className="flex items-center justify-between border-b pb-4">
                <div>
                    <h2 className="text-2xl font-bold tracking-tight text-slate-900">Expenses & Balances</h2>
                    <p className="text-muted-foreground mt-1">Add new expenses, see who paid what, and track balances.</p>
                </div>
                <div className="bg-amber-100/50 text-amber-700 px-3 py-1.5 rounded-full text-sm font-medium border border-amber-200">
                    Group: Weekend Trip
                </div>
            </div>

            <div className="grid lg:grid-cols-2 gap-8">

                {/* ADD EXPENSE FORM */}
                <Card className="border-t-4 border-t-amber-500 shadow-md h-fit">
                    <CardHeader>
                        <div className="flex items-center gap-2">
                            <div className="p-2 bg-amber-100 rounded-lg text-amber-600">
                                <Receipt size={20} />
                            </div>
                            <div>
                                <CardTitle>Add New Expense</CardTitle>
                                <CardDescription>Will be split equally among all members.</CardDescription>
                            </div>
                        </div>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleAddExpense} className="space-y-5">
                            <div className="space-y-2">
                                <Label htmlFor="desc">What was this for?</Label>
                                <Input
                                    id="desc"
                                    placeholder="e.g. Dinner at Italian Restaurant"
                                    value={description}
                                    onChange={e => setDescription(e.target.value)}
                                />
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div className="space-y-2">
                                    <Label htmlFor="amount">Amount ($)</Label>
                                    <Input
                                        id="amount"
                                        type="number"
                                        step="0.01"
                                        min="0.01"
                                        placeholder="0.00"
                                        value={amount}
                                        onChange={e => setAmount(e.target.value)}
                                        className="font-mono text-lg"
                                    />
                                </div>
                                <div className="space-y-2">
                                    <Label htmlFor="paidBy">Paid By</Label>
                                    <Select value={paidBy} onChange={e => setPaidBy(e.target.value)}>
                                        <option value="" disabled>Select payer...</option>
                                        {users.map(u => (
                                            <option key={u.id} value={u.id}>{u.name}</option>
                                        ))}
                                    </Select>
                                </div>
                            </div>

                            <Button type="submit" disabled={submitting} className="w-full bg-amber-500 hover:bg-amber-600">
                                {submitting ? 'Adding...' : 'Save Expense & Split Equally'}
                            </Button>
                        </form>
                    </CardContent>
                </Card>

                {/* BALANCES LIST */}
                <div className="space-y-4">
                    <div>
                        <h3 className="text-lg font-semibold text-slate-900 border-b pb-2">Current Net Balances</h3>
                        <p className="text-sm text-muted-foreground mt-1">
                            Positive means they get money back. Negative means they owe money.
                        </p>
                    </div>

                    <div className="grid gap-3">
                        {balances.length === 0 && !loading && (
                            <div className="text-center py-8 text-slate-500 bg-slate-50 border border-dashed rounded-lg">
                                No balances calculated yet.
                            </div>
                        )}

                        {balances.map(b => {
                            const amount = parseFloat(b.balance);
                            const isCreditor = amount > 0;
                            const isDebtor = amount < 0;
                            const isSettled = amount === 0;

                            return (
                                <div
                                    key={b.user_id}
                                    className={`flex items-center justify-between p-4 rounded-xl border bg-white shadow-sm transition-transform hover:-translate-y-0.5
                    ${isCreditor ? 'border-l-4 border-l-emerald-500' : ''}
                    ${isDebtor ? 'border-l-4 border-l-red-500' : ''}
                    ${isSettled ? 'border-l-4 border-l-slate-300 opacity-60' : ''}
                  `}
                                >
                                    <div className="flex items-center gap-3">
                                        <div className="h-10 w-10 rounded-full bg-slate-100 flex items-center justify-center font-semibold text-slate-600">
                                            {b.name.charAt(0)}
                                        </div>
                                        <div>
                                            <p className="font-semibold text-slate-900">{b.name}</p>
                                            <p className="text-xs text-slate-500 flex items-center gap-1">
                                                {isSettled ? <CheckCircle2 size={12} className="text-slate-400" /> : <AlertCircle size={12} />}
                                                {isCreditor ? 'Owed back' : isDebtor ? 'Owes others' : 'All Settled up'}
                                            </p>
                                        </div>
                                    </div>

                                    <div className={`text-right ${isCreditor ? 'text-emerald-600' : isDebtor ? 'text-red-600' : 'text-slate-500'}`}>
                                        <p className={`font-mono text-lg font-bold ${isSettled ? 'text-slate-400' : ''}`}>
                                            {isCreditor ? '+' : ''}${Math.abs(amount).toFixed(2)}
                                        </p>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                </div>

            </div>

            {/* EXPENSE HISTORY */}
            <div className="space-y-4">
                <div className="border-b pb-3">
                    <h3 className="text-lg font-semibold text-slate-900 flex items-center gap-2">
                        <Clock size={20} className="text-slate-500" />
                        Expense History
                    </h3>
                    <p className="text-sm text-muted-foreground mt-1">
                        Full log of all expenses — who paid, how much, and the per-person split.
                    </p>
                </div>

                {expenses.length === 0 && !loading ? (
                    <div className="text-center py-12 text-slate-400 bg-slate-50 border border-dashed rounded-xl">
                        <Receipt className="h-10 w-10 mx-auto mb-3 text-slate-300" />
                        <p className="text-sm font-medium">No expenses recorded yet</p>
                        <p className="text-xs mt-1">Add an expense above to see the history here</p>
                    </div>
                ) : (
                    <div className="space-y-3">
                        {expenses.map(expense => {
                            const isExpanded = expandedExpense === expense.id;
                            return (
                                <Card
                                    key={expense.id}
                                    className="overflow-hidden hover:shadow-md transition-shadow cursor-pointer"
                                    onClick={() => setExpandedExpense(isExpanded ? null : expense.id)}
                                >
                                    <div className="flex items-center gap-4 p-4">
                                        {/* Payer Avatar */}
                                        <div className="h-12 w-12 rounded-full bg-gradient-to-br from-amber-400 to-amber-600 flex items-center justify-center text-white font-bold text-lg shrink-0 shadow-sm">
                                            {expense.paid_by?.name?.charAt(0) || '?'}
                                        </div>

                                        {/* Expense Info */}
                                        <div className="flex-1 min-w-0">
                                            <p className="text-sm font-bold text-slate-900">{expense.description}</p>
                                            <p className="text-xs text-slate-500 mt-0.5">
                                                Paid by <span className="font-semibold text-amber-700">{expense.paid_by?.name || 'Unknown'}</span>
                                                <span className="mx-1.5">•</span>
                                                {formatDate(expense.created_at)}
                                            </p>
                                        </div>

                                        {/* Amount & Expand */}
                                        <div className="text-right shrink-0 flex items-center gap-2">
                                            <div>
                                                <p className="font-mono text-lg font-bold text-slate-900">${parseFloat(expense.amount).toFixed(2)}</p>
                                                <p className="text-xs text-slate-400">{expense.splits?.length || 0} way split</p>
                                            </div>
                                            {isExpanded ?
                                                <ChevronUp size={16} className="text-slate-400" /> :
                                                <ChevronDown size={16} className="text-slate-400" />
                                            }
                                        </div>
                                    </div>

                                    {/* Expanded Split Details */}
                                    {isExpanded && expense.splits && (
                                        <div className="border-t bg-slate-50/80 px-4 py-3">
                                            <p className="text-xs font-medium text-slate-500 uppercase tracking-wider mb-2">Split Breakdown</p>
                                            <div className="grid gap-2 sm:grid-cols-2">
                                                {expense.splits.map(split => (
                                                    <div
                                                        key={split.id}
                                                        className={`flex items-center justify-between p-2.5 rounded-lg bg-white border text-sm ${split.user_id === expense.paid_by_id
                                                                ? 'border-amber-200 bg-amber-50/50'
                                                                : 'border-slate-200'
                                                            }`}
                                                    >
                                                        <div className="flex items-center gap-2">
                                                            <div className="h-7 w-7 rounded-full bg-slate-100 flex items-center justify-center text-xs font-semibold text-slate-600">
                                                                {split.user?.name?.charAt(0) || '?'}
                                                            </div>
                                                            <span className="font-medium text-slate-700">{split.user?.name || 'Unknown'}</span>
                                                            {split.user_id === expense.paid_by_id && (
                                                                <span className="text-[10px] bg-amber-100 text-amber-700 px-1.5 py-0.5 rounded-full font-semibold">PAYER</span>
                                                            )}
                                                        </div>
                                                        <span className="font-mono font-semibold text-slate-900">${parseFloat(split.amount).toFixed(2)}</span>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    )}
                                </Card>
                            );
                        })}
                    </div>
                )}
            </div>
        </div>
    );
}
