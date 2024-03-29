import DefaultLayout from "../../layouts/default";
import { Navigate} from "react-router-dom";

const Product = (props) => {
    if(!props.user){
        return <Navigate to="/signin" />
    }
    
    return (
        <DefaultLayout
            {...props}>
            <h1>Product</h1>
        </DefaultLayout>
    )
}

export default Product;