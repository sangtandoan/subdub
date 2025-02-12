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
        const existingUser = await User.findOne({ email }).lean();
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
        ).lean();

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

export const signIn = async (req, res) => {
    const { email, password } = req.body;

    // Check user exists
    const existingUser = await User.findOne({ email }).lean();
    if (!existingUser) {
        const error = new Error("User not found!");
        error.statusCode = 404;

        throw error;
    }

    if (!bcrypt.compareSync(password, existingUser.password)) {
        const error = new Error("Invalid password");
        error.statusCode = 401;

        throw error;
    }

    const token = jwt.sign({ id: existingUser._id }, JWT_SECRET, {
        expiresIn: JWT_EXPIRE_TIME,
    });

    res.status(200).json({
        success: true,
        message: "Login successfully!",
        data: {
            token,
        },
    });
};
