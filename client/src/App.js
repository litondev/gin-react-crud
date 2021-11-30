import React,{useState,useEffect} from 'react';
import { BrowserRouter,Route,Routes} from "react-router-dom";
import MyRoutes from "./routes/index.js";

const App = (props) => {
    const [user,setUser] = useState(false);

    const ToastContainer = window.$ToastContainer;

    const setAuth = (user) => {      
      setUser(user)
    }

     useEffect(() => {
       if(props.auth){
        setAuth(props.auth)
       }
    },[props.auth])

    return (        
        <BrowserRouter>   
          <React.Suspense fallback={ <span>Loading</span> }>    
            <Routes>
            {
              MyRoutes.map((route,indexRoute) => {
                return <Route
                  path={ route.path }
                  key={ indexRoute }
                  element={
                    <route.component 
                     user={user}
                     setUser={setAuth}/>                  
                  } />
              })
            }           
            </Routes>
          </React.Suspense>
          <ToastContainer/>
        </BrowserRouter>
    )
}

export default App;