import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ChangeEvent, FormEvent, useState } from "react";
import api from "@/lib/axios.config";

interface FormData {
	email: string,
	firstName: string,
	lastName: string,
	password: string,
	passwordConfirm: string
}

interface RegisterFormProps extends React.ComponentPropsWithoutRef<"form"> {
	toggleMode: () => void;
}

export function RegisterForm({
	className,
	toggleMode,
	...props
}: RegisterFormProps) {
	const [formData, setFormData] = useState<FormData>({
		email: '',
		firstName: '',
		lastName: '',
		password: '',
		passwordConfirm: ''
	});

	const [errors, setErrors] = useState<Partial<FormData>>({});

	const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
		const { name, value } = e.target;
		setFormData({
			...formData,
			[name]: value,
		});
	};

	const validateForm = () => {
		const newErrors: Partial<FormData> = {};

		if (!formData.email) {
			newErrors.email = "Email is required";
		} else if (!/\S+@\S+\.\S+/.test(formData.email)) {
			newErrors.email = "Email address is invalid";
		}

		if (!formData.firstName) {
			newErrors.firstName = "First name is required";
		}

		if (!formData.lastName) {
			newErrors.lastName = "Last name is required";
		}

		if (!formData.password) {
			newErrors.password = "Password is required";
		} else if (formData.password.length < 8) {
			newErrors.password = "Password must be at least 8 characters long";
		}

		if (!formData.passwordConfirm) {
			newErrors.passwordConfirm = "Please confirm your password";
		} else if (formData.passwordConfirm !== formData.password) {
			newErrors.passwordConfirm = "Passwords do not match";
		}

		setErrors(newErrors);
		return Object.keys(newErrors).length === 0;
	};

	const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		const isValid = validateForm();
		if (isValid) {
			// Proceed with form submission
			console.log("Form data submitted:", formData);
			api.post('api/auth/login', { formData }).then(res => {
				console.log(res)
				console.log(res.data)
			})
		}
	};

	return (
		<form onSubmit={handleSubmit} className={cn("flex flex-col gap-6", className)} {...props}>
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
						name="firstName"
						placeholder="Jason"
						onChange={handleChange}
						required
					/>
					{errors.firstName && <span className="text-sm text-red-500">{errors.firstName}</span>}
				</div>
				<div className="grid gap-2">
					<Label htmlFor="lastName">Last Name</Label>
					<Input
						autoComplete="family-name"
						id="lastName"
						type="text"
						name="lastName"
						placeholder="Millman"
						onChange={handleChange}
						required
					/>
					{errors.lastName && <span className="text-sm text-red-500">{errors.lastName}</span>}
				</div>
				<div className="grid gap-2">
					<Label htmlFor="email">Email</Label>
					<Input
						autoComplete="email"
						id="email"
						type="email"
						name="email"
						placeholder="jasonrmillman@gmail.com"
						onChange={handleChange}
						required
					/>
					{errors.email && <span className="text-sm text-red-500">{errors.email}</span>}
				</div>
				<div className="grid gap-2">
					<div className="flex items-center">
						<Label htmlFor="password">Password</Label>
					</div>
					<Input id="password" name="password" type="password" onChange={handleChange} required />
					{errors.password && <span className="text-sm text-red-500">{errors.password}</span>}
				</div>
				<div className="grid gap-2">
					<div className="flex items-center">
						<Label htmlFor="confirmPassword">
							Confirm Password
						</Label>
					</div>
					<Input id="confirmPassword" name="passwordConfirm" type="password" onChange={handleChange} required />
					{errors.passwordConfirm && <span className="text-sm text-red-500">{errors.passwordConfirm}</span>}
				</div>
				<Button type="submit" className="w-full">
					Signup
				</Button>
				<div className="relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t after:border-border">
					<span className="relative z-10 bg-background px-2 text-muted-foreground">
						Or continue with
					</span>
				</div>
				<Button type="button" variant="outline" className="w-full">
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
