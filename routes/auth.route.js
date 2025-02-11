import { Router } from "express"

const authRouter = Router()


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
    return function (req, res, next) {
        fn(req, res, next).catch(next)
    }
}

authRouter.get("/", catchErrors(async (req, res) => {
    throw new Error("Test catch errors")
}))

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
