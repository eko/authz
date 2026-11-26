import { isEnabled } from './env';

export const isDemoEnabled = (): boolean => isEnabled(import.meta.env.REACT_APP_DEMO_ENABLED);
