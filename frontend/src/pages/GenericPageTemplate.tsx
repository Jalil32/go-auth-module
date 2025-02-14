import { AppSidebar } from "@/components/app-sidebar";
import {
	Breadcrumb,
	BreadcrumbItem,
	BreadcrumbLink,
	BreadcrumbList,
	BreadcrumbPage,
	BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import {
	SidebarInset,
	SidebarProvider,
	SidebarTrigger,
} from "@/components/ui/sidebar";
import { Separator } from "@radix-ui/react-separator";
import React from "react";

/*
 * Component used to create pages with consistent layout.
 * When creating a new page, use this component as the foundation.
 */

interface GenericPageTemplateProps {
	pageContent: React.ReactNode;
}

const GenericPageTemplate = ({ pageContent }: GenericPageTemplateProps) => {
	const currUrl = window.location.href;
	const pathSegments: string[] = currUrl.split("/");

	// Dynamically generate breadcrumb links based on the current URL
	const getBreadCrumbLinks = (items: string[]) => {
		const domain = items.slice(0, 3).join("/"); // e.g. https://wealthscope.com
		const breadCrumbNames: string[] = items.slice(3, items.length).map(
			(c) =>
				c
					.replace(/-/g, " ")
					.replace(/\b\w/g, (char) => char.toUpperCase()), // Don't ask me how this regex works, it was Copilot
		); // e.g. "https://wealthscope.com/dashboard/assets" -> ["Dashboard", "Assets"]
		const breadCrumbLinks: { [key: string]: string } = {};

		breadCrumbNames.forEach((name, index) => {
			const path = items.slice(3, 4 + index).join("/");
			breadCrumbLinks[name] = `${domain}/${path}`;
		});

		// e.g. {dashboard: 'http://localhost:5173/dashboard', assets: 'http://localhost:5173/dashboard/assets'}
		return breadCrumbLinks;
	};

	const breadCrumbLinks = getBreadCrumbLinks(pathSegments);

	return (
		<SidebarProvider>
			<AppSidebar />
			<SidebarInset>
				<header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
					<div className="flex items-center gap-2 px-4">
						<SidebarTrigger className="-ml-1" />
						<Separator
							orientation="vertical"
							className="mr-2 h-4"
						/>
						<Breadcrumb>
							<BreadcrumbList>
								{Object.entries(breadCrumbLinks).map(
									([name, link], index) => (
										<React.Fragment
											key={`breadcrumb-${name}=${link}`}
										>
											{index !==
											Object.keys(breadCrumbLinks)
												.length -
												1 ? (
												// If not the last breadcrumb item, render as a link with separator
												<>
													<BreadcrumbItem className="hidden md:block">
														<BreadcrumbLink
															href={link}
														>
															{name}
														</BreadcrumbLink>
													</BreadcrumbItem>
													<BreadcrumbSeparator className="hidden md:block" />
												</>
											) : (
												// If the last breadcrumb item, render as a page
												<BreadcrumbItem>
													<BreadcrumbPage>
														{name}
													</BreadcrumbPage>
												</BreadcrumbItem>
											)}
										</React.Fragment>
									),
								)}
							</BreadcrumbList>
						</Breadcrumb>
					</div>
				</header>
				{pageContent}
			</SidebarInset>
		</SidebarProvider>
	);
};

export default GenericPageTemplate;
