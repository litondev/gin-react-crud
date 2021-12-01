import DefaultLayout from "../../layouts/default";
import { Navigate} from "react-router-dom";
import { useEffect,useState } from "react";
import { Formik, Form, Field, ErrorMessage,useFormik} from 'formik';
import * as Yup from 'yup';


const DataSchema = Yup.object()
    .shape({      
        name : Yup.string()
            .required('Required')
            .max(50, 'Too Long!'),
        phone : Yup.string() 
    });

const Data = (props) => {
    // const [form,setForm] = useState({
    //     id : 0,
    //     name : '',
    //     phone : ''
    // })

    const [data,setData] =  useState({
        indexLoadingDelete : '',
        isLoadingDelete : false,
        isEditable : false,
        isLoadPage : true,
        search : '',  
        per_page : 10,
        page : 1,
        total_page : 0,
        paginate : []
    })

    const { values, errors, handleChange, setFieldValue,setValues } = useFormik({
        initialValues: {
            id : 0,
            name : 's',
            phone : ''
        }
    });

    const onLoad = (isSearch = false) => {
        setData(prevState => {
            return {
                ...prevState,
                isLoadPage : true,
                page : isSearch ? 1 : data.page 
            }
        })
        
        window.$axios.get("/data",{
            params : {
                search : data.search,
                per_page : data.per_page,
                page : isSearch ? 1 : data.page 
            },
        })
        .then(res => {
            setData(prevState => {
                return {
                    ...prevState,
                    isLoadPage : false,
                    paginate : res.data.data,   
                    total_page : res.data.total_page             
                }
            })
            console.log(res);
        })
        .catch(err => {
            console.log(err);
        })
    }

    const onDelete = (item) => {
        if(data.isLoadingDelete) return 

        let sure = window.confirm("Anda yakin");

        if(!sure) return

        setData((prevState) => {
            return {
                ...prevState,
                isLoadingDelete : true,
                indexLoadingDelete : item.id
            }
        })

        window.$axios.delete("/data/"+item.id)
        .then(() => {
            window.$toastr("success","Berhasil Mengahapus Data")
            onLoad()
        })
        .catch(err => {
            console.log(err)
            window.$globalErrorToaster(window.$toastr,err)     
        })
        .finally(() => {
            setData((prevState) => {
                return {
                    ...prevState,
                    isLoadingDelete : false,
                    indexLoadingDelete : ''
                }
            })
        })
    }

    const onAdd = () => {
        setData(prevState => {
            return {
                ...prevState,
                isEditable : false,
                isShowForm : true
            }
        })

        setValues({
            id : 0,
            name : '',
            phone : ''
        })        
    }

    const onEdit = (item) => {
        setData(prevState => {
            return {
                ...prevState,
                isEditable : true,
                isShowForm : true
            }
        })
        
        setValues({...item})        
 
    }
    
    const onResetForm = () => {
        setValues({
            id : 0,
            name : '',
            phone : ''
        })    
    }

    const onSubmit = (values,{setSubmitting}) => {            
        let url = "/data";

        if(data.isEditable){
            url += "/"+values.id;
        }

        window.$axios({
            method : data.isEditable ? "put" : "post",
            url : url,
            data : values
        })
        .then(() => {        
            window.$toastr("success","Berhasil")
            onLoad()
        })
        .catch(err => {        
            console.log(err)
            window.$globalErrorToaster(window.$toastr,err)        
        })
        .finally(() => {
            setSubmitting(false)
        })            
    }

    const onSearch = (event) => {
        if(event.key === 'Enter'){
            setData({
                ...data,
                search : event.target.value,
            })                      
        }
    }

    const onPage = (isNext) => {
        let page = data.page;

        if(isNext){
            page = page+1;
        }else{
            page = page-1;
        }

        setData({
            ...data,
            page : page
        })
    }

    const onExport = (type) => {
        window.$axios({
            url : "/data/export/"+type,
        })
        .then(res => {             
            var byteString = atob(res.data);
            var ab = new ArrayBuffer(byteString.length);
            var ia = new Uint8Array(ab);
            
            for (var i = 0; i < byteString.length; i++) {
                ia[i] = byteString.charCodeAt(i);
            }

            if(type == "excel"){
                var blob = new Blob([ab], { type: 'application/wps-office' });            
            }else{
                var blob = new Blob([ab], { type: 'application/pdf' });            
            }

            console.log(res.data);
            const link = document.createElement('a');
            link.href = window.URL.createObjectURL(blob);
            link.setAttribute('download', 'report.'+ (type == "excel" ? "xlsx" : "pdf"));
            document.body.appendChild(link);
            link.click();            
        })
        .catch(err => {
            console.log(err)
            window.$globalErrorToaster(window.$toastr,err)      
        })
    }


    const onSubmitExcel = (event) => {
        event.preventDefault();

        let formData = new FormData(document.getElementById("upload"))
        // formData.append("_method","PUT");

        window.$axios.post("/data/import",formData)
        .then(() => {
            window.$toastr("success","Berhasil")
            onLoad()
        })       
        .catch(err => {
            console.log(err)
            window.$globalErrorToaster(window.$toastr,err) 
        })
        .finally(() => {
            // end loading
        })
    }

    useEffect(() => {
        if(data.isLoadPage == false){
            onLoad()                        
        }
    },[data.page])

    useEffect(() => {
        if(data.isLoadPage == false){
            onLoad(true)            
        }
    },[data.search])

    useEffect(() => {
        onLoad();
    }, [])

    if(!props.user){
        return <Navigate to="/signin" />
    }
    
    return (
        <DefaultLayout {...props}>
            <h1>Data</h1>

            <div>
                <form id="upload">                    
                    <input type="file" name="excel" id="excel"/>
                    <button onClick={onSubmitExcel}>Upload</button>
                </form>
            </div>

            {data.isShowForm && <>
                <div>
                <Formik
                    initialValues={values}
                    enableReinitialize={true}
                    validationSchema={DataSchema}
                    onSubmit={onSubmit}>
                    {({isSubmitting,resetForm}) => (                
                        <Form>
                           
                            
                            <div>
                                <Field 
                                    type="text" 
                                    name="name"
                                    placeholder="Name . . ." />

                                <ErrorMessage  
                                    name="name" 
                                    component="div" />
                            </div>

                            <div>
                                <Field 
                                    type="text" 
                                    name="phone"
                                    placeholder="Phone . . ." />

                                <ErrorMessage  
                                    name="phone" 
                                    component="div" />
                            </div>

                            <div>
                                <button 
                                    type="submit" 
                                    disabled={isSubmitting}>
                                    {isSubmitting ? '...' : 'Submit'}
                                </button>
                                <button type="reset"
                                    onClick={onResetForm}>
                                    Reset
                                </button>
                            </div>                          
                        </Form>
                        )
                    }
                </Formik>
                </div>
            </>}

            <div>
                <input type="text" name="search"
                    onKeyPress={onSearch}/>
                <br/>
                <button onClick={onAdd}>Add</button>
                <button onClick={() => onExport('pdf')}>ExportPdf</button>
                <button onClick={() => onExport('excel')}>ExportExcel</button>
            </div>

            <table>
                <thead>
                    <tr>
                        <td>Id</td>
                        <td>Nama</td>
                        <td>Phone</td>
                        <td>Opsi</td>
                    </tr>                
                </thead>
                <tbody>
                    { data.paginate.map((item,index) => (
                        <tr key={index}>
                            <td>{item.id}</td>
                            <td>{item.name}</td>
                            <td>{item.phone || '-'}</td>
                            <td>
                                <button onClick={() => onDelete(item)}>
                                    {data.isLoadingDelete && data.indexLoadingDelete  == item.id ? '...' : 'Delete' }
                                </button>

                                <button onClick={() => onEdit(item)}>
                                    Edit
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>

            <div>
                { data.page > 1 &&
                    <button onClick={() => onPage(false)}>Before</button>
                }
                <br/>
                { data.page < data.total_page && 
                    <button onClick={() => onPage(true)}>Next</button>                
                }
            </div>
        </DefaultLayout>
    )
}

export default Data;