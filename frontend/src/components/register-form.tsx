import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import axios from "axios";
import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";

type Inputs = {
	email: string;
	firstName: string;
	lastName: string;
	password: string;
	confirmPassword: string;
};

interface RegisterFormProps extends React.ComponentPropsWithoutRef<"form"> {
	toggleMode: () => void;
}

export function RegisterForm({
	className,
	toggleMode,
	...props
}: RegisterFormProps) {
	const navigate = useNavigate();
	const {
		register,
		handleSubmit,
		watch,
		formState: { errors },
	} = useForm<Inputs>();
	const [submitError, setSubmitError] = useState<string | null>(null);
	const [loading, setLoading] = useState(false);
	const password = watch("password");

	const handleGoogleAuth = async () => {
		window.location.href = "/api/auth/google";
		// will need backend to redirect to custom error page
	};

	const onSubmit: SubmitHandler<Inputs> = async (data) => {
		setLoading(true);
		try {
			const res = await axios.post("/api/auth/register", data);

			if (res.status === 201) {
				navigate("/auth/otp", { state: { email: data.email } });
			}
		} catch (error) {
			if (axios.isAxiosError(error)) {
				if (error.response?.status === 400) {
					setSubmitError("Invalid request");
				} else if (error.response?.status === 401) {
					setSubmitError("Invalid email or password");
				} else if (error.response?.status === 409) {
					setSubmitError(
						"An account with this email already exists. Please sign in",
					);
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
				<h1 className="text-2xl font-bold">Signup to Wealth Scope</h1>
				<p className="text-balance text-sm text-muted-foreground">
					Enter your information below to signup
				</p>
			</div>
			<div className="grid gap-6">
				<div className="grid gap-2">
					<Label htmlFor="firstName">First Name</Label>
					<Input
						autoComplete="given-name webauthn"
						id="firstName"
						type="text"
						placeholder="John"
						{...register("firstName", {
							required: "First name is required",
							pattern: {
								value: /^[A-Za-z'-]+$/,
								message: "First name is invalid",
							},
						})}
					/>
					{errors.firstName && (
						<span className="text-sm text-red-500">
							{errors.firstName.message}
						</span>
					)}
				</div>
				<div className="grid gap-2">
					<Label htmlFor="lastName">Last Name</Label>
					<Input
						autoComplete="family-name"
						id="lastName"
						type="text"
						placeholder="Doe"
						{...register("lastName", {
							required: "Last name is required",
							pattern: {
								value: /^[A-Za-z'-]+$/,
								message: "Last name is invalid",
							},
						})}
					/>
					{errors.lastName && (
						<span className="text-sm text-red-500">
							{errors.lastName.message}
						</span>
					)}
				</div>
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
					</div>
					<Input
						id="password"
						type="password"
						{...register("password", {
							required: "Password is required",
							pattern: {
								value: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/,
								message:
									"Password must be at least 8 characters long and include one uppercase letter, one lowercase letter, one number, and one special character",
							},
						})}
					/>
					{errors.password && (
						<span className="text-sm text-red-500">
							{errors.password.message}
						</span>
					)}
				</div>
				<div className="grid gap-2">
					<div className="flex items-center">
						<Label htmlFor="confirmPassword">
							Confirm Password
						</Label>
					</div>
					<Input
						id="confirmPassword"
						type="password"
						{...register("confirmPassword", {
							required: "Please confirm your password",
							validate: (value) =>
								value === password || "Passwords do not match",
						})}
					/>
					{errors.confirmPassword && (
						<span className="text-sm text-red-500">
							{errors.confirmPassword.message}
						</span>
					)}
					{submitError && (
						<span className="text-sm text-red-500">
							{submitError}
						</span>
					)}
				</div>
				<Button type="submit" disabled={loading} className="w-full">
					{loading ? "Signing up..." : "Signup"}
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
					Signup with Google
				</Button>
			</div>
			<div className="text-center text-sm">
				Already have an account?{" "}
				<button
					type="button"
					onClick={toggleMode}
					className="underline underline-offset-4"
				>
					Login
				</button>
			</div>
		</form>
	);
}
