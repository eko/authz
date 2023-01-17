import { useState, useEffect, createContext, ReactNode } from 'react';
import { useNavigate } from 'react-router';
import { AuthResponseUser } from 'service/auth/model';
import Cookies from 'js-cookie';
import jwt_decode from 'jwt-decode';

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

export const OAuthAccessTokenCookieName = 'authz_access_token';
export const OAuthExpiresInCookieName = 'authz_expires_in';

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
        let authToken = Cookies.get(OAuthAccessTokenCookieName) || null;
        const authExpiresIn = Cookies.get(OAuthExpiresInCookieName) || null;

        if (authToken === null || authExpiresIn === null) {
            authToken = localStorage.getItem(AuthenticationToken);
        } else {
            const expireAt = new Date();
            expireAt.setSeconds(expireAt.getSeconds() + Number(authExpiresIn));

            localStorage.setItem(AuthenticationToken, authToken);
            localStorage.setItem(AuthenticationExpiration, expireAt.toISOString());

            const decodedHeader = jwt_decode(authToken);
            if (decodedHeader !== undefined) {
                localStorage.setItem(AuthenticationUser, JSON.stringify({
                    username: (decodedHeader as any)?.sub,
                }));
            }

            Cookies.remove(OAuthAccessTokenCookieName);
            Cookies.remove(OAuthExpiresInCookieName);
        }

        const authUser = localStorage.getItem(AuthenticationUser);
        const authExpiration = localStorage.getItem(AuthenticationExpiration);

        if (
            (authToken === undefined || authToken === null)
            || (authUser === undefined || authUser === null)
        ) {
            navigate('/signin');
            return;
        }

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
    // eslint-disable-next-line
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