import "./index.css";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import LoginPage from "./pages/login";
import ChartComponent from "./components/chart";
import RegisterPage from "./pages/register";

const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/register",
    element: <RegisterPage />,
  },
  {
    path: "/",
    element: <ChartComponent></ChartComponent>,
  },
]);

function App() {
  return <RouterProvider router={router} />;
}

export default App;
