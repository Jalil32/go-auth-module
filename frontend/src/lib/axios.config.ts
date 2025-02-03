import axios from 'axios';


const api = axios.create({
	baseURL: '/api', // Uses the proxy from Vite
	headers: {
		'Content-Type': 'application/json',
	}
});

export default api;
