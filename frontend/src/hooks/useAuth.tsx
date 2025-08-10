import { useContext } from 'react';
import AuthContext from '../context/AuthContext';
import { api } from '../app/axios';

export const useAuth = () => {
    const context = useContext(AuthContext);

    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }

    const logout = async () => {
        try {
            await api.post('/logout');
            context.setUser(null);
        } catch (err) { }
    };

    return { ...context, logout };
};
