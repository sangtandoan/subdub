import { Router } from "express";
import { signIn, signUp } from "../controllers/auth.controller.js";
import { catchErrors } from "../utils/error.utils.js";

const authRouter = Router();

// authRouter.post("/sign-in", async (req, res, next) => {
//     try {
//         const data = await signInHandler(req.username, req.password)
//         res.send({ title: "Sign-in" })
//     } catch (error) {
//         next(error)
//     }
// })

authRouter.post("/sign-in", catchErrors(signIn));

authRouter.post("/sign-out", (req, res) => {
	res.send({ title: "Sign-out" });
});

authRouter.post("/sign-up", catchErrors(signUp));

export default authRouter;
