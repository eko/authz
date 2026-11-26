export const getOauthButtonLabel = (): string | undefined => import.meta.env.REACT_APP_OAUTH_BUTTON_LABEL || 'Sign-in wigh Single Sign-On (SSO)';
export const getOauthLogoUrl = (): string | undefined => import.meta.env.REACT_APP_OAUTH_LOGO_URL;
export const isOauthEnabled = (): boolean => import.meta.env.REACT_APP_OAUTH_ENABLED ? true : false;
