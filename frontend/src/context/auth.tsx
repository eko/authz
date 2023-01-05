import { useState, useEffect, createContext, ReactNode } from 'react';
import { useNavigate } from 'react-router';
import { AuthResponseUser } from 'service/auth/model';

export type User = {
    username: string
    token: string
}

const AuthContext = createContext<AuthContextType>({
    user: undefined,
    logout: () => {},
});

type AuthContextType = {
    user: User | undefined
    logout: Function
};

type AuthContextProviderType = {
    children: ReactNode
}

export const AuthenticationExpiration = 'authzExpiration';
export const AuthenticationToken = 'authzToken';
export const AuthenticationUser = 'authzUser';

const AuthContextProvider = ({ children }: AuthContextProviderType) => {
    const [user, setUser] = useState<User | undefined>();
    const navigate = useNavigate();
    
    const logout = () => {
        localStorage.removeItem(AuthenticationToken);
        localStorage.removeItem(AuthenticationExpiration);
        localStorage.removeItem(AuthenticationUser);
        setUser(undefined);
        navigate('/signin');
    }

    useEffect(() => {
        const authToken = localStorage.getItem(AuthenticationToken);
        const authUser = localStorage.getItem(AuthenticationUser);
 
        if (
            (authToken === undefined || authToken === null)
            || (authUser === undefined || authUser === null)
        ) {
            navigate('/signin');
            return;
        }

        const authExpiration = localStorage.getItem(AuthenticationExpiration);
        const expirationDate = new Date(authExpiration!);

        if (expirationDate < new Date()) {
            logout();
            return;
        }

        const userObject = JSON.parse(authUser) as AuthResponseUser;

        const userData: User = {
            username: userObject.username,
            token: authToken,
        }

        setUser(userData);
    }, [navigate]);

    return (
        <AuthContext.Provider value={{
            user,
            logout,
        }}>
            {user?.token ? (
                <>{children}</>
            ) : null}
        </AuthContext.Provider>
    );
}

export {
    AuthContext,
    AuthContextProvider,
};