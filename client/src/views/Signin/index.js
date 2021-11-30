import React, { useState } from 'react';
import { useNavigate ,Navigate,Link} from "react-router-dom";
import { Formik, Form, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';

const SigninSchema = Yup.object()
    .shape({      
        password: Yup.string()
            .min(8, 'Too Short!')
            .max(50, 'Too Long!')
            .required('Required'),
        email: Yup.string()
            .email('Invalid email')
            .required('Required'),
    });

const Signin = (props) => {

    const navigate = useNavigate();

    const [form] = useState({
        email : '', 
        password : ''
    })

    const onSubmit = (values,{setSubmitting}) => {            
        window.$axios.post("/auth/signin",values)
        .then(res => {
            localStorage.setItem('user-token',res.data.access_token);            
            return window.$axios.get("/me");           
        })
        .then(res => {
            setSubmitting(false)
            props.setUser(res.data);
            window.$toastr("Success","Berhasil Masuk")            
            navigate('/')
        })        
        .catch(err => {         
            setSubmitting(false)   
            console.log(err)
            window.$globalErrorToaster(window.$toastr,err)        
        })
    }

    if(props.user){
        return <Navigate to="/" />
    }

    return (
        <div>
            <h1>Signin</h1>

            <Formik
                initialValues={form}
                validationSchema={SigninSchema}
                onSubmit={onSubmit}>
                {({isSubmitting,resetForm}) => (                
                    <Form>
                        <div>
                            <Field 
                                type="email" 
                                name="email"
                                placeholder="Email . . ." />

                            <ErrorMessage  
                                name="email" 
                                component="div" />
                        </div>

                        <div>
                            <Field 
                                type="password" 
                                name="password" 
                                placeholder="Password . . ."/>

                            <ErrorMessage  
                                name="password" 
                                component="div" />
                        </div>

                        <div>
                            <button 
                                type="submit" 
                                disabled={isSubmitting}>
                                {isSubmitting ? '...' : 'Submit'}
                            </button>
                            <button type="reset"
                                onClick={resetForm}>
                                Reset
                            </button>
                        </div>

                        <div>
                            <Link to="/forgot-password">Forgot Password</Link>
                            <br/>
                            <Link to="/signup">Daftar</Link>
                        </div>
                    </Form>
                    )
                }
            </Formik>
        </div>
    );
}
 
 export default Signin;