import http from 'http';
import app from './app';
import { env } from './config/env';
import { logger } from './utils/logger';

const server = http.createServer(app);

server.listen(env.PORT, () => {
  logger.info(`API Gateway démarrée sur le port ${env.PORT} [${env.NODE_ENV}]`);
});

function gracefulShutdown(signal: string): void {
  logger.info(`Signal ${signal} reçu — arrêt gracieux en cours...`);
  server.close(() => {
    logger.info('Serveur HTTP fermé');
    process.exit(0);
  });

  setTimeout(() => {
    logger.warn('Arrêt forcé après timeout');
    process.exit(1);
  }, 10_000);
}

process.on('SIGTERM', () => gracefulShutdown('SIGTERM'));
process.on('SIGINT', () => gracefulShutdown('SIGINT'));

process.on('unhandledRejection', (reason) => {
  logger.error('Unhandled promise rejection', { reason });
});

process.on('uncaughtException', (err) => {
  logger.error('Uncaught exception', { message: err.message, stack: err.stack });
  process.exit(1);
});
