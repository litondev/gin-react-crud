import React, { useState,useEffect } from 'react';
import { Navigate} from "react-router-dom";
import { Formik, Form, Field, ErrorMessage,useFormik } from 'formik';
import * as Yup from 'yup';
import DefaultLayout from "../../layouts/default";

const ProfilSchema = Yup.object()
    .shape({      
        name : Yup.string()
            .required('Required')
            .max(50, 'Too Long!'),
        password: Yup.string()
            .min(8, 'Too Short!')
            .max(50, 'Too Long!'),
        password_confirm: Yup.string()
            .min(8,'Too short!')
            .required('Required'),
        email: Yup.string()
            .email('Invalid email')
            .required('Required'),
    });

const Profil = (props) => {         
 
    const { values, errors, handleChange, setFieldValue,setValues } = useFormik({
        initialValues: {            
            photo : props.user.photo,
            name : props.user.name,
            email : props.user.email, 
            password : '',
            password_confirm : ''
        }
    });

    const onChangePhoto = (evt) => {
        if(!evt.target.files[0]){
            return false;
        } 

        if(!['image/jpeg','image/jpg','image/png'].includes(evt.target.files[0].type)){
            evt.target.value = "";              
            return false;
        }
    
        setValues((prevState) => ({
            ...prevState,
            photo : URL.createObjectURL(evt.target.files[0])
        }));
    }

    const onSubmitPhoto = (event) => {
        event.preventDefault();

        let formData = new FormData(document.getElementById("upload"))
        // formData.append("_method","PUT");

        window.$axios.post("/profil/upload",formData)
        .then(() => {
            return window.$axios.get("/me")
        })
        .then((res) => {
            document.getElementById("upload").reset();
            props.setUser(res.data.user)
            window.$toastr("success","Berhasil")
        })       
        .catch(err => {
            console.log(err)
            window.$globalErrorToaster(window.$toastr,err) 
        })
        .finally(() => {
            // end loading
        })
    }

    const onSubmit = (values,{setSubmitting}) => {            
        window.$axios.put("/profil/update",values)
        .then(() => {
            return window.$axios.get("/me")
        })
        .then((res) => {
            props.setUser(res.data.user)
            window.$toastr("success","Berhasil")
        })
        .catch(err => {
            console.log(err)
            window.$globalErrorToaster(window.$toastr,err)        
        })
        .finally(() => {
            setSubmitting(false)
        })
    }

    if(!props.user){
        return <Navigate to="/" />
    }

    return (
        <DefaultLayout {...props}>
            <h1>Profil</h1>

            <div>
                <form id="upload">
                    <img src={values.photo} 
                        width="100px"/>
                    <input type="file" name="photo" id="photo"
                        onChange={onChangePhoto}/>
                    <button onClick={onSubmitPhoto}>Upload</button>
                </form>
            </div>

            <Formik
                initialValues={values}
                validationSchema={ProfilSchema}
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
                            <Field 
                                type="password" 
                                name="password_confirm" 
                                placeholder="Password Konfirmasi . . ."/>

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
        </DefaultLayout>
    );
}
 
 export default Profil;