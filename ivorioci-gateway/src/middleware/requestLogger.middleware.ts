import morgan from 'morgan';
import { logger } from '../utils/logger';
import { env } from '../config/env';

const stream = {
  write: (message: string) => logger.http(message.trimEnd()),
};

export const requestLogger = morgan(
  env.NODE_ENV === 'production' ? 'combined' : 'dev',
  { stream },
);
