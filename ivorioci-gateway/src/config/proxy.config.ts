export interface ProxyRouteConfig {
  path: string;
  target: string;
  pathRewrite?: Record<string, string>;
  requireAuth: boolean;
  changeOrigin?: boolean;
}

export const proxyRoutes: ProxyRouteConfig[] = [
  // ── ivorioci-stream-service ──────────────────────────────
  {
    path: '/api/videos',
    target: process.env.STREAM_SERVICE_URL ?? 'http://localhost:8080',
    pathRewrite: { '^/api': '' },
    requireAuth: false,
    changeOrigin: true,
  },
  {
    path: '/api/categories',
    target: process.env.STREAM_SERVICE_URL ?? 'http://localhost:8080',
    pathRewrite: { '^/api': '' },
    requireAuth: false,
    changeOrigin: true,
  },

  {
    path: '/api/stream',
    target: process.env.STREAM_SERVICE_URL ?? 'http://localhost:8080',
    pathRewrite: { '^/api': '' },
    requireAuth: true,
    changeOrigin: true,
  },

  {
    path: '/api/thumbnails',
    target: process.env.STREAM_SERVICE_URL ?? 'http://localhost:8080',
    pathRewrite: { '^/api': '' },
    requireAuth: false,
    changeOrigin: true,
  },
];
