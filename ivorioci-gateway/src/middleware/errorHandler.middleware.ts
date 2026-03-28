import { Request, Response, NextFunction } from 'express';
import { ZodError } from 'zod';
import { errorResponse } from '../utils/response.utils';
import { logger } from '../utils/logger';

interface AppError extends Error {
  statusCode?: number;
}

export function errorHandler(
  err: AppError,
  req: Request,
  res: Response,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  _next: NextFunction,
): void {
  // Erreurs de validation Zod
  if (err instanceof ZodError) {
    const message = err.errors.map((e) => `${e.path.join('.')}: ${e.message}`).join('; ');
    res.status(400).json(errorResponse('VALIDATION_ERROR', message));
    return;
  }

  const statusCode = err.statusCode ?? 500;

  if (statusCode >= 500) {
    logger.error('Erreur interne', { message: err.message, stack: err.stack, path: req.path });
  }

  res.status(statusCode).json(
    errorResponse(
      statusCode === 500 ? 'INTERNAL_ERROR' : 'ERROR',
      statusCode === 500 ? 'Une erreur interne est survenue' : err.message,
    ),
  );
}
