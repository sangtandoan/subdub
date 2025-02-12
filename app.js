import express from "express";
import { PORT } from "./config/env.js";
import connectToDatabase from "./database/mongodb.js";
import errorMiddleware from "./middlewares/error.middleware.js";
import cookieParser from "cookie-parser";
import compression from "compression";
import router from "./routes/index.js";

const app = express();
// Using middlewares
app.use(express.json());
// parse form body to req.body
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(compression());

app.use("/api/v1", router);

// Using global error handler
app.use(errorMiddleware);

app.listen(PORT, async () => {
    console.log(
        `Subscription Tracker API is running on http://localhost:${PORT}`,
    );

    await connectToDatabase();
});

export default app;
