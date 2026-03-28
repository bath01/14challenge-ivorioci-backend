import rateLimit from 'express-rate-limit';
import { errorResponse } from '../utils/response.utils';

export const globalLimiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  max: 200,
  standardHeaders: 'draft-7',
  legacyHeaders: false,
  handler: (_req, res) => {
    res.status(429).json(errorResponse('RATE_LIMIT', 'Trop de requêtes, réessayez dans quelques minutes'));
  },
});

export const authLimiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  max: 10,
  standardHeaders: 'draft-7',
  legacyHeaders: false,
  handler: (_req, res) => {
    res.status(429).json(errorResponse('AUTH_RATE_LIMIT', 'Trop de tentatives d\'authentification, réessayez dans 15 minutes'));
  },
});
