import React, { useState } from 'react';
import { Navigate,Link} from "react-router-dom";
import { Formik, Form, Field, ErrorMessage} from 'formik';
import * as Yup from 'yup';

const ForgotPasswordSchema = Yup.object()
    .shape({      
        email: Yup.string()
            .email('Invalid email')
            .required('Required'),
    });

const ForgotPassword = (props) => {
    const [form] = useState({
        email : '', 
    })

    const onSubmit = (values,{setSubmitting}) => {            
        window.$axios.post("/auth/forgot-password",values)        
        .then(res => {
            window.$toastr("Success","Berhasil Kirim Ke Email")                    
        })        
        .catch(err => {         
            console.log(err)
            window.$globalErrorToaster(window.$toastr,err)        
        })
        .finally(() => {
            setSubmitting(false)   
        });
    }

    if(props.user){
        return <Navigate to="/" />
    }

    return (
        <div>
            <h1>ForgotPassword</h1>

            <Formik
                initialValues={form}
                validationSchema={ForgotPasswordSchema}
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
                            <Link to="/signin">Signin</Link>
                        </div>
                    </Form>
                    )
                }
            </Formik>
        </div>
    );
}
 
 export default ForgotPassword;