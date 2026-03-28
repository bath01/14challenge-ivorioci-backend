import { Request, Response, NextFunction } from 'express';
import { verifyAccessToken as verifyToken } from '../utils/jwt.utils';
import { errorResponse } from '../utils/response.utils';

export function verifyAccessToken(req: Request, res: Response, next: NextFunction): void {
  const authHeader = req.headers.authorization;

  if (!authHeader?.startsWith('Bearer ')) {
    res.status(401).json(errorResponse('UNAUTHORIZED', 'Token d\'accès manquant'));
    return;
  }

  const token = authHeader.slice(7);

  try {
    req.user = verifyToken(token);
    next();
  } catch (err: any) {
    if (err?.name === 'TokenExpiredError') {
      res.status(401).json(errorResponse('TOKEN_EXPIRED', 'Token expiré'));
    } else {
      res.status(401).json(errorResponse('TOKEN_INVALID', 'Token invalide'));
    }
  }
}
