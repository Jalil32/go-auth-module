import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import {
	Table,
	TableBody,
	TableCaption,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuRadioGroup,
	DropdownMenuRadioItem,
	DropdownMenuTrigger,
} from "@radix-ui/react-dropdown-menu";
import { ChevronDown, Eye, Upload } from "lucide-react";
import Papa from "papaparse";
import { useState } from "react";
import GenericPageTemplate from "./GenericPageTemplate";

// Supported file types for bank statement upload - Add more as needed in the future
const SUPPORTED_FILE_TYPES = ["text/csv", "application/vnd.ms-excel"];

// Number of rows to preview from the uploaded file - Too many rows can cause CSS issues
const NUM_ROWS_TO_PREVIEW = 10;

// Header options for the bank statement upload table
const HEADER_OPTIONS = {
	DATE: "Date",
	DESCRIPTION: "Description",
	AMOUNT: "Amount",
};

// Sample statement data to show users what the uploaded file should look like
const SAMPLE_STATEMENT_DATA = [
	{
		DATE: "30/01/2025",
		AMOUNT: -39.93,
		DESCRIPTION: "Purchase at Starbucks",
	},
	{
		DATE: "1/02/2025",
		AMOUNT: 100,
		DESCRIPTION: "Freelance Project Payment",
	},
	{
		DATE: "3/02/2025",
		AMOUNT: -22.75,
		DESCRIPTION: "Grocery Store Purchase",
	},
	{ DATE: "5/02/2025", AMOUNT: -19.99, DESCRIPTION: "Netflix Subscription" },
	{ DATE: "7/02/2025", AMOUNT: -15.4, DESCRIPTION: "Fast Food Order" },
];

const UploadBankStatementPage = () => {
	const [statementData, setStatementData] = useState<string[][]>([]);
	const [selectedHeaders, setSelectedHeaders] = useState<{
		[key: number]: string;
	}>({});

	const uploadHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
		const file = event.target.files?.[0];
		if (!file || !SUPPORTED_FILE_TYPES.includes(file.type)) {
			// TODO: Add error notification here
			return;
		}

		Papa.parse(file, {
			header: false,
			skipEmptyLines: true,
			preview: NUM_ROWS_TO_PREVIEW,
			complete: (results) => {
				setStatementData(results.data as string[][]);
			},
		});
	};

	const bankStatementUpload = (
		<div className="flex flex-1 p-4 pt-0">
			<div className="flex flex-col flex-1 rounded-xl bg-muted/50 space-y-4">
				<div className="flex flex-col space-y-10 m-5">
					<div className="space-y-2">
						<h1 className="flex flex-row items-center text-2xl font-semibold">
							<Upload className="mr-2" /> Upload File
						</h1>
						<p>
							Upload your bank statement here to record your
							transactions
						</p>
					</div>
					<div className="flex flex-col items-center justify-center space-y-2">
						<Input
							id="csv"
							type="file"
							accept=".csv"
							onChange={uploadHandler}
						/>
						<p className="text-xs">Supported file types: .csv</p>
					</div>
				</div>
				<div className="flex items-center justify-center">
					<Separator
						className="w-[95%] border-muted/50"
						orientation="horizontal"
					/>
				</div>
				<div className="flex-1 m-5">
					{statementData.length > 0 ? (
						<Table>
							<TableHeader>
								<TableRow>
									{statementData[0].map((_, index) => (
										<TableHead key={index}>
											<DropdownMenu>
												<DropdownMenuTrigger asChild>
													<Button variant="ghost">
														<ChevronDown />
													</Button>
												</DropdownMenuTrigger>
												<DropdownMenuContent>
													<DropdownMenuRadioGroup>
														{Object.values(
															HEADER_OPTIONS,
														).map(
															(
																option: string,
															) => (
																<DropdownMenuRadioItem
																	key={option}
																	value={
																		option
																	}
																>
																	<DropdownMenuItem
																		onSelect={() =>
																			setSelectedHeaders(
																				{
																					...selectedHeaders,
																					[index]:
																						option,
																				},
																			)
																		}
																	>
																		{option}
																	</DropdownMenuItem>
																</DropdownMenuRadioItem>
															),
														)}
													</DropdownMenuRadioGroup>
												</DropdownMenuContent>
											</DropdownMenu>
										</TableHead>
									))}
								</TableRow>
							</TableHeader>
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
							<TableCaption>
								Preview of the first {NUM_ROWS_TO_PREVIEW} rows
							</TableCaption>
						</Table>
					) : (
						<div className="flex flex-col justify-center items-center h-full text-muted space-y-4">
							<div className="flex flex-col justify-center items-center space-y-2">
								<h1 className="flex flex-row items-center text-2xl font-semibold">
									<Eye className="mr-2" /> File Preview
								</h1>
								<p>Preview of uploaded file will appear here</p>
							</div>
							<div>
								<Table>
									<TableCaption className="text-muted">
										Sample Statement Data
									</TableCaption>
									<TableHeader>
										<TableRow>
											<TableHead className="text-muted">
												{HEADER_OPTIONS.DATE}
											</TableHead>
											<TableHead className="text-muted">
												{HEADER_OPTIONS.AMOUNT}
											</TableHead>
											<TableHead className="text-muted">
												{HEADER_OPTIONS.DESCRIPTION}
											</TableHead>
										</TableRow>
									</TableHeader>
									<TableBody>
										{SAMPLE_STATEMENT_DATA.map(
											(statement) => (
												<TableRow key={statement.DATE}>
													<TableCell>
														{statement.DATE}
													</TableCell>
													<TableCell>
														{statement.AMOUNT}
													</TableCell>
													<TableCell>
														{statement.DESCRIPTION}
													</TableCell>
												</TableRow>
											),
										)}
									</TableBody>
								</Table>
							</div>
						</div>
					)}
				</div>
			</div>
		</div>
	);

	return <GenericPageTemplate pageContent={bankStatementUpload} />;
};

export default UploadBankStatementPage;
