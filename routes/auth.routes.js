import { Router } from "express"

const authRouter = Router()

authRouter.post("/sign-in", (req, res) => {
    res.send({ title: "Sign-in" })
})

authRouter.post("/sign-out", (req, res) => {
    res.send({ title: "Sign-out" })
})

authRouter.post("/sign-up", (req, res) => {
    res.send("Sign-up")
})

export default authRouter
