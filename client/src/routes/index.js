import React from "react";

const myRoutes = [
	{
		path : "/signin",
		component : React.lazy(() =>  import('../views/Signin'))
	},
	{
		path : "/signup",
		component : React.lazy(() =>  import('../views/Signup'))
	},
	{
		path : "/forgot-password",
		component : React.lazy(() =>  import('../views/ForgotPassword'))
	},
	{
		path : '/data',
		component : React.lazy(() => import('../views/Data'))
	},
	{
		path : "/product",
		component : React.lazy(() => import('../views/Product'))
	},
	{
		path : '/profil',
		component : React.lazy(() => import('../views/Profil'))
	},
	{
		path : '/',
		component : React.lazy(() => import('../views/Home'))
	},
	{
		path : '*',
		component : React.lazy(() => import('../views/Error/P404'))
	}
];

export default myRoutes;