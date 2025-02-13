import jwt from "jsonwebtoken";
import { JWT_SECRET } from "../config/env.js";
import User from "../models/user.model.js";

export const authorize = async (req, res, next) => {
	try {
		let token;

		if (req.headers.authorization?.startsWith("Bearer")) {
			console.log(req.headers.authorization);
			token = req.headers.authorization.split(" ")[1];
		}

		if (!token) {
			throw new Error("Invalid authorization header!");
		}

		const payload = jwt.verify(token, JWT_SECRET);
		if (!payload) {
			throw new Error("Invalid token");
		}

		const user = await User.findOne({ _id: payload.id }).lean();
		if (!user) {
			throw new Error("Invalid user");
		}

		req.user = user;

		next();
	} catch (error) {
		res.status(401).json({ message: "Unauthorized", error: error.message });
	}
};
