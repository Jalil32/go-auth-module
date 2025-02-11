import { Button } from "@/components/ui/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuLabel,
	DropdownMenuRadioGroup,
	DropdownMenuRadioItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
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
import { ChevronDown, Eye, Upload } from "lucide-react";
import Papa from "papaparse";
import { useState } from "react";
import GenericPageTemplate from "./GenericPageTemplate";

// Supported file types for bank statement upload - Add more as needed in the future
const SUPPORTED_FILE_TYPES = ["text/csv", "application/vnd.ms-excel"];

// Number of rows to preview from the uploaded file - Too many rows can cause CSS issues
const NUM_ROWS_TO_PREVIEW = 5;

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
	const [file, setFile] = useState<File | null>(null);
	const [statementData, setStatementData] = useState<string[][]>([]);
	const [previewData, setPreviewData] = useState<string[][]>([]);
	const [dropdownValue, setDropdownValue] = useState<string[]>([]);

	const selectFileHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
		const file = event.target.files?.[0];
		if (!file || !SUPPORTED_FILE_TYPES.includes(file.type)) {
			// TODO: Add error notification here: File type not supported
			return;
		}

		setFile(file);
	};

	const previewFileHandler = () => {
		if (file) {
			Papa.parse(file, {
				header: false,
				skipEmptyLines: true,
				complete: (results) => {
					setStatementData(results.data as string[][]);
					const previewData = results.data.slice(
						0,
						NUM_ROWS_TO_PREVIEW,
					);
					setPreviewData(previewData as string[][]);
				},
				error: (error) => {
					// TODO: Add error notification here: Error parsing file
				},
			});
		} else {
			// TODO: Add error notification here: No file selected
		}
	};

	const changeHeaderHandler = (index: number, value: string) => {
		const newDropdownValue = [...dropdownValue];
		const existingIndex = newDropdownValue.indexOf(value); // Check if the value is already selected

		// If the value is already selected, remove it from the previous index
		if (existingIndex !== -1) {
			newDropdownValue[existingIndex] = "";
		}

		newDropdownValue[index] = value;
		setDropdownValue(newDropdownValue);
	};

	const uploadFileHandler = () => {};

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
						<div className="flex flex-row space-x-2">
							<Input
								id="csv"
								type="file"
								accept=".csv"
								className="max-w-lg cursor-pointer"
								onChange={selectFileHandler}
							/>
							<Button
								variant="secondary"
								size="default"
								onClick={() => {
									previewFileHandler();
								}}
							>
								Preview File
							</Button>
						</div>
						<p className="text-xs">Supported file types: .csv</p>
					</div>
				</div>
				<div className="flex items-center justify-center">
					<Separator
						className="w-[95%] border-muted/50"
						orientation="horizontal"
					/>
				</div>
				<div className="flex-1 m-5 max-w-[80vw]">
					{statementData.length > 0 ? (
						<div className="flex flex-col items-center space-y-5">
							<p className="text-sm">
								Preview of the uploaded bank statement is shown
								below. Allocate the headers:{" "}
								<span className="font-semibold italic">
									Date, Description, Amount
								</span>{" "}
								to the respective columns to upload your
								transaction data.
							</p>
							<Table>
								<TableHeader>
									<TableRow>
										{statementData[0].map((_, index) => {
											return (
												<TableHead key={index}>
													<DropdownMenu>
														<DropdownMenuTrigger
															asChild
														>
															<Button
																variant="ghost"
																className="w-20"
															>
																{dropdownValue[
																	index
																] || "Options"}
																<ChevronDown />
															</Button>
														</DropdownMenuTrigger>
														<DropdownMenuContent className="w-56">
															<DropdownMenuLabel>
																Select Header
															</DropdownMenuLabel>
															<DropdownMenuSeparator />
															<DropdownMenuRadioGroup
																value={
																	dropdownValue[
																		index
																	]
																}
																onValueChange={(
																	value: string,
																) => {
																	changeHeaderHandler(
																		index,
																		value,
																	);
																}}
															>
																{Object.values(
																	HEADER_OPTIONS,
																).map(
																	(
																		headerOption,
																		headerOptionIndex,
																	) => (
																		<DropdownMenuRadioItem
																			key={
																				headerOptionIndex
																			}
																			value={
																				headerOption
																			}
																		>
																			{
																				headerOption
																			}
																		</DropdownMenuRadioItem>
																	),
																)}
															</DropdownMenuRadioGroup>
														</DropdownMenuContent>
													</DropdownMenu>
												</TableHead>
											);
										})}
									</TableRow>
								</TableHeader>
								<TableBody>
									{previewData.map((row, rowIndex) => (
										<TableRow key={rowIndex}>
											{row.map((cell, cellIndex) => (
												<TableCell
													key={cellIndex}
													className="p-0 md:p-3"
												>
													{cell}
												</TableCell>
											))}
										</TableRow>
									))}
								</TableBody>
								<TableCaption>
									Preview of the first {NUM_ROWS_TO_PREVIEW}{" "}
									rows
								</TableCaption>
							</Table>
							<Button
								variant="default"
								size="default"
								onClick={() => {
									uploadFileHandler();
								}}
							>
								Save & Upload File
							</Button>
						</div>
					) : (
						<div className="flex flex-col justify-center items-center h-full text-muted space-y-10">
							<div className="flex flex-col justify-center items-center space-y-2">
								<h1 className="flex flex-row items-center text-2xl font-semibold">
									<Eye className="mr-2" /> File Preview
								</h1>
								<p>Preview of uploaded file will appear here</p>
							</div>
							<div>
								<h1 className="text-center font-semibold">
									Your file should look like this...
								</h1>
								<Table>
									<TableCaption className="text-muted">
										Sample Statement Data
									</TableCaption>
									<TableBody>
										{SAMPLE_STATEMENT_DATA.map(
											(statement, index) => (
												<TableRow key={index}>
													<TableCell className="p-0 md:p-3">
														{statement.DATE}
													</TableCell>
													<TableCell className="p-0 md:p-3">
														{statement.AMOUNT}
													</TableCell>
													<TableCell className="p-0 md:p-3">
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
