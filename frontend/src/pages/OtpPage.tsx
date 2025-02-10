import {
	InputOTP,
	InputOTPGroup,
	InputOTPSeparator,
	InputOTPSlot,
} from "@/components/ui/input-otp";
import { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

export function InputOtpPage() {
	const location = useLocation();
	const navigate = useNavigate();

	const email = location.state?.email || "";

	if (!location.state || email === "") {
		navigate("/error");
		return null;
	}

	return (
		<div className="w-full h-screen overflow-y-hidden flex items-center justify-center flex-col space-y-6">
			<p className="text-balance text-sm text-muted-foreground">
				Please enter the one time password send to {email}
			</p>
			<InputOTP maxLength={6}>
				<InputOTPGroup>
					<InputOTPSlot index={0} />
					<InputOTPSlot index={1} />
					<InputOTPSlot index={2} />
				</InputOTPGroup>
				<InputOTPSeparator />
				<InputOTPGroup>
					<InputOTPSlot index={3} />
					<InputOTPSlot index={4} />
					<InputOTPSlot index={5} />
				</InputOTPGroup>
			</InputOTP>
		</div>
	);
}
