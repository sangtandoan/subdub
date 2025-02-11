import mongoose from "mongoose";

const subscriptionSchema = new mongoose.Schema({
    name: {
        type: String,
        required: [true, "Subscription Name is required"],
        trim: true,
        minLength: 2,
        maxLength: 100,
    },
    price: {
        type: Number,
        required: [true, "Subscription Price is required"],
        min: [0, "Price must be greater than 0"]
    },
    currency: {
        type: String,
        enum: ["USD", "EUR", "GBP", "VND"],
        default: "VND"
    },
    frequency: {
        type: String,
        enum: ["daily", "weekly", "monthly", "yearly"],
        default: "monthly"
    },
    category: {
        type: String,
        enum: ["sport", "new", "entertainment", "lifestyle", "technology", "finance", "other"],
        required: true,
    },
    paymentMethod: {
        type: String,
        required: true,
        trim: true
    },
    status: {
        type: String,
        enum: ["active", "cancelled", "expired"],
        default: "active"
    },
    startDate: {
        type: Date,
        required: true,
        // Validate input for this field
        validate: {
            validator: (value) => value <= new Date(),
            message: "Start Date must be in the past"
        }
    },
    renewalDate: {
        type: Date,
        validate: {
            // Using function keyword to enable context of this
            validator: function (value) {
                return value > this.startDate
            },
            message: "Renewal Date must be after the Start Date"
        }
    },
    user: {
        type: mongoose.Schema.Types.ObjectId,
        required: true,
        ref: "User",
        index: true, // Index to optimize query for this field
    }
}, { timestamps: true })

// Middleware before save document
subscriptionSchema.pre("save", function (next) {
    // Auto calculate renewal date if missing
    if (!this.renewalDate) {
        const renewalPeriods = {
            daily: 1,
            weekly: 7,
            monthly: 30,
            yearly: 365,
        }

        this.renewalDate = new Date(this.startDate)
        this.renewalDate.setDate(this.renewalDate.getDate() + renewalPeriods[this.frequency])
    }

    // Auto update status if renewal date has passed
    if (this.renewalDate < new Date()) {
        this.status = "expired"
    }

    next()
})


const Subscription = mongoose.model("Subscription", subscriptionSchema)

export default Subscription
