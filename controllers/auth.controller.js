import mongoose from "mongoose";
import User from "../models/user.model.js";
import bcrypt from "bcryptjs";
import jwt from "jsonwebtoken";
import { JWT_SECRET, JWT_EXPIRE_TIME } from "../config/env.js";

export const signUp = async (req, res, next) => {
    // Starts mongoose session so that can mongoose can do transaction
    const session = await mongoose.startSession();
    session.startTransaction();

    try {
        const { name, email, password } = req.body;

        // Check if user exists
        const existingUser = await User.findOne({ email });
        if (existingUser) {
            const error = new Error("User already exists!");
            error.statusCode = 409;
            throw error;
        }

        // Hash password
        const salt = await bcrypt.genSalt();
        const hashedPassword = await bcrypt.hash(password, salt);

        const newUsers = await User.create(
            [{ name, email, password: hashedPassword }],
            { session },
        );

        // Create access token
        const token = jwt.sign({ id: newUsers[0]._id }, JWT_SECRET, {
            expiresIn: JWT_EXPIRE_TIME,
        });

        await session.commitTransaction();
        session.endSession();

        res.status(201).json({
            success: true,
            message: "Create User successfully",
            data: {
                token,
                user: newUsers[0],
            },
        });
    } catch (error) {
        await session.abortTransaction();
        session.endSession();

        next(error);
    }
};
