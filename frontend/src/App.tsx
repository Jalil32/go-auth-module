import { RouterProvider, createBrowserRouter } from "react-router-dom";
import "./index.css";
import AuthPage from "./pages/AuthPage";
import Dashboard from "./pages/Dashboard";
import ErrorPage from "./pages/ErrorPage";
import { Landing } from "./pages/Landing";
import { InputOtpPage } from "./pages/OtpPage";

const router = createBrowserRouter([
	{
		path: "/",
		element: <Landing />,
	},
	{
		path: "/auth",
		element: <AuthPage />,
	},
	{
		path: "/auth/otp",
		element: <InputOtpPage />,
	},
	{
		path: "/dashboard",
		element: <Dashboard />,
	},
	{
		path: "*",
		element: <ErrorPage />,
	},
]);

function App() {
	return <RouterProvider router={router} />;
}

export default App;
