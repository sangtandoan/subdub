import { Router } from "express";
import { signUp } from "../controllers/auth.controller.js";

const authRouter = Router();

// authRouter.post("/sign-in", async (req, res, next) => {
//     try {
//         const data = await signInHandler(req.username, req.password)
//         res.send({ title: "Sign-in" })
//     } catch (error) {
//         next(error)
//     }
// })

// HOF to handle throw, unhandled, syntax errors that does not call next()
const catchErrors = (fn) => {
	return (req, res, next) => {
		return fn(req, res, next).catch(next);
	};
};

authRouter.post("/sign-in", (req, res) => {
	res.send({ title: "Sign-in" });
});

authRouter.post("/sign-out", (req, res) => {
	res.send({ title: "Sign-out" });
});

authRouter.post("/sign-up", catchErrors(signUp));

export default authRouter;
