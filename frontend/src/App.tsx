import { RouterProvider, createBrowserRouter } from "react-router-dom";
import "./index.css";
import AuthPage from "./pages/AuthPage";
import Dashboard from "./pages/Dashboard";
import { InputOtpPage } from "./pages/OtpPage";
import ErrorPage from "./pages/ErrorPage";

const router = createBrowserRouter([
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
		element: <ErrorPage />
	}
]);

function App() {
	return <RouterProvider router={router} />;
}

export default App;
