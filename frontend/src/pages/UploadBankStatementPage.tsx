import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Table } from "@/components/ui/table";
import Papa from "papaparse";
import { useState } from "react";
import GenericPageTemplate from "./GenericPageTemplate";

const UploadBankStatementPage = () => {
	const [statementData, setStatementData] = useState<string[]>([]);

	const uploadHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
		const file = event.target.files?.[0];
		if (!file) return; // TODO: Add error notification here

		Papa.parse(file, {
			header: false,
			skipEmptyLines: true,
			complete: (results) => {
				// TODO: Add success notification here
				setStatementData(results.data as string[]);
			},
		});
	};

	const bankStatementUpload = (
		<div className="flex flex-1 p-4 pt-0">
			<div className="flex flex-col md:flex-row flex-1 rounded-xl bg-muted/50 space-y-4 md:space-x-4">
				<div className="flex-1">
					<Label htmlFor="csv">Upload your bank statement</Label>
					<Input
						id="csv"
						type="file"
						accept=".csv"
						onChange={uploadHandler}
					/>
				</div>
				<div className="flex items-center justify-center">
					<Separator
						className="hidden md:block h-[90%] border-muted/50"
						orientation="vertical"
					/>
					<Separator
						className="block md:hidden w-[90%] border-muted/50"
						orientation="horizontal"
					/>
				</div>
				<div className="flex-1">
					<Table></Table>
				</div>
			</div>
		</div>
	);

	return <GenericPageTemplate pageContent={bankStatementUpload} />;
};

export default UploadBankStatementPage;
