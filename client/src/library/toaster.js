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
		progress: undefined,
	};
		
	(action === 'error' ? window.$toast.error(args,options) : window.$toast.success(args,options))
}