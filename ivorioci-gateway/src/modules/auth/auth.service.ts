import { hashPassword, comparePassword } from '../../utils/password.utils';
import { signAccessToken, signRefreshToken, verifyRefreshToken } from '../../utils/jwt.utils';
import { RegisterDto, LoginDto, AuthTokens } from './auth.types';
import { prisma } from '../../utils/prisma';

export class AuthService {

  async register(dto: RegisterDto): Promise<AuthTokens> {
    const existing = await prisma.user.findUnique({
      where: { email: dto.email }
    });
    if (existing) {
      throw Object.assign(new Error('Un compte avec cet email existe déjà'), { statusCode: 409 });
    }

    const passwordHash = await hashPassword(dto.password);
    const user = await prisma.user.create({
      data: {
        email: dto.email,
        passwordHash,
        firstName: dto.firstName,
        lastName: dto.lastName,
    }});

    return this.issueTokens(user.id, user.email);
  }

  async login(dto: LoginDto): Promise<AuthTokens> {
    const user = await prisma.user.findUnique({
      where: { email: dto.email }
    });
    if (!user) {
      throw Object.assign(new Error('Email ou mot de passe incorrect'), { statusCode: 401 });
    }

    const valid = await comparePassword(dto.password, user.passwordHash);
    if (!valid) {
      throw Object.assign(new Error('Email ou mot de passe incorrect'), { statusCode: 401 });
    }

    return this.issueTokens(user.id, user.email);
  }

  async refresh(refreshToken: string): Promise<AuthTokens> {
    let payload;
    try {
      payload = verifyRefreshToken(refreshToken);
    } catch {
      throw Object.assign(new Error('Refresh token invalide ou expiré'), { statusCode: 401 });
    }

    const user = await prisma.user.findUnique({
      where: { id: payload.sub }
    });
    if (!user) {
      throw Object.assign(new Error('Refresh token révoqué'), { statusCode: 401 });
    }

    // Rotation : on révoque l'ancien token avant d'en émettre un nouveau
    await prisma.refreshToken.delete({
      where: { token: refreshToken }
    });

    return this.issueTokens(user.id, user.email);
  }

  async logout(refreshToken: string): Promise<void> {
    await prisma.refreshToken.delete({
      where: { token: refreshToken }
    });
  }

  async logoutAll(userId: string): Promise<void> {
    await prisma.refreshToken.deleteMany({
      where: { userId }
    });
  }

  private async issueTokens(userId: string, email: string): Promise<AuthTokens> {
    const accessToken = signAccessToken({ sub: userId, email });
    const refreshToken = signRefreshToken({ sub: userId, email });
    await prisma.refreshToken.create({
      data: {
        token: refreshToken,
        userId,
        expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000) // 7 days
      }
    });
    return { accessToken, refreshToken };
  }
}
