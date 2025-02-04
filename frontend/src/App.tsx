import { RouterProvider, createBrowserRouter } from "react-router-dom";
import "./index.css";
import AuthPage from "./pages/AuthPage";
import Dashboard from "./pages/Dashboard";
import UploadBankStatementPage from "./pages/UploadBankStatementPage";

const router = createBrowserRouter([
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
