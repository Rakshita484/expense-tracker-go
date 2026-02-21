import axios from 'axios';

// Create an Axios instance configured to talk to the Go backend
const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
    headers: {
        'Content-Type': 'application/json',
    },
});

export const userService = {
    getUsers: () => api.get('/users').then(res => res.data.data),
    createUser: (data) => api.post('/users', data).then(res => res.data.data),
};

export const groupService = {
    getGroups: () => api.get('/groups').then(res => res.data.data),
    createGroup: (data) => api.post('/groups', data).then(res => res.data.data),
    addMember: (groupId, userId) => api.post(`/groups/${groupId}/members`, { user_id: userId }).then(res => res.data.data),
    getMembers: (groupId) => api.get(`/groups/${groupId}/members`).then(res => res.data.data),
    addExpense: (groupId, data) => api.post(`/groups/${groupId}/expenses`, data).then(res => res.data.data),
    getExpenses: (groupId) => api.get(`/groups/${groupId}/expenses`).then(res => res.data.data),
    getBalances: (groupId) => api.get(`/groups/${groupId}/balances`).then(res => res.data.data),
    getSettlements: (groupId) => api.get(`/groups/${groupId}/settlements`).then(res => res.data.data),
};

export const dashboardService = {
    getStats: () => api.get('/dashboard/stats').then(res => res.data.data),
};

export default api;
