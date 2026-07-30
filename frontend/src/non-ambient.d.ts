declare global {
	type ToastMessage = {
		id: string;
		message: string;
		type: 'success' | 'error' | 'warning' | 'info';
		duration?: number;
	};
	interface CountryItem {
		name: string;
		code: string;
		emoji: string;
		unicode: string;
		image: string;
		dial_code: string;
		minLength: number;
		maxLength: number;
		regexPattern: string; // Regex pattern for phone number validation
	}
	interface TimezoneOption {
		zone: string;
		gmt: string;
		name: string;
	}
	interface DateRangeChartLabel {
		dates: Date[];
		labels: string[];
	}
	interface MenuItem {
		id: number;
		title: string;
		url: string;
		icon?: Component<SVGAttributes<SVGSVGElement>, {}, string>;
		child?: {
			title: string;
			url: string;
			icon?: Component<SVGAttributes<SVGSVGElement>, {}, string>;
		}[];
	}
	interface BreadcrumbItem {
		name?: string;
		title: string;
		href: string;
	}
	interface NavItem {
		id: number;
		title: string;
		href: string;
		icon?: Component<SVGAttributes<SVGSVGElement>, {}, string>;
		isActive?: boolean;
	}
	interface NavMenu {
		id: number;
		title: string;
		url: string;
		icon: Component<SVGAttributes<SVGSVGElement>, {}, string>;
		isActive?: boolean;
		navItem: NavItem[];
	}
	type BreadcrumbItemType = BreadcrumbItem;

	type StatCard = {
		title: string;
		description?: string;
		icon: Component<LucideProps, {}, string>;
		value?: number;
		variation?: number;
		bgColor?: string;
		borderColor?: string;
		iconColor?: string;
		textColor?: string;
	};
	interface TimeComponents {
		hours: number;
		minutes: number;
		seconds: number;
	}
}

export {};
