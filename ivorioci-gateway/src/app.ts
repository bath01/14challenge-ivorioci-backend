import express from 'express';
import helmet from 'helmet';
import cors from 'cors';
import { corsOrigins } from './config/env';
import { globalLimiter } from './middleware/rateLimiter.middleware';
import { requestLogger } from './middleware/requestLogger.middleware';
import { errorHandler } from './middleware/errorHandler.middleware';
import { verifyAccessToken } from './middleware/auth.middleware';
import { successResponse } from './utils/response.utils';

import authRouter from './modules/auth/auth.router';
import usersRouter from './modules/users/users.router';
import proxyRouter from './proxy/proxy.router';

const app = express();

app.use(helmet());
app.use(
  cors({
    origin: corsOrigins,
    methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS'],
    allowedHeaders: ['Content-Type', 'Authorization'],
    credentials: true,
  }),
);

app.use(globalLimiter);

app.use(requestLogger);

app.use('/auth', express.json({ limit: '10kb' }));
app.use('/users', express.json({ limit: '10kb' }));

app.get('/health', (_req, res) => {
  res.status(200).json(successResponse({ status: 'ok' }));
});

app.use('/auth', authRouter);
app.use('/users', verifyAccessToken, usersRouter);

app.use(proxyRouter);

app.use((_req, res) => {
  res.status(404).json({
    success: false,
    error: { code: 'NOT_FOUND', message: 'Route introuvable' },
    timestamp: new Date().toISOString(),
  });
});

app.use(errorHandler);

export default app;
