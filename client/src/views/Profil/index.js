import DefaultLayout from "../../layouts/default";
import { Navigate} from "react-router-dom";

const Profil = (props) => {
    if(!props.user){
        return <Navigate to="/signin" />
    }
    
    return (
        <DefaultLayout  {...props}>
            <h1>Profil</h1>
        </DefaultLayout>
    )
}

export default Profil;