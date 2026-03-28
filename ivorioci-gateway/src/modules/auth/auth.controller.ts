import { Request, Response, NextFunction } from 'express';
import { AuthService } from './auth.service';
import { RegisterSchema, LoginSchema, RefreshSchema } from './auth.types';
import { successResponse } from '../../utils/response.utils';

export class AuthController {
  constructor(private readonly authService: AuthService) {}

  register = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      const dto = RegisterSchema.parse(req.body);
      const tokens = await this.authService.register(dto);
      res.status(201).json(successResponse(tokens));
    } catch (err) {
      next(err);
    }
  };

  login = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      const dto = LoginSchema.parse(req.body);
      const tokens = await this.authService.login(dto);
      res.status(200).json(successResponse(tokens));
    } catch (err) {
      next(err);
    }
  };

  refresh = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { refreshToken } = RefreshSchema.parse(req.body);
      const tokens = await this.authService.refresh(refreshToken);
      res.status(200).json(successResponse(tokens));
    } catch (err) {
      next(err);
    }
  };

  logout = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { refreshToken } = RefreshSchema.parse(req.body);
      await this.authService.logout(refreshToken);
      res.status(204).send();
    } catch (err) {
      next(err);
    }
  };

  logoutAll = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      await this.authService.logoutAll(req.user!.sub);
      res.status(204).send();
    } catch (err) {
      next(err);
    }
  };
}
