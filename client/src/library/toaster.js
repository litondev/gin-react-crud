import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

window.$toast = toast;
window.$ToastContainer = ToastContainer;
window.$toastr = (action,args) => {
 	let options = {
		position: "top-right",
		autoClose: 5000,
		hideProgressBar: true,
		closeOnClick: true,
		pauseOnHover: true,
		draggable: true,
		progress: true,
		theme : 'colored'
	};
		
	(action === 'error' ? window.$toast.error(args,options) : window.$toast.success(args,options))
}

window.$globalErrorToaster = ($toastErr,$err) => {
	if($err.response && $err.response.status === 422){
		$toastErr('error',$err.response.data.message || 'Terjadi Kesalahan');
	}else if($err.response && $err.response.status === 500){
		$toastErr('error',$err.response.data.message || 'Terjadi Kesalahan');
	}else if($err.response && $err.response.status === 503){
		$toastErr('error',"Maintenance");
	}else if($err.response && $err.response.status === 401){
		$toastErr('error',$err.response.data.message || 'Terjadi Kesalahan');
	}else{
		$toastErr('error','Terjadi Kesalahan');
	}  
}