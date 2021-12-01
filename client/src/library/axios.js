import axios from "axios";
import moment from "moment";

axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';
axios.defaults.baseURL = process.env.REACT_APP_API_URL;

axios.interceptors.request.use(config => { 	
  	config.headers['Authorization'] = localStorage.getItem('user-token') ? "Bearer "+localStorage.getItem("user-token") : null
	return config;
});

function refreshToken(){	
	console.log("Refresh Token");

	if(localStorage.getItem("user-token")){				
		let token = localStorage.getItem("user-token").split(".");
		const decodedJWT = JSON.parse(atob(token[1]));		
		const dateNow = new Date()
		const miliseconds = dateNow.getTime() / 1000
		
		// EXPIRED TOKEN z
		console.log(moment(decodedJWT.exp*1000).format("hh:mm:ss"))
		// TIME TO REFRESH TOKEN
		console.log(moment(decodedJWT.exp*1000).subtract(30,'minutes').format("hh:mm:ss"))
		// TIME NOW
		console.log(moment(miliseconds*1000).format("hh:mm:ss"))

		if(moment(miliseconds*1000).isAfter(moment(decodedJWT.exp*1000).subtract(30,'seconds'))){
			let sendRefreshToken = window.$axios;

			sendRefreshToken.defaults.headers = {
				Authorization : "Bearer " + localStorage.getItem("user-token"),
			}		

			return sendRefreshToken.post("/refresh-token")
			.then(res => {
				localStorage.setItem('user-token',res.data.access_token);            
			})
			.catch(err => {				
				localStorage.removeItem("user-token");
				window.location = "/signin";
			})
		}
	}
}

axios.interceptors.response.use(res => {	
	if(!res.data.access_token){	
		refreshToken();
	}

	return res;
});

window.$axios = axios;
