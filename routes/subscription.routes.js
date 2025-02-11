import { Router } from "express"

const subscriptionRouter = Router()

subscriptionRouter.get("/", (req, res) => {
    res.send({ title: "GET all subscriptions" })
})

subscriptionRouter.get("/:id", (req, res) => {
    res.send({ title: "GET a specific subscription" })
})

subscriptionRouter.post("/", (req, res) => {
    res.send({ title: "CREATE a subscription" })
})

subscriptionRouter.put("/:id", (req, res) => {
    res.send({ title: "UPDATE a specific subscription" })
})

subscriptionRouter.delete("/:id", (req, res) => {
    res.send({ title: "DELETE a specific subscription" })
})

subscriptionRouter.get("/user/:id", (req, res) => {
    res.send({ title: "GET all subscriptions of user" })
})

subscriptionRouter.put("/:id/cancel", (req, res) => {
    res.send({ title: "CANCEL subscription" })
})

subscriptionRouter.get("/upcoming-renewals", (req, res) => {
    res.send({ title: "GET upcoming renewals" })
})

export default subscriptionRouter

