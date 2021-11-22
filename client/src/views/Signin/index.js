import React, { useState } from 'react';
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

const Signin = () => {
    const [form] = useState({
        email : '', 
        password : ''
    })

    const onSubmit = (values,{setSubmitting}) => {            
        setTimeout(() => {
            alert(JSON.stringify(values, null, 2));
            setSubmitting(false);
        }, 400);        
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
                                Submit
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
 
 export default Signin;