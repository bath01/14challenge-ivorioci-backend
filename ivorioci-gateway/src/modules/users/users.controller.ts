import { Request, Response, NextFunction } from 'express';
import { UsersService } from './users.service';
import { UpdateProfileSchema } from './users.types';
import { successResponse } from '../../utils/response.utils';

export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  getMe = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      const profile = await this.usersService.getProfile(req.user!.sub);
      res.status(200).json(successResponse(profile));
    } catch (err) {
      next(err);
    }
  };

  updateMe = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      const dto = UpdateProfileSchema.parse(req.body);
      const profile = await this.usersService.updateProfile(req.user!.sub, dto);
      res.status(200).json(successResponse(profile));
    } catch (err) {
      next(err);
    }
  };

  deleteMe = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      await this.usersService.deleteAccount(req.user!.sub);
      res.status(204).send();
    } catch (err) {
      next(err);
    }
  };
}
