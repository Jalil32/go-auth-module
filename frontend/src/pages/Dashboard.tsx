import { useEffect, useState } from "react";
import GenericPageTemplate from "./GenericPageTemplate";

const API_BASE_URL = import.meta.env.DEV
	? import.meta.env.VITE_API_DEV
	: import.meta.env.VITE_API_FLY;

export default function Dashboard() {
	const [backendData, setBackendData] = useState(null);
	const [error, setError] = useState(null);

	console.log(API_BASE_URL); // remove after use

	// Temporary for ya'll to see backend response - remove when updating this component
	useEffect(() => {
		fetch(`${API_BASE_URL}/api/test`)
			.then((response) => {
				if (!response.ok) {
					throw new Error(`Network error: ${response.statusText}`);
				}
				return response.json();
			})
			.then((data) => setBackendData(data))
			.catch((err) => setError(err.message));
	}, []);

	const dashboardContent = (
		<div className="flex flex-1 flex-col gap-4 p-4 pt-0">
			<div className="grid auto-rows-min gap-4 md:grid-cols-3">
				<div className="aspect-video rounded-xl bg-muted/50" />
				<div className="aspect-video rounded-xl bg-muted/50" />
				<div className="aspect-video rounded-xl bg-muted/50" />
			</div>
			<div className="min-h-[100vh] flex-1 rounded-xl bg-muted/50 md:min-h-min p-4">
				{error ? (
					<div>Error: {error}</div>
				) : backendData ? ( // shouldnt chain these conditionals but its temporary so meh
					<div>
						<h2 className="text-xl font-bold mb-2">
							Backend Response:
						</h2>
						<pre className="bg-black p-2 rounded">
							{JSON.stringify(backendData, null, 2)}
						</pre>
					</div>
				) : (
					<div>Loading backend data...</div>
				)}
			</div>
		</div>
	);

	return <GenericPageTemplate pageContent={dashboardContent} />;
}
