import React, { useState } from 'react';
import { useNavigate,Navigate, useSearchParams} from "react-router-dom";
import { Formik, Form, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';

const ResetPasswordSchema = Yup.object()
    .shape({      
        password_confirm : Yup.string()
            .min(8, 'Too Short!')
            .max(50, 'Too Long!')
            .required('Required'),
        new_password: Yup.string()
            .min(8, 'Too Short!')
            .max(50, 'Too Long!')
            .required('Required'),
        email: Yup.string()
            .email('Invalid email')
            .required('Required'),
    });

const ResetPassword = (props) => {

    const navigate = useNavigate();

    const [searchParams] = useSearchParams();

    const [form] = useState({
        password_confirm : '',
        email : searchParams.get("email") || "user1@gmail.com", 
        new_password : '',
        token : searchParams.get("token") || ""
    })

    const onSubmit = (values,{setSubmitting}) => {            
        window.$axios.post("/auth/reset-password",values)
        .then(() => {
            setSubmitting(false)
            window.$toastr("success","Berhasil Reset Password")
            navigate('/signin')
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
            <h1>ResetPassword</h1>

            <Formik
                initialValues={form}
                validationSchema={ResetPasswordSchema}
                onSubmit={onSubmit}>
                {({isSubmitting,resetForm}) => (                
                    <Form>                                      
                        <div>
                            <Field 
                                type="password" 
                                name="new_password" 
                                placeholder="Password . . ."/>

                            <ErrorMessage  
                                name="new_password" 
                                component="div" />
                        </div>

                        <div>
                            <Field 
                                type="password" 
                                name="password_confirm" 
                                placeholder="Password  Konfirmasi . . ."/>

                            <ErrorMessage  
                                name="password_confirm" 
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
                    </Form>
                    )
                }
            </Formik>
        </div>
    );
}
 
 export default ResetPassword;