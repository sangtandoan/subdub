const errorMiddleware = (err, req, res, next) => {
    try {
        // shallow copy (nested objects do not copy) 
        let error = { ...err }
        // deep copy (nested objects do copy)
        // let error = JSON.parse(JSON.stringify(err))

        error.message = err.message
        console.error(error)

        // Mongoose bad ObjectId
        if (err.name === "CastError") {
            const message = "Resource not found"

            error = new Error(message)
            error.statusCode = 404
        }

        // Mongoose duplicate key
        if (err.code === 11000) {
            const message = "Duplicate field value entered"

            error = new Error(message)
            error.statusCode = 400
        }

        // Mongoose validation errors
        if (err.name === "ValidationError") {
            const message = Object.values(err.errors).map(val => val.message)

            error = new Error(message.join(", "))
            error.statusCode = 400
        }

        res.status(error.statusCode || 500).json({ success: false, error: error.message || "Server error" })

    } catch (error) {
        next(error)
    }
}

export default errorMiddleware
