/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly REACT_APP_API_BASE_URI?: string;
  readonly REACT_APP_DEMO_ENABLED?: string;
  readonly REACT_APP_OAUTH_BUTTON_LABEL?: string;
  readonly REACT_APP_OAUTH_ENABLED?: string;
  readonly REACT_APP_OAUTH_LOGO_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
