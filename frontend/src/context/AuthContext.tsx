import {
    createContext,
    useEffect,
    useState,
} from 'react';
import type { ReactNode } from 'react';
import { api } from '../app/axios';

export interface User {
    id: string;
    email: string;
    name: string;
}

interface AuthContextType {
    user: User | null;
    loading: boolean;
    error: string | null;
    setUser: (user: User | null) => void;
}

const AuthContext = createContext<AuthContextType>({
    user: null,
    loading: true,
    error: null,
    setUser: () => { },
});

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const getProfile = async () => {
            try {
                const res = await api.get('/profile');
                setUser(res.data);
            } catch (err: any) {
                setUser(null);
                if (err.response?.status !== 401) {
                    setError('Error receiving profile');
                } else {
                    setError('Unauthorized');
                }
            } finally {
                setLoading(false);
            }
        };

        getProfile();
    }, []);

    return (
        <AuthContext.Provider value={{ user, loading, error, setUser }}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthContext;
