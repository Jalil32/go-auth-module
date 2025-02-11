import { RouterProvider, createBrowserRouter } from "react-router-dom";
import "./index.css";
import AuthPage from "./pages/AuthPage";
import Dashboard from "./pages/Dashboard";
import { Landing } from "./pages/Landing";
import UploadBankStatementPage from "./pages/UploadBankStatementPage";

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
		path: "/dashboard",
		element: <Dashboard />,
	},
	{
		path: "/upload-bank-statement",
		element: <UploadBankStatementPage />,
	},
]);

function App() {
	return <RouterProvider router={router} />;
}

export default App;
