import { User } from '@prisma/client';
import { z } from 'zod';

export const UpdateProfileSchema = z.object({
  firstName: z.string().min(1).max(50).optional(),
  lastName: z.string().min(1).max(50).optional(),
  email: z.string().email('Email invalide').optional(),
});

export type UpdateProfileDto = z.infer<typeof UpdateProfileSchema>;

export type PublicUser = Omit<User, 'passwordHash'>;

export function toPublicUser(user: User): PublicUser {
  const { passwordHash: _, ...publicUser } = user;
  return publicUser;
}
