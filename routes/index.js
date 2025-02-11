import { Router } from "express"
import authRouter from "./auth.route.js"
import userRouter from "./user.route.js"
import subscriptionRouter from "./subscription.route.js"

const router = Router()


router.use("/auth", authRouter)
router.use("/users", userRouter)
router.use("/subscriptions", subscriptionRouter)

export default router

