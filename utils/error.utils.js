// HOF to handle throw, unhandled, syntax errors that does not call next()
const catchErrors = (fn) => {
	return (req, res, next) => {
		return fn(req, res, next).catch(next);
	};
};

export { catchErrors };
