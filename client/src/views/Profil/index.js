import { useEffect } from "react";
import DefaultLayout from "../../layouts/default";

const Profil = () => {
    useEffect(() => {
        window.$axios.get("/p").then((res) => {
            console.log(res);
        });
    },[]);

    return (
        <DefaultLayout>
            <h1>Profil</h1>
        </DefaultLayout>
    )
}

export default Profil;