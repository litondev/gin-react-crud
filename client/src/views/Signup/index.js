import React, { useState } from 'react';
import { useNavigate ,Navigate} from "react-router-dom";
import { Formik, Form, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';

const SignupSchema = Yup.object()
    .shape({      
        name : Yup.string()
            .required('Required')
            .max(50, 'Too Long!'),
        password: Yup.string()
            .min(8, 'Too Short!')
            .max(50, 'Too Long!')
            .required('Required'),
        email: Yup.string()
            .email('Invalid email')
            .required('Required'),
    });

const Signup = (props) => {
    const navigate = useNavigate();
  
    const [form] = useState({
        name : '',
        email : '', 
        password : ''
    })

    const onSubmit = (values,{setSubmitting}) => {            
        window.$axios.post("/auth/signup",values)
        .then(() => {
            setSubmitting(false)
            window.$toastr("success","Berhasil Membuat User")
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
            <h1>Signup</h1>

            <Formik
                initialValues={form}
                validationSchema={SignupSchema}
                onSubmit={onSubmit}>
                {({isSubmitting,resetForm}) => (                
                    <Form>
                        <div>
                            <Field
                                type="name"
                                name="name"
                                placeholder="Name . . ." />

                            <ErrorMessage  
                                name="name" 
                                component="div" />
                        </div>
                        
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
                    </Form>
                    )
                }
            </Formik>
        </div>
    );
}
 
 export default Signup;