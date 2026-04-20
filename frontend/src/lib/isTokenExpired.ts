import { jwtDecode } from 'jwt-decode';

export const isTokenExpired = (token: string): boolean => {
	try {
		const decoded = jwtDecode(token);

		const currentTime = Date.now() / 1000;

		if (decoded.exp && decoded.exp < currentTime) {
			return true;
		}

		return false;
	} catch (error) {
		console.error(error);

		return true;
	}
};
