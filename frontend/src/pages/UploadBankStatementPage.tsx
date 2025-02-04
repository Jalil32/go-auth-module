import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";
import { Upload } from "lucide-react";
import Papa from "papaparse";
import { useState } from "react";
import GenericPageTemplate from "./GenericPageTemplate";

const UploadBankStatementPage = () => {
	const [statementData, setStatementData] = useState<string[][]>([]);

	const uploadHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
		const file = event.target.files?.[0];
		if (!file) return; // TODO: Add error notification here

		Papa.parse(file, {
			header: false,
			skipEmptyLines: true,
			complete: (results) => {
				// TODO: Add success notification here
				setStatementData(results.data as string[][]);
			},
		});
	};

	const bankStatementUpload = (
		<div className="flex flex-1 p-4 pt-0">
			<div className="flex flex-col md:flex-row flex-1 rounded-xl bg-muted/50 space-y-4 md:space-x-4">
				<div className="flex flex-col flex-1 space-y-10 m-5">
					<div className="space-y-2">
						<h1 className="flex flex-row items-center text-2xl font-semibold">
							<Upload className="mr-2" /> Upload File
						</h1>
						<p>
							Upload your bank statement here to record your
							transactions
						</p>
					</div>
					<div>
						<Input
							id="csv"
							type="file"
							accept=".csv"
							onChange={uploadHandler}
						/>
						<p className="text-xs text-center">
							Supported file types: .csv
						</p>
					</div>
				</div>
				<div className="flex items-center justify-center">
					<Separator
						className="hidden md:block h-[95%] border-muted/50"
						orientation="vertical"
					/>
					<Separator
						className="block md:hidden w-[95%] border-muted/50"
						orientation="horizontal"
					/>
				</div>
				<div className="flex-1 m-5 overflow-auto">
					<Table>
						<TableBody>
							{statementData.map(
								(row: string[], rowIndex: number) => (
									<TableRow key={rowIndex}>
										{row.map(
											(
												cell: string,
												cellIndex: number,
											) => (
												<TableCell key={cellIndex}>
													{cell}
												</TableCell>
											),
										)}
									</TableRow>
								),
							)}
						</TableBody>
					</Table>
				</div>
			</div>
		</div>
	);

	return <GenericPageTemplate pageContent={bankStatementUpload} />;
};

export default UploadBankStatementPage;
