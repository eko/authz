export const getOauthButtonLabel = (): string | undefined => process.env.REACT_APP_OAUTH_BUTTON_LABEL || 'Sign-in wigh Single Sign-On (SSO)';
export const getOauthLogoUrl = (): string | undefined => process.env.REACT_APP_OAUTH_LOGO_URL;
export const isOauthEnabled = (): boolean => process.env.REACT_APP_OAUTH_ENABLED ? true : false;
