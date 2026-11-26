import dayjs from 'dayjs';
import localizedFormat from 'dayjs/plugin/localizedFormat';
import relativeTime from 'dayjs/plugin/relativeTime';

// Enables localized format tokens (L, LT, LTS, LLL...) and fromNow() used across the UI.
dayjs.extend(localizedFormat);
dayjs.extend(relativeTime);

export default dayjs;
