import { Router, Request, Response, NextFunction } from 'express';
import { createProxyMiddleware } from 'http-proxy-middleware';
import { proxyRoutes } from '../config/proxy.config';
import { verifyAccessToken } from '../middleware/auth.middleware';
import { logger } from '../utils/logger';

const router = Router();

for (const route of proxyRoutes) {
  if (!route.target) {
    logger.warn(`Proxy route ${route.path} ignorée — target manquant`);
    continue;
  }

  const proxyMw = createProxyMiddleware({
    target: route.target,
    changeOrigin: route.changeOrigin ?? true,
    pathFilter: route.path,
    pathRewrite: route.pathRewrite,
    on: {
      error: (err, _req, res) => {
        logger.error(`Proxy error [${route.path} -> ${route.target}]`, { message: (err as Error).message });
        const httpRes = res as import('http').ServerResponse;
        if (!httpRes.headersSent) {
          httpRes.writeHead(502, { 'Content-Type': 'application/json' });
          httpRes.end(JSON.stringify({
            success: false,
            error: { code: 'BAD_GATEWAY', message: 'Le microservice est indisponible' },
            timestamp: new Date().toISOString(),
          }));
        }
      },
      proxyReq: (_proxyReq, req) => {
        logger.debug(`Proxy -> ${route.target}${(req as any).url}`);
      },
    },
  });

  if (route.requireAuth) {
    const conditionalAuth = (req: Request, res: Response, next: NextFunction) => {
      if (req.path === route.path || req.path.startsWith(route.path + '/')) {
        return verifyAccessToken(req, res, next);
      }
      next();
    };
    router.use(conditionalAuth, proxyMw as any);
  } else {
    router.use(proxyMw as any);
  }

  logger.info(`Proxy enregistré: ${route.path} -> ${route.target} (auth: ${route.requireAuth})`);
}

export default router;
