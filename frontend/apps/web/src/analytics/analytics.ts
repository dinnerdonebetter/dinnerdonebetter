import { browserAnalyticsWrapper, serverAnalyticsWrapper } from '@dinnerdonebetter/analytics';

export const browserSideAnalytics = new browserAnalyticsWrapper(process.env.NEXT_PUBLIC_SEGMENT_API_TOKEN || '');
export const serverSideAnalytics = new serverAnalyticsWrapper(process.env.NEXT_SEGMENT_API_TOKEN || '');
