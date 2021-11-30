import {useState} from 'react';
import { Link ,useNavigate} from "react-router-dom";

const DefaultLayout = (props) => {
    const [isLodingLogout,setIsLoadingLogout] = useState(false);
    const navigate = useNavigate();

    function onLogut(){
        setIsLoadingLogout(true)

        window.$axios.post("/logout")
        .then(() => {
            window.$toastr("Success","Berhasil Keluar")            
            setIsLoadingLogout(false)
            localStorage.removeItem('user-token');            
            props.setUser(false);
            navigate('/')
        })
        .catch(err => {            
            console.log(err)
            setIsLoadingLogout(false)
            window.$globalErrorToaster(window.$toastr,err)        
        })     
    }

    return (
        <>
            <div>
                <li>
                    <Link to="/">Home</Link>
                </li>

                {props.user && (
                    <>
                        <li>
                            <Link to="/profil">
                                Profil
                            </Link>
                        </li>                
                        <li>
                            <Link to="/data">
                                Data
                            </Link>
                        </li>
                        <li>
                            <Link to="#" onClick={onLogut}>
                                { isLodingLogout ? '...' : 'Keluar' }
                            </Link>
                        </li>
                    </>
                    )
                }   

                {!props.user && (
                    <>
                        <li>
                            <Link to="/signin">
                                Signin
                            </Link>
                        </li>
                        <li>
                            <Link to="/signup">
                                Signup
                            </Link>
                        </li>
                        <li>
                            <Link to="/forgot-password">
                                Forgot Password
                            </Link>
                        </li>
                    </>
                    )
                }
            </div>

            <div className="conatiner">
                {props.children}
            </div>
        </>
    )
}

export default DefaultLayout;