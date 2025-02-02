import { RouterProvider, createBrowserRouter } from "react-router-dom";
import "./index.css";
import AuthPage from "./pages/AuthPage";
import Dashboard from "./pages/Dashboard";

const router = createBrowserRouter([
	{
		path: "/auth",
		element: <AuthPage />,
	},
	{
		path: "/dashboard",
		element: <Dashboard />,
	},
]);

function App() {
	return <RouterProvider router={router} />;
}

export default App;
