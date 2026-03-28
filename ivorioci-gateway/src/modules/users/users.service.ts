import { prisma } from '../../utils/prisma';
import { PublicUser, UpdateProfileDto, toPublicUser } from './users.types';

export class UsersService {

  async getProfile(userId: string): Promise<PublicUser> {
    const user = await prisma.user.findUnique({
      where: { id: userId }
    });
    if (!user) {
      throw Object.assign(new Error('Utilisateur introuvable'), { statusCode: 404 });
    }
    return toPublicUser(user);
  }

  async updateProfile(userId: string, dto: UpdateProfileDto): Promise<PublicUser> {
    if (dto.email) {
      const existing = await prisma.user.findUnique({
        where: { email: dto.email }
      });
      if (existing && existing.id !== userId) {
        throw Object.assign(new Error('Cet email est déjà utilisé'), { statusCode: 409 });
      }
    }

    const updated = await prisma.user.update({
      where: { id: userId },
      data: dto
    });
    return toPublicUser(updated);
  }

  async deleteAccount(userId: string): Promise<void> {
    const deleted = await prisma.user.delete({
      where: { id: userId }
    });
    if (!deleted) {
      throw Object.assign(new Error('Utilisateur introuvable'), { statusCode: 404 });
    }
  }
}
