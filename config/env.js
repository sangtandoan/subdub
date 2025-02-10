import { config } from "dotenv"

console.log(process.env.NODE_ENV)

config({ path: `.env.${process.env.NODE_ENV || 'development'}.local` })

export const { PORT, NODE_ENV } = process.env

