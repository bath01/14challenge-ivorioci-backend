import { Router } from 'express';
import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import { verifyAccessToken } from '../../middleware/auth.middleware';
import { authLimiter } from '../../middleware/rateLimiter.middleware';

const authService = new AuthService();
const authController = new AuthController(authService);

const router = Router();

router.post('/register', authLimiter, authController.register);
router.post('/login', authLimiter, authController.login);
router.post('/refresh', authLimiter, authController.refresh);

router.post('/logout', verifyAccessToken, authController.logout);
router.post('/logout-all', verifyAccessToken, authController.logoutAll);

export default router;
