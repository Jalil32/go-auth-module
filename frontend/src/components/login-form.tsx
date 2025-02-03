import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ChangeEvent, FormEvent, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

interface LoginFormProps extends React.ComponentPropsWithoutRef<"form"> {
	toggleMode: () => void;
}

export function LoginForm({
	className,
	toggleMode,
	...props
}: LoginFormProps) {
	const navigate = useNavigate()
	const [errors, setErrors] = useState<{ email?: string; password?: string, submitError?: string }>({});
	const [formData, setFormData] = useState({
		email: '',
		password: ''
	});


	const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
		const { name, value } = e.target;
		setFormData({
			...formData,
			[name]: value,
		});
	};

	const validateForm = () => {
		const newErrors: { email?: string; password?: string } = {};

		if (!formData.email) {
			newErrors.email = "Email is required";
		} else if (!/\S+@\S+\.\S+/.test(formData.email)) {
			newErrors.email = "Email address is invalid";
		}

		if (!formData.password) {
			newErrors.password = "Password is required";
		}

		setErrors(newErrors);
		return Object.keys(newErrors).length === 0;
	};

	const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		const isValid = validateForm();
		if (isValid) {
			axios.post('api/auth/login', formData).then(res => {
				setErrors({})
				if (res.status == 200) {
					console.log("Login Successful")
					navigate('/')
				}
			}).catch((error) => {
				const newErrors: { submitError?: string } = {};
				if (error.response && error.response.status === 401) {
					newErrors.submitError = "Invalid email or password"
					setErrors(newErrors)
				} else {
					newErrors.submitError = "Something went wrong. Please try again."
				}
			})
		}
	};

	return (
		<form onSubmit={handleSubmit} className={cn("flex flex-col gap-6", className)} {...props}>
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
						id="email"
						type="email"
						name="email"
						placeholder="jasonrmillman@gmail.com"
						value={formData.email}
						onChange={handleChange}
						required
					/>
					{errors.email && <span className="text-sm text-red-500">{errors.email}</span>}
				</div>
				<div className="grid gap-2">
					<div className="flex items-center">
						<Label htmlFor="password">Password</Label>
						<a
							href="test"
							className="ml-auto text-sm underline-offset-4 hover:underline"
						>
							Forgot your password?
						</a>
					</div>
					<Input
						id="password"
						type="password"
						name="password"
						value={formData.password}
						onChange={handleChange}
						required
					/>
					{errors.password && <span className="text-sm text-red-500">{errors.password}</span>}
					{errors.submitError && <span className="text-sm text-red-500">{errors.submitError}</span>}
				</div>
				<Button type="submit" className="w-full">
					Login
				</Button>
				<div className="relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t after:border-border">
					<span className="relative z-10 bg-background px-2 text-muted-foreground">
						Or continue with
					</span>
				</div>
				<Button type="button" variant="outline" className="w-full">
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
