import DefaultLayout from "../../layouts/default";
import { Navigate} from "react-router-dom";

const Data = (props) => {
    if(!props.user){
        return <Navigate to="/signin" />
    }

    return (
        <DefaultLayout {...props}>
            <h1>Data</h1>
        </DefaultLayout>
    )
}

export default Data;