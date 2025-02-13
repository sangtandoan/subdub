import { config } from "dotenv";

console.log(process.cwd());
// path = process.cwd() + ".env.development.local"
config({ path: `.env.${process.env.NODE_ENV || "development"}.local` });

console.log(process.env.DB_URI);
export const {
	PORT,
	DB_URI,
	JWT_SECRET,
	JWT_EXPIRE_TIME,
	NODE_ENV,
	ARCJET_KEY,
	ARCJET_ENV,
} = process.env;
