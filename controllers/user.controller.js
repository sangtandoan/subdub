import User from "../models/user.model.js";

export const getAllUsers = async (req, res) => {
	const users = await User.find({}).lean();

	res.status(200).json({
		success: true,
		message: "Get all users successfully!",
		data: {
			users,
		},
	});
};

export const getUser = async (req, res) => {
	if (req.user._id.toString() !== req.params.id) {
		const error = new Error("Unauthorized!");
		error.statusCode = 401;
		throw error;
	}

	const user = await User.findOne({ _id: req.params.id })
		.select("-password")
		.lean();
	if (!user) {
		const error = new Error("There is no user of this id!");
		error.statusCode = 404;
		throw error;
	}

	res.status(200).json({
		success: true,
		message: "Get a user successfully!",
		data: {
			user,
		},
	});
};
