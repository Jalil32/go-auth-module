import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import GenericPageTemplate from "./GenericPageTemplate";

const UploadBankStatementPage = () => {
	const bankStatementUpload = (
		<div className="flex flex-1 flex-col p-4 pt-0">
			<div className="flex-1 rounded-xl bg-muted/50">
				<div className="">
					<Label htmlFor="picture">Picture</Label>
					<Input id="picture" type="file" />
				</div>
			</div>
		</div>
	);

	return <GenericPageTemplate pageContent={bankStatementUpload} />;
};

export default UploadBankStatementPage;
