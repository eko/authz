import { isEnabled } from './env';

export const getOauthButtonLabel = (): string | undefined => import.meta.env.REACT_APP_OAUTH_BUTTON_LABEL || 'Sign-in with Single Sign-On (SSO)';
export const getOauthLogoUrl = (): string | undefined => import.meta.env.REACT_APP_OAUTH_LOGO_URL;
export const isOauthEnabled = (): boolean => isEnabled(import.meta.env.REACT_APP_OAUTH_ENABLED);
