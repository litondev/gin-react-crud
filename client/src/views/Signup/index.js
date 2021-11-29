import React, { useState } from 'react';
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

const Signup = () => {
    const [form] = useState({
        name : '',
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
 
 export default Signup;