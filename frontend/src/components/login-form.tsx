import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import axios from "axios";
import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";

interface LoginFormProps extends React.ComponentPropsWithoutRef<"form"> {
	toggleMode: () => void;
}

type Inputs = {
	email: string;
	password: string;
};

export function LoginForm({ className, toggleMode, ...props }: LoginFormProps) {
	const navigate = useNavigate();
	const {
		register,
		handleSubmit,
		formState: { errors },
	} = useForm<Inputs>();
	const [submitError, setSubmitError] = useState<string | null>(null);
	const [loading, setLoading] = useState(false);

	const handleGoogleAuth = async () => {
		window.location.href = "/api/auth/google";
	};

	const onSubmit: SubmitHandler<Inputs> = async (data) => {
		setLoading(true);
		try {
			const res = await axios.post("/api/auth/login", data);

			if (res.status === 200) {
				navigate("/dashboard");
			}
		} catch (error) {
			if (axios.isAxiosError(error)) {
				if (error.response?.status === 401) {
					if (
						error.response?.data?.message ===
						"User is not verified."
					) {
						setSubmitError(
							"Please verify your email to access the application",
						);
						navigate("/auth/otp", { state: { email: data.email } });
					} else {
						setSubmitError("Invalid email or password");
					}
				} else {
					setSubmitError("Something went wrong. Please try again");
				}
			} else {
				setSubmitError("An unexpected error occurred");
			}
		} finally {
			setLoading(false);
		}
	};

	return (
		<form
			onSubmit={handleSubmit(onSubmit)}
			className={cn("flex flex-col gap-6", className)}
			{...props}
		>
			<div className="flex flex-col items-center gap-2 text-center">
				<h1 className="text-2xl font-bold">Login to Wealth Scope</h1>
				<p className="text-balance text-sm text-muted-foreground">
					Enter your email and password below to login
				</p>
			</div>
			<div className="grid gap-6">
				<div className="grid gap-2">
					<Label htmlFor="email">Email</Label>
					<Input
						autoComplete="email"
						id="email"
						type="email"
						placeholder="ws@gmail.com"
						{...register("email", {
							required: "Email is required",
							pattern: {
								value: /\S+@\S+\.\S+/,
								message: "Email address is invalid",
							},
						})}
					/>
					{errors.email && (
						<span className="text-sm text-red-500">
							{errors.email.message}
						</span>
					)}
				</div>
				<div className="grid gap-2">
					<div className="flex items-center">
						<Label htmlFor="password">Password</Label>
						<a
							href="/forgotpassword"
							className="ml-auto text-sm underline-offset-4 hover:underline"
						>
							Forgot your password?
						</a>
					</div>
					<Input
						id="password"
						type="password"
						{...register("password", {
							required: "Password is required",
						})}
					/>
					{errors.password && (
						<span className="text-sm text-red-500">
							{errors.password.message}
						</span>
					)}
					{submitError && (
						<span className="text-sm text-red-500">
							{submitError}
						</span>
					)}
				</div>
				<Button type="submit" disabled={loading} className="w-full">
					{loading ? "Logging in..." : "Login"}
				</Button>
				<div className="relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t after:border-border">
					<span className="relative z-10 bg-background px-2 text-muted-foreground">
						Or continue with
					</span>
				</div>
				<Button
					onClick={handleGoogleAuth}
					type="button"
					variant="outline"
					className="w-full"
				>
					Login with Google
				</Button>
			</div>
			<div className="text-center text-sm">
				Don&apos;t have an account?{" "}
				<button
					type="button"
					onClick={toggleMode}
					className="underline underline-offset-4"
				>
					Signup
				</button>
			</div>
		</form>
	);
}
