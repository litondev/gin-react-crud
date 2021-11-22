import { Link } from "react-router-dom";

const DefaultLayout = (props) => {
    return (
        <>
            <div>
                <li>
                    <Link to="/">Home</Link>
                </li>
                <li>
                    <Link to="/profil">Profil</Link>
                </li>
                <li>
                    <Link to="/data">Data</Link>
                </li>
                <li>
                    <Link to="/signin">Signin</Link>
                </li>
                <li>
                    <Link to="/signup">Signup</Link>
                </li>
                <li>
                    <Link to="/forgot-password">Forgot Password</Link>
                </li>
            </div>
            <div class="conatiner">
                {props.children}
            </div>
        </>
    )
}

export default DefaultLayout;