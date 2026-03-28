import { Router } from 'express';
import { UsersController } from './users.controller';
import { UsersService } from './users.service';

const usersService = new UsersService();
const usersController = new UsersController(usersService);

const router = Router();

router.get('/me', usersController.getMe);
router.put('/me', usersController.updateMe);
router.delete('/me', usersController.deleteMe);

export default router;
