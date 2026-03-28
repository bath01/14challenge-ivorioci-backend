import { z } from 'zod';
import dotenv from 'dotenv';

dotenv.config();

const envSchema = z.object({
  PORT: z.coerce.number().default(3000),
  NODE_ENV: z.enum(['development', 'production', 'test']).default('development'),
  LOG_LEVEL: z.enum(['error', 'warn', 'info', 'debug']).default('info'),

  JWT_ACCESS_SECRET: z.string().min(16, 'JWT_ACCESS_SECRET doit faire au moins 16 caractères'),
  JWT_REFRESH_SECRET: z.string().min(16, 'JWT_REFRESH_SECRET doit faire au moins 16 caractères'),
  JWT_ACCESS_EXPIRES_IN: z.string().default('15m'),
  JWT_REFRESH_EXPIRES_IN: z.string().default('7d'),

  CORS_ALLOWED_ORIGINS: z.string().default('http://localhost:3000'),
});

const parsed = envSchema.safeParse(process.env);

if (!parsed.success) {
  console.error('Variables d\'environnement invalides:');
  console.error(parsed.error.flatten().fieldErrors);
  process.exit(1);
}

export const env = Object.freeze(parsed.data);

export const corsOrigins = env.CORS_ALLOWED_ORIGINS.split(',').map((o) => o.trim());
