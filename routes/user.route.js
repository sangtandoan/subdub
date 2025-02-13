import { Router } from "express";
import { catchErrors } from "../utils/error.utils.js";
import { getAllUsers, getUser } from "../controllers/user.controller.js";
import { authorize } from "../middlewares/auth.middleware.js";

const userRouter = Router();

userRouter.get("/", catchErrors(getAllUsers));

userRouter.get("/:id", authorize, catchErrors(getUser));

userRouter.post("/", (req, res) => {
	res.send({ title: "CREATE a user" });
});

userRouter.put("/:id", (req, res) => {
	res.send({ title: "UPDATE a specific user" });
});

userRouter.delete("/:id", (req, res) => {
	res.send({ title: "DELETE a specific user" });
});

export default userRouter;
