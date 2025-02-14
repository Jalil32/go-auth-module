import { GalleryVerticalEnd } from "lucide-react";
import { AnimatePresence, motion } from "motion/react";
import { useLocation, useNavigate } from "react-router-dom";
import { LoginForm } from "../components/login-form.tsx";
import { RegisterForm } from "../components/register-form.tsx";

export default function AuthPage() {
	const location = useLocation();
	const navigate = useNavigate();

	// Extract mode query parameter
	const searchParams = new URLSearchParams(location.search);
	const mode = searchParams.get("mode") || "login"; // Default to login

	// Function to toggle between login and register
	const toggleMode = () => {
		const newMode = mode === "register" ? "login" : "register";
		navigate(`/auth?mode=${newMode}`, { replace: true });
	};

	return (
		<div className="grid min-h-svh">
			<div className="flex flex-col gap-4 p-6 md:p-10">
				<div className="flex justify-center gap-2 md:justify-start">
					<a
						href="test"
						className="flex items-center gap-2 font-medium"
					>
						<div className="flex h-6 w-6 items-center justify-center rounded-md bg-primary text-primary-foreground">
							<GalleryVerticalEnd className="size-4" />
						</div>
						Wealth Scope.
					</a>
				</div>
				<div className="flex flex-1 items-center justify-center">
					<div className="w-full max-w-xs">
						<AnimatePresence mode="wait">
							{mode === "register" ? (
								<motion.div
									key="register"
									initial={{ opacity: 0, y: 10 }}
									animate={{ opacity: 1, y: 0 }}
									exit={{ opacity: 0, y: -10 }}
									transition={{ duration: 0.1 }}
								>
									<RegisterForm toggleMode={toggleMode} />
								</motion.div>
							) : (
								<motion.div
									key="login"
									initial={{ opacity: 0, y: 10 }}
									animate={{ opacity: 1, y: 0 }}
									exit={{ opacity: 0, y: -10 }}
									transition={{ duration: 0.1 }}
								>
									<LoginForm toggleMode={toggleMode} />
								</motion.div>
							)}
						</AnimatePresence>
					</div>
				</div>
			</div>
		</div>
	);
}
