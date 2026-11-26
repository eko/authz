// isEnabled parses boolean build-time flags: only "true" and "1" enable them,
// so that "false" or "0" don't turn a feature on.
export const isEnabled = (value: string | undefined): boolean =>
  ['true', '1'].includes((value ?? '').trim().toLowerCase());
