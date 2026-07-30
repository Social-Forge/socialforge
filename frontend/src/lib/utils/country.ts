export const CountryList: CountryItem[] = [
	{
		name: 'Andorra',
		code: 'AD',
		emoji: '🇦🇩',
		unicode: 'U+1F1E6 U+1F1E9',
		image: 'https://country-code-au6g.vercel.app/AD.svg',
		dial_code: '+376',
		minLength: 6, // Umumnya 6 hingga 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{6,9})$'
	},
	{
		name: 'United Arab Emirates',
		code: 'AE',
		emoji: '🇦🇪',
		unicode: 'U+1F1E6 U+1F1EA',
		image: 'https://country-code-au6g.vercel.app/AE.svg',
		dial_code: '+971',
		minLength: 8, // Umumnya 8 hingga 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Afghanistan',
		code: 'AF',
		emoji: '🇦🇫',
		unicode: 'U+1F1E6 U+1F1EB',
		image: 'https://country-code-au6g.vercel.app/AF.svg',
		dial_code: '+93',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Antigua & Barbuda',
		code: 'AG',
		emoji: '🇦🇬',
		unicode: 'U+1F1E6 U+1F1EC',
		image: 'https://country-code-au6g.vercel.app/AG.svg',
		dial_code: '+1268',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Anguilla',
		code: 'AI',
		emoji: '🇦🇮',
		unicode: 'U+1F1E6 U+1F1EE',
		image: 'https://country-code-au6g.vercel.app/AI.svg',
		dial_code: '+1264',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Albania',
		code: 'AL',
		emoji: '🇦🇱',
		unicode: 'U+1F1E6 U+1F1F1',
		image: 'https://country-code-au6g.vercel.app/AL.svg',
		dial_code: '+355',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Armenia',
		code: 'AM',
		emoji: '🇦🇲',
		unicode: 'U+1F1E6 U+1F1F2',
		image: 'https://country-code-au6g.vercel.app/AM.svg',
		dial_code: '+374',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Angola',
		code: 'AO',
		emoji: '🇦🇴',
		unicode: 'U+1F1E6 U+1F1F4',
		image: 'https://country-code-au6g.vercel.app/AO.svg',
		dial_code: '+244',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Antarctica',
		code: 'AQ',
		emoji: '🇦🇶',
		unicode: 'U+1F1E6 U+1F1F6',
		image: 'https://country-code-au6g.vercel.app/AQ.svg',
		dial_code: '+672',
		minLength: 6, // Bervariasi, menggunakan 6-9
		maxLength: 9,
		regexPattern: '^(\\d{6,9})$'
	},
	{
		name: 'Argentina',
		code: 'AR',
		emoji: '🇦🇷',
		unicode: 'U+1F1E6 U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/AR.svg',
		dial_code: '+54',
		minLength: 8, // Bervariasi, umumnya 8 atau 10 digit (untuk seluler)
		maxLength: 10,
		regexPattern: '^(\\d{8,10})$'
	},
	{
		name: 'American Samoa',
		code: 'AS',
		emoji: '🇦🇸',
		unicode: 'U+1F1E6 U+1F1F8',
		image: 'https://country-code-au6g.vercel.app/AS.svg',
		dial_code: '+1684',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Austria',
		code: 'AT',
		emoji: '🇦🇹',
		unicode: 'U+1F1E6 U+1F1F9',
		image: 'https://country-code-au6g.vercel.app/AT.svg',
		dial_code: '+43',
		minLength: 4, // Umumnya 4 hingga 13 digit
		maxLength: 13,
		regexPattern: '^(\\d{4,13})$'
	},
	{
		name: 'Australia',
		code: 'AU',
		emoji: '🇦🇺',
		unicode: 'U+1F1E6 U+1F1FA',
		image: 'https://country-code-au6g.vercel.app/AU.svg',
		dial_code: '+61',
		minLength: 8, // Umumnya 8 hingga 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Aruba',
		code: 'AW',
		emoji: '🇦🇼',
		unicode: 'U+1F1E6 U+1F1FC',
		image: 'https://country-code-au6g.vercel.app/AW.svg',
		dial_code: '+297',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Åland Islands',
		code: 'AX',
		emoji: '🇦🇽',
		unicode: 'U+1F1E6 U+1F1FD',
		image: 'https://country-code-au6g.vercel.app/AX.svg',
		dial_code: '+358',
		minLength: 6, // Sama dengan Finlandia (6-12)
		maxLength: 12,
		regexPattern: '^(\\d{6,12})$'
	},
	{
		name: 'Azerbaijan',
		code: 'AZ',
		emoji: '🇦🇿',
		unicode: 'U+1F1E6 U+1F1FF',
		image: 'https://country-code-au6g.vercel.app/AZ.svg',
		dial_code: '+994',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Bosnia & Herzegovina',
		code: 'BA',
		emoji: '🇧🇦',
		unicode: 'U+1F1E7 U+1F1E6',
		image: 'https://country-code-au6g.vercel.app/BA.svg',
		dial_code: '+387',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Barbados',
		code: 'BB',
		emoji: '🇧🇧',
		unicode: 'U+1F1E7 U+1F1E7',
		image: 'https://country-code-au6g.vercel.app/BB.svg',
		dial_code: '+1246',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Bangladesh',
		code: 'BD',
		emoji: '🇧🇩',
		unicode: 'U+1F1E7 U+1F1E9',
		image: 'https://country-code-au6g.vercel.app/BD.svg',
		dial_code: '+880',
		minLength: 10, // Umumnya 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Belgium',
		code: 'BE',
		emoji: '🇧🇪',
		unicode: 'U+1F1E7 U+1F1EA',
		image: 'https://country-code-au6g.vercel.app/BE.svg',
		dial_code: '+32',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Burkina Faso',
		code: 'BF',
		emoji: '🇧🇫',
		unicode: 'U+1F1E7 U+1F1EB',
		image: 'https://country-code-au6g.vercel.app/BF.svg',
		dial_code: '+226',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Bulgaria',
		code: 'BG',
		emoji: '🇧🇬',
		unicode: 'U+1F1E7 U+1F1EC',
		image: 'https://country-code-au6g.vercel.app/BG.svg',
		dial_code: '+359',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Bahrain',
		code: 'BH',
		emoji: '🇧🇭',
		unicode: 'U+1F1E7 U+1F1ED',
		image: 'https://country-code-au6g.vercel.app/BH.svg',
		dial_code: '+973',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Burundi',
		code: 'BI',
		emoji: '🇧🇮',
		unicode: 'U+1F1E7 U+1F1EE',
		image: 'https://country-code-au6g.vercel.app/BI.svg',
		dial_code: '+257',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Benin',
		code: 'BJ',
		emoji: '🇧🇯',
		unicode: 'U+1F1E7 U+1F1EF',
		image: 'https://country-code-au6g.vercel.app/BJ.svg',
		dial_code: '+229',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'St. Barthélemy',
		code: 'BL',
		emoji: '🇧🇱',
		unicode: 'U+1F1E7 U+1F1F1',
		image: 'https://country-code-au6g.vercel.app/BL.svg',
		dial_code: '+590',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Bermuda',
		code: 'BM',
		emoji: '🇧🇲',
		unicode: 'U+1F1E7 U+1F1F2',
		image: 'https://country-code-au6g.vercel.app/BM.svg',
		dial_code: '+1441',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Brunei',
		code: 'BN',
		emoji: '🇧🇳',
		unicode: 'U+1F1E7 U+1F1F3',
		image: 'https://country-code-au6g.vercel.app/BN.svg',
		dial_code: '+673',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Bolivia',
		code: 'BO',
		emoji: '🇧🇴',
		unicode: 'U+1F1E7 U+1F1F4',
		image: 'https://country-code-au6g.vercel.app/BO.svg',
		dial_code: '+591',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Brazil',
		code: 'BR',
		emoji: '🇧🇷',
		unicode: 'U+1F1E7 U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/BR.svg',
		dial_code: '+55',
		minLength: 8, // Umumnya 8 sampai 11 digit
		maxLength: 11,
		regexPattern: '^(\\d{8,11})$'
	},
	{
		name: 'Bahamas',
		code: 'BS',
		emoji: '🇧🇸',
		unicode: 'U+1F1E7 U+1F1F8',
		image: 'https://country-code-au6g.vercel.app/BS.svg',
		dial_code: '+1242',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Bhutan',
		code: 'BT',
		emoji: '🇧🇹',
		unicode: 'U+1F1E7 U+1F1F9',
		image: 'https://country-code-au6g.vercel.app/BT.svg',
		dial_code: '+975',
		minLength: 7, // Umumnya 7 atau 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{7,8})$'
	},
	{
		name: 'Botswana',
		code: 'BW',
		emoji: '🇧🇼',
		unicode: 'U+1F1E7 U+1F1FC',
		image: 'https://country-code-au6g.vercel.app/BW.svg',
		dial_code: '+267',
		minLength: 7, // Umumnya 7 atau 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{7,8})$'
	},
	{
		name: 'Belarus',
		code: 'BY',
		emoji: '🇧🇾',
		unicode: 'U+1F1E7 U+1F1FE',
		image: 'https://country-code-au6g.vercel.app/BY.svg',
		dial_code: '+375',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Belize',
		code: 'BZ',
		emoji: '🇧🇿',
		unicode: 'U+1F1E7 U+1F1FF',
		image: 'https://country-code-au6g.vercel.app/BZ.svg',
		dial_code: '+501',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Canada',
		code: 'CA',
		emoji: '🇨🇦',
		unicode: 'U+1F1E8 U+1F1E6',
		image: 'https://country-code-au6g.vercel.app/CA.svg',
		dial_code: '+1',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Cocos (Keeling) Islands',
		code: 'CC',
		emoji: '🇨🇨',
		unicode: 'U+1F1E8 U+1F1E8',
		image: 'https://country-code-au6g.vercel.app/CC.svg',
		dial_code: '+61',
		minLength: 8, // Sama dengan Australia (8-9)
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Congo - Kinshasa',
		code: 'CD',
		emoji: '🇨🇩',
		unicode: 'U+1F1E8 U+1F1E9',
		image: 'https://country-code-au6g.vercel.app/CD.svg',
		dial_code: '+243',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Central African Republic',
		code: 'CF',
		emoji: '🇨🇫',
		unicode: 'U+1F1E8 U+1F1EB',
		image: 'https://country-code-au6g.vercel.app/CF.svg',
		dial_code: '+236',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Congo - Brazzaville',
		code: 'CG',
		emoji: '🇨🇬',
		unicode: 'U+1F1E8 U+1F1EC',
		image: 'https://country-code-au6g.vercel.app/CG.svg',
		dial_code: '+242',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Switzerland',
		code: 'CH',
		emoji: '🇨🇭',
		unicode: 'U+1F1E8 U+1F1ED',
		image: 'https://country-code-au6g.vercel.app/CH.svg',
		dial_code: '+41',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Côte d’Ivoire',
		code: 'CI',
		emoji: '🇨🇮',
		unicode: 'U+1F1E8 U+1F1EE',
		image: 'https://country-code-au6g.vercel.app/CI.svg',
		dial_code: '+225',
		minLength: 8, // Umumnya 8 hingga 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{8,10})$'
	},
	{
		name: 'Cook Islands',
		code: 'CK',
		emoji: '🇨🇰',
		unicode: 'U+1F1E8 U+1F1F0',
		image: 'https://country-code-au6g.vercel.app/CK.svg',
		dial_code: '+682',
		minLength: 5, // Umumnya 5 digit
		maxLength: 5,
		regexPattern: '^(\\d{5})$'
	},
	{
		name: 'Chile',
		code: 'CL',
		emoji: '🇨🇱',
		unicode: 'U+1F1E8 U+1F1F1',
		image: 'https://country-code-au6g.vercel.app/CL.svg',
		dial_code: '+56',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Cameroon',
		code: 'CM',
		emoji: '🇨🇲',
		unicode: 'U+1F1E8 U+1F1F2',
		image: 'https://country-code-au6g.vercel.app/CM.svg',
		dial_code: '+237',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'China',
		code: 'CN',
		emoji: '🇨🇳',
		unicode: 'U+1F1E8 U+1F1F3',
		image: 'https://country-code-au6g.vercel.app/CN.svg',
		dial_code: '+86',
		minLength: 10, // Umumnya 10 atau 11 digit
		maxLength: 11,
		regexPattern: '^(\\d{10,11})$'
	},
	{
		name: 'Colombia',
		code: 'CO',
		emoji: '🇨🇴',
		unicode: 'U+1F1E8 U+1F1F4',
		image: 'https://country-code-au6g.vercel.app/CO.svg',
		dial_code: '+57',
		minLength: 8, // Umumnya 8 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{8,10})$'
	},
	{
		name: 'Costa Rica',
		code: 'CR',
		emoji: '🇨🇷',
		unicode: 'U+1F1E8 U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/CR.svg',
		dial_code: '+506',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Cuba',
		code: 'CU',
		emoji: '🇨🇺',
		unicode: 'U+1F1E8 U+1F1FA',
		image: 'https://country-code-au6g.vercel.app/CU.svg',
		dial_code: '+53',
		minLength: 6, // Umumnya 6 sampai 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{6,8})$'
	},
	{
		name: 'Cape Verde',
		code: 'CV',
		emoji: '🇨🇻',
		unicode: 'U+1F1E8 U+1F1FB',
		image: 'https://country-code-au6g.vercel.app/CV.svg',
		dial_code: '+238',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Christmas Island',
		code: 'CX',
		emoji: '🇨🇽',
		unicode: 'U+1F1E8 U+1F1FD',
		image: 'https://country-code-au6g.vercel.app/CX.svg',
		dial_code: '+61',
		minLength: 8, // Sama dengan Australia (8-9)
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Cyprus',
		code: 'CY',
		emoji: '🇨🇾',
		unicode: 'U+1F1E8 U+1F1FE',
		image: 'https://country-code-au6g.vercel.app/CY.svg',
		dial_code: '+357',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Czechia',
		code: 'CZ',
		emoji: '🇨🇿',
		unicode: 'U+1F1E8 U+1F1FF',
		image: 'https://country-code-au6g.vercel.app/CZ.svg',
		dial_code: '+420',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Germany',
		code: 'DE',
		emoji: '🇩🇪',
		unicode: 'U+1F1E9 U+1F1EA',
		image: 'https://country-code-au6g.vercel.app/DE.svg',
		dial_code: '+49',
		minLength: 4, // Bervariasi, umumnya 4 sampai 12 digit
		maxLength: 12,
		regexPattern: '^(\\d{4,12})$'
	},
	{
		name: 'Djibouti',
		code: 'DJ',
		emoji: '🇩🇯',
		unicode: 'U+1F1E9 U+1F1EF',
		image: 'https://country-code-au6g.vercel.app/DJ.svg',
		dial_code: '+253',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Denmark',
		code: 'DK',
		emoji: '🇩🇰',
		unicode: 'U+1F1E9 U+1F1F0',
		image: 'https://country-code-au6g.vercel.app/DK.svg',
		dial_code: '+45',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Dominica',
		code: 'DM',
		emoji: '🇩🇲',
		unicode: 'U+1F1E9 U+1F1F2',
		image: 'https://country-code-au6g.vercel.app/DM.svg',
		dial_code: '+1767',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Dominican Republic',
		code: 'DO',
		emoji: '🇩🇴',
		unicode: 'U+1F1E9 U+1F1F4',
		image: 'https://country-code-au6g.vercel.app/DO.svg',
		dial_code: '+1849',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Algeria',
		code: 'DZ',
		emoji: '🇩🇿',
		unicode: 'U+1F1E9 U+1F1FF',
		image: 'https://country-code-au6g.vercel.app/DZ.svg',
		dial_code: '+213',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Ecuador',
		code: 'EC',
		emoji: '🇪🇨',
		unicode: 'U+1F1EA U+1F1E8',
		image: 'https://country-code-au6g.vercel.app/EC.svg',
		dial_code: '+593',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Estonia',
		code: 'EE',
		emoji: '🇪🇪',
		unicode: 'U+1F1EA U+1F1EA',
		image: 'https://country-code-au6g.vercel.app/EE.svg',
		dial_code: '+372',
		minLength: 7, // Umumnya 7 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{7,10})$'
	},
	{
		name: 'Egypt',
		code: 'EG',
		emoji: '🇪🇬',
		unicode: 'U+1F1EA U+1F1EC',
		image: 'https://country-code-au6g.vercel.app/EG.svg',
		dial_code: '+20',
		minLength: 8, // Umumnya 8 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{8,10})$'
	},
	{
		name: 'Eritrea',
		code: 'ER',
		emoji: '🇪🇷',
		unicode: 'U+1F1EA U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/ER.svg',
		dial_code: '+291',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Spain',
		code: 'ES',
		emoji: '🇪🇸',
		unicode: 'U+1F1EA U+1F1F8',
		image: 'https://country-code-au6g.vercel.app/ES.svg',
		dial_code: '+34',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Ethiopia',
		code: 'ET',
		emoji: '🇪🇹',
		unicode: 'U+1F1EA U+1F1F9',
		image: 'https://country-code-au6g.vercel.app/ET.svg',
		dial_code: '+251',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Finland',
		code: 'FI',
		emoji: '🇫🇮',
		unicode: 'U+1F1EB U+1F1EE',
		image: 'https://country-code-au6g.vercel.app/FI.svg',
		dial_code: '+358',
		minLength: 6, // Umumnya 6 sampai 12 digit
		maxLength: 12,
		regexPattern: '^(\\d{6,12})$'
	},
	{
		name: 'Fiji',
		code: 'FJ',
		emoji: '🇫🇯',
		unicode: 'U+1F1EB U+1F1EF',
		image: 'https://country-code-au6g.vercel.app/FJ.svg',
		dial_code: '+679',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Falkland Islands',
		code: 'FK',
		emoji: '🇫🇰',
		unicode: 'U+1F1EB U+1F1F0',
		image: 'https://country-code-au6g.vercel.app/FK.svg',
		dial_code: '+500',
		minLength: 5, // Umumnya 5 digit
		maxLength: 5,
		regexPattern: '^(\\d{5})$'
	},
	{
		name: 'Micronesia',
		code: 'FM',
		emoji: '🇫🇲',
		unicode: 'U+1F1EB U+1F1F2',
		image: 'https://country-code-au6g.vercel.app/FM.svg',
		dial_code: '+691',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Faroe Islands',
		code: 'FO',
		emoji: '🇫🇴',
		unicode: 'U+1F1EB U+1F1F4',
		image: 'https://country-code-au6g.vercel.app/FO.svg',
		dial_code: '+298',
		minLength: 6, // Umumnya 6 digit
		maxLength: 6,
		regexPattern: '^(\\d{6})$'
	},
	{
		name: 'France',
		code: 'FR',
		emoji: '🇫🇷',
		unicode: 'U+1F1EB U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/FR.svg',
		dial_code: '+33',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Gabon',
		code: 'GA',
		emoji: '🇬🇦',
		unicode: 'U+1F1EC U+1F1E6',
		image: 'https://country-code-au6g.vercel.app/GA.svg',
		dial_code: '+241',
		minLength: 7, // Umumnya 7 atau 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{7,8})$'
	},
	{
		name: 'United Kingdom',
		code: 'GB',
		emoji: '🇬🇧',
		unicode: 'U+1F1EC U+1F1E7',
		image: 'https://country-code-au6g.vercel.app/GB.svg',
		dial_code: '+44',
		minLength: 9, // Umumnya 9 atau 10 digit (tanpa 0 awal)
		maxLength: 10,
		regexPattern: '^(\\d{9,10})$'
	},
	{
		name: 'Grenada',
		code: 'GD',
		emoji: '🇬🇩',
		unicode: 'U+1F1EC U+1F1E9',
		image: 'https://country-code-au6g.vercel.app/GD.svg',
		dial_code: '+1473',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Georgia',
		code: 'GE',
		emoji: '🇬🇪',
		unicode: 'U+1F1EC U+1F1EA',
		image: 'https://country-code-au6g.vercel.app/GE.svg',
		dial_code: '+995',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'French Guiana',
		code: 'GF',
		emoji: '🇬🇫',
		unicode: 'U+1F1EC U+1F1EB',
		image: 'https://country-code-au6g.vercel.app/GF.svg',
		dial_code: '+594',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Guernsey',
		code: 'GG',
		emoji: '🇬🇬',
		unicode: 'U+1F1EC U+1F1EC',
		image: 'https://country-code-au6g.vercel.app/GG.svg',
		dial_code: '+44',
		minLength: 9, // Sama dengan UK (9-10)
		maxLength: 10,
		regexPattern: '^(\\d{9,10})$'
	},
	{
		name: 'Ghana',
		code: 'GH',
		emoji: '🇬🇭',
		unicode: 'U+1F1EC U+1F1ED',
		image: 'https://country-code-au6g.vercel.app/GH.svg',
		dial_code: '+233',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Gibraltar',
		code: 'GI',
		emoji: '🇬🇮',
		unicode: 'U+1F1EC U+1F1EE',
		image: 'https://country-code-au6g.vercel.app/GI.svg',
		dial_code: '+350',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Greenland',
		code: 'GL',
		emoji: '🇬🇱',
		unicode: 'U+1F1EC U+1F1F1',
		image: 'https://country-code-au6g.vercel.app/GL.svg',
		dial_code: '+299',
		minLength: 6, // Umumnya 6 digit
		maxLength: 6,
		regexPattern: '^(\\d{6})$'
	},
	{
		name: 'Gambia',
		code: 'GM',
		emoji: '🇬🇲',
		unicode: 'U+1F1EC U+1F1F2',
		image: 'https://country-code-au6g.vercel.app/GM.svg',
		dial_code: '+220',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Guinea',
		code: 'GN',
		emoji: '🇬🇳',
		unicode: 'U+1F1EC U+1F1F3',
		image: 'https://country-code-au6g.vercel.app/GN.svg',
		dial_code: '+224',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Guadeloupe',
		code: 'GP',
		emoji: '🇬🇵',
		unicode: 'U+1F1EC U+1F1F5',
		image: 'https://country-code-au6g.vercel.app/GP.svg',
		dial_code: '+590',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Equatorial Guinea',
		code: 'GQ',
		emoji: '🇬🇶',
		unicode: 'U+1F1EC U+1F1F6',
		image: 'https://country-code-au6g.vercel.app/GQ.svg',
		dial_code: '+240',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Greece',
		code: 'GR',
		emoji: '🇬🇷',
		unicode: 'U+1F1EC U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/GR.svg',
		dial_code: '+30',
		minLength: 10, // Umumnya 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'South Georgia & South Sandwich Islands',
		code: 'GS',
		emoji: '🇬🇸',
		unicode: 'U+1F1EC U+1F1F8',
		image: 'https://country-code-au6g.vercel.app/GS.svg',
		dial_code: '+500',
		minLength: 5, // Umumnya 5 digit
		maxLength: 5,
		regexPattern: '^(\\d{5})$'
	},
	{
		name: 'Guatemala',
		code: 'GT',
		emoji: '🇬🇹',
		unicode: 'U+1F1EC U+1F1F9',
		image: 'https://country-code-au6g.vercel.app/GT.svg',
		dial_code: '+502',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Guam',
		code: 'GU',
		emoji: '🇬🇺',
		unicode: 'U+1F1EC U+1F1FA',
		image: 'https://country-code-au6g.vercel.app/GU.svg',
		dial_code: '+1671',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Guinea-Bissau',
		code: 'GW',
		emoji: '🇬🇼',
		unicode: 'U+1F1EC U+1F1FC',
		image: 'https://country-code-au6g.vercel.app/GW.svg',
		dial_code: '+245',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Guyana',
		code: 'GY',
		emoji: '🇬🇾',
		unicode: 'U+1F1EC U+1F1FE',
		image: 'https://country-code-au6g.vercel.app/GY.svg',
		dial_code: '+595',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Hong Kong SAR China',
		code: 'HK',
		emoji: '🇭🇰',
		unicode: 'U+1F1ED U+1F1F0',
		image: 'https://country-code-au6g.vercel.app/HK.svg',
		dial_code: '+852',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Honduras',
		code: 'HN',
		emoji: '🇭🇳',
		unicode: 'U+1F1ED U+1F1F3',
		image: 'https://country-code-au6g.vercel.app/HN.svg',
		dial_code: '+504',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Croatia',
		code: 'HR',
		emoji: '🇭🇷',
		unicode: 'U+1F1ED U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/HR.svg',
		dial_code: '+385',
		minLength: 8, // Umumnya 8 sampai 11 digit
		maxLength: 11,
		regexPattern: '^(\\d{8,11})$'
	},
	{
		name: 'Haiti',
		code: 'HT',
		emoji: '🇭🇹',
		unicode: 'U+1F1ED U+1F1F9',
		image: 'https://country-code-au6g.vercel.app/HT.svg',
		dial_code: '+509',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Hungary',
		code: 'HU',
		emoji: '🇭🇺',
		unicode: 'U+1F1ED U+1F1FA',
		image: 'https://country-code-au6g.vercel.app/HU.svg',
		dial_code: '+36',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Indonesia',
		code: 'ID',
		emoji: '🇮🇩',
		unicode: 'U+1F1EE U+1F1E9',
		image: 'https://country-code-au6g.vercel.app/ID.svg',
		dial_code: '+62',
		minLength: 9, // Umumnya 9 sampai 12 digit (misalnya: 812xxxxxx)
		maxLength: 12,
		regexPattern: '^(\\d{9,12})$'
	},
	{
		name: 'Ireland',
		code: 'IE',
		emoji: '🇮🇪',
		unicode: 'U+1F1EE U+1F1EA',
		image: 'https://country-code-au6g.vercel.app/IE.svg',
		dial_code: '+353',
		minLength: 7, // Umumnya 7 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{7,10})$'
	},
	{
		name: 'Israel',
		code: 'IL',
		emoji: '🇮🇱',
		unicode: 'U+1F1EE U+1F1F1',
		image: 'https://country-code-au6g.vercel.app/IL.svg',
		dial_code: '+972',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Isle of Man',
		code: 'IM',
		emoji: '🇮🇲',
		unicode: 'U+1F1EE U+1F1F2',
		image: 'https://country-code-au6g.vercel.app/IM.svg',
		dial_code: '+44',
		minLength: 9, // Sama dengan UK (9-10)
		maxLength: 10,
		regexPattern: '^(\\d{9,10})$'
	},
	{
		name: 'India',
		code: 'IN',
		emoji: '🇮🇳',
		unicode: 'U+1F1EE U+1F1F3',
		image: 'https://country-code-au6g.vercel.app/IN.svg',
		dial_code: '+91',
		minLength: 10, // Umumnya 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'British Indian Ocean Territory',
		code: 'IO',
		emoji: '🇮🇴',
		unicode: 'U+1F1EE U+1F1F4',
		image: 'https://country-code-au6g.vercel.app/IO.svg',
		dial_code: '+246',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Iraq',
		code: 'IQ',
		emoji: '🇮🇶',
		unicode: 'U+1F1EE U+1F1F6',
		image: 'https://country-code-au6g.vercel.app/IQ.svg',
		dial_code: '+964',
		minLength: 8, // Umumnya 8 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{8,10})$'
	},
	{
		name: 'Iran',
		code: 'IR',
		emoji: '🇮🇷',
		unicode: 'U+1F1EE U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/IR.svg',
		dial_code: '+98',
		minLength: 10, // Umumnya 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Iceland',
		code: 'IS',
		emoji: '🇮🇸',
		unicode: 'U+1F1EE U+1F1F8',
		image: 'https://country-code-au6g.vercel.app/IS.svg',
		dial_code: '+354',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Italy',
		code: 'IT',
		emoji: '🇮🇹',
		unicode: 'U+1F1EE U+1F1F9',
		image: 'https://country-code-au6g.vercel.app/IT.svg',
		dial_code: '+39',
		minLength: 8, // Umumnya 8 sampai 11 digit
		maxLength: 11,
		regexPattern: '^(\\d{8,11})$'
	},
	{
		name: 'Jersey',
		code: 'JE',
		emoji: '🇯🇪',
		unicode: 'U+1F1EF U+1F1EA',
		image: 'https://country-code-au6g.vercel.app/JE.svg',
		dial_code: '+44',
		minLength: 9, // Sama dengan UK (9-10)
		maxLength: 10,
		regexPattern: '^(\\d{9,10})$'
	},
	{
		name: 'Jamaica',
		code: 'JM',
		emoji: '🇯🇲',
		unicode: 'U+1F1EF U+1F1F2',
		image: 'https://country-code-au6g.vercel.app/JM.svg',
		dial_code: '+1876',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Jordan',
		code: 'JO',
		emoji: '🇯🇴',
		unicode: 'U+1F1EF U+1F1F4',
		image: 'https://country-code-au6g.vercel.app/JO.svg',
		dial_code: '+962',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Japan',
		code: 'JP',
		emoji: '🇯🇵',
		unicode: 'U+1F1EF U+1F1F5',
		image: 'https://country-code-au6g.vercel.app/JP.svg',
		dial_code: '+81',
		minLength: 9, // Umumnya 9 atau 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{9,10})$'
	},
	{
		name: 'Kenya',
		code: 'KE',
		emoji: '🇰🇪',
		unicode: 'U+1F1F0 U+1F1EA',
		image: 'https://country-code-au6g.vercel.app/KE.svg',
		dial_code: '+254',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Kyrgyzstan',
		code: 'KG',
		emoji: '🇰🇬',
		unicode: 'U+1F1F0 U+1F1EC',
		image: 'https://country-code-au6g.vercel.app/KG.svg',
		dial_code: '+996',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Cambodia',
		code: 'KH',
		emoji: '🇰🇭',
		unicode: 'U+1F1F0 U+1F1ED',
		image: 'https://country-code-au6g.vercel.app/KH.svg',
		dial_code: '+855',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Kiribati',
		code: 'KI',
		emoji: '🇰🇮',
		unicode: 'U+1F1F0 U+1F1EE',
		image: 'https://country-code-au6g.vercel.app/KI.svg',
		dial_code: '+686',
		minLength: 5, // Umumnya 5 digit
		maxLength: 5,
		regexPattern: '^(\\d{5})$'
	},
	{
		name: 'Comoros',
		code: 'KM',
		emoji: '🇰🇲',
		unicode: 'U+1F1F0 U+1F1F2',
		image: 'https://country-code-au6g.vercel.app/KM.svg',
		dial_code: '+269',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'St. Kitts & Nevis',
		code: 'KN',
		emoji: '🇰🇳',
		unicode: 'U+1F1F0 U+1F1F3',
		image: 'https://country-code-au6g.vercel.app/KN.svg',
		dial_code: '+1869',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'North Korea',
		code: 'KP',
		emoji: '🇰🇵',
		unicode: 'U+1F1F0 U+1F1F5',
		image: 'https://country-code-au6g.vercel.app/KP.svg',
		dial_code: '+850',
		minLength: 7, // Umumnya 7 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{7,10})$'
	},
	{
		name: 'South Korea',
		code: 'KR',
		emoji: '🇰🇷',
		unicode: 'U+1F1F0 U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/KR.svg',
		dial_code: '+82',
		minLength: 7, // Umumnya 7 sampai 11 digit
		maxLength: 11,
		regexPattern: '^(\\d{7,11})$'
	},
	{
		name: 'Kuwait',
		code: 'KW',
		emoji: '🇰🇼',
		unicode: 'U+1F1F0 U+1F1FC',
		image: 'https://country-code-au6g.vercel.app/KW.svg',
		dial_code: '+965',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Cayman Islands',
		code: 'KY',
		emoji: '🇰🇾',
		unicode: 'U+1F1F0 U+1F1FE',
		image: 'https://country-code-au6g.vercel.app/KY.svg',
		dial_code: '+ 345',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Kazakhstan',
		code: 'KZ',
		emoji: '🇰🇿',
		unicode: 'U+1F1F0 U+1F1FF',
		image: 'https://country-code-au6g.vercel.app/KZ.svg',
		dial_code: '+77',
		minLength: 10, // Umumnya 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Laos',
		code: 'LA',
		emoji: '🇱🇦',
		unicode: 'U+1F1F1 U+1F1E6',
		image: 'https://country-code-au6g.vercel.app/LA.svg',
		dial_code: '+856',
		minLength: 8, // Umumnya 8 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{8,10})$'
	},
	{
		name: 'Lebanon',
		code: 'LB',
		emoji: '🇱🇧',
		unicode: 'U+1F1F1 U+1F1E7',
		image: 'https://country-code-au6g.vercel.app/LB.svg',
		dial_code: '+961',
		minLength: 7, // Umumnya 7 atau 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{7,8})$'
	},
	{
		name: 'St. Lucia',
		code: 'LC',
		emoji: '🇱🇨',
		unicode: 'U+1F1F1 U+1F1E8',
		image: 'https://country-code-au6g.vercel.app/LC.svg',
		dial_code: '+1758',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Liechtenstein',
		code: 'LI',
		emoji: '🇱🇮',
		unicode: 'U+1F1F1 U+1F1EE',
		image: 'https://country-code-au6g.vercel.app/LI.svg',
		dial_code: '+423',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Sri Lanka',
		code: 'LK',
		emoji: '🇱🇰',
		unicode: 'U+1F1F1 U+1F1F0',
		image: 'https://country-code-au6g.vercel.app/LK.svg',
		dial_code: '+94',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Liberia',
		code: 'LR',
		emoji: '🇱🇷',
		unicode: 'U+1F1F1 U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/LR.svg',
		dial_code: '+231',
		minLength: 7, // Umumnya 7 atau 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{7,8})$'
	},
	{
		name: 'Lesotho',
		code: 'LS',
		emoji: '🇱🇸',
		unicode: 'U+1F1F1 U+1F1F8',
		image: 'https://country-code-au6g.vercel.app/LS.svg',
		dial_code: '+266',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Lithuania',
		code: 'LT',
		emoji: '🇱🇹',
		unicode: 'U+1F1F1 U+1F1F9',
		image: 'https://country-code-au6g.vercel.app/LT.svg',
		dial_code: '+370',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Luxembourg',
		code: 'LU',
		emoji: '🇱🇺',
		unicode: 'U+1F1F1 U+1F1FA',
		image: 'https://country-code-au6g.vercel.app/LU.svg',
		dial_code: '+352',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Latvia',
		code: 'LV',
		emoji: '🇱🇻',
		unicode: 'U+1F1F1 U+1F1FB',
		image: 'https://country-code-au6g.vercel.app/LV.svg',
		dial_code: '+371',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Libya',
		code: 'LY',
		emoji: '🇱🇾',
		unicode: 'U+1F1F1 U+1F1FE',
		image: 'https://country-code-au6g.vercel.app/LY.svg',
		dial_code: '+218',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Morocco',
		code: 'MA',
		emoji: '🇲🇦',
		unicode: 'U+1F1F2 U+1F1E6',
		image: 'https://country-code-au6g.vercel.app/MA.svg',
		dial_code: '+212',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Monaco',
		code: 'MC',
		emoji: '🇲🇨',
		unicode: 'U+1F1F2 U+1F1E8',
		image: 'https://country-code-au6g.vercel.app/MC.svg',
		dial_code: '+377',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Moldova',
		code: 'MD',
		emoji: '🇲🇩',
		unicode: 'U+1F1F2 U+1F1E9',
		image: 'https://country-code-au6g.vercel.app/MD.svg',
		dial_code: '+373',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Montenegro',
		code: 'ME',
		emoji: '🇲🇪',
		unicode: 'U+1F1F2 U+1F1EA',
		image: 'https://country-code-au6g.vercel.app/ME.svg',
		dial_code: '+382',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'St. Martin',
		code: 'MF',
		emoji: '🇲🇫',
		unicode: 'U+1F1F2 U+1F1EB',
		image: 'https://country-code-au6g.vercel.app/MF.svg',
		dial_code: '+590',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Madagascar',
		code: 'MG',
		emoji: '🇲🇬',
		unicode: 'U+1F1F2 U+1F1EC',
		image: 'https://country-code-au6g.vercel.app/MG.svg',
		dial_code: '+261',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Marshall Islands',
		code: 'MH',
		emoji: '🇲🇭',
		unicode: 'U+1F1F2 U+1F1ED',
		image: 'https://country-code-au6g.vercel.app/MH.svg',
		dial_code: '+692',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'North Macedonia',
		code: 'MK',
		emoji: '🇲🇰',
		unicode: 'U+1F1F2 U+1F1F0',
		image: 'https://country-code-au6g.vercel.app/MK.svg',
		dial_code: '+389',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Mali',
		code: 'ML',
		emoji: '🇲🇱',
		unicode: 'U+1F1F2 U+1F1F1',
		image: 'https://country-code-au6g.vercel.app/ML.svg',
		dial_code: '+223',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Myanmar (Burma)',
		code: 'MM',
		emoji: '🇲🇲',
		unicode: 'U+1F1F2 U+1F1F2',
		image: 'https://country-code-au6g.vercel.app/MM.svg',
		dial_code: '+95',
		minLength: 7, // Umumnya 7 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{7,10})$'
	},
	{
		name: 'Mongolia',
		code: 'MN',
		emoji: '🇲🇳',
		unicode: 'U+1F1F2 U+1F1F3',
		image: 'https://country-code-au6g.vercel.app/MN.svg',
		dial_code: '+976',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Macao SAR China',
		code: 'MO',
		emoji: '🇲🇴',
		unicode: 'U+1F1F2 U+1F1F4',
		image: 'https://country-code-au6g.vercel.app/MO.svg',
		dial_code: '+853',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Northern Mariana Islands',
		code: 'MP',
		emoji: '🇲🇵',
		unicode: 'U+1F1F2 U+1F1F5',
		image: 'https://country-code-au6g.vercel.app/MP.svg',
		dial_code: '+1670',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Martinique',
		code: 'MQ',
		emoji: '🇲🇶',
		unicode: 'U+1F1F2 U+1F1F6',
		image: 'https://country-code-au6g.vercel.app/MQ.svg',
		dial_code: '+596',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Mauritania',
		code: 'MR',
		emoji: '🇲🇷',
		unicode: 'U+1F1F2 U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/MR.svg',
		dial_code: '+222',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Montserrat',
		code: 'MS',
		emoji: '🇲🇸',
		unicode: 'U+1F1F2 U+1F1F8',
		image: 'https://country-code-au6g.vercel.app/MS.svg',
		dial_code: '+1664',
		minLength: 10, // NANP - 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{10})$'
	},
	{
		name: 'Malta',
		code: 'MT',
		emoji: '🇲🇹',
		unicode: 'U+1F1F2 U+1F1F9',
		image: 'https://country-code-au6g.vercel.app/MT.svg',
		dial_code: '+356',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Mauritius',
		code: 'MU',
		emoji: '🇲🇺',
		unicode: 'U+1F1F2 U+1F1FA',
		image: 'https://country-code-au6g.vercel.app/MU.svg',
		dial_code: '+230',
		minLength: 7, // Umumnya 7 atau 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{7,8})$'
	},
	{
		name: 'Maldives',
		code: 'MV',
		emoji: '🇲🇻',
		unicode: 'U+1F1F2 U+1F1FB',
		image: 'https://country-code-au6g.vercel.app/MV.svg',
		dial_code: '+960',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Malawi',
		code: 'MW',
		emoji: '🇲🇼',
		unicode: 'U+1F1F2 U+1F1FC',
		image: 'https://country-code-au6g.vercel.app/MW.svg',
		dial_code: '+265',
		minLength: 7, // Umumnya 7 sampai 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{7,9})$'
	},
	{
		name: 'Mexico',
		code: 'MX',
		emoji: '🇲🇽',
		unicode: 'U+1F1F2 U+1F1FD',
		image: 'https://country-code-au6g.vercel.app/MX.svg',
		dial_code: '+52',
		minLength: 8, // Umumnya 8 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{8,10})$'
	},
	{
		name: 'Malaysia',
		code: 'MY',
		emoji: '🇲🇾',
		unicode: 'U+1F1F2 U+1F1FE',
		image: 'https://country-code-au6g.vercel.app/MY.svg',
		dial_code: '+60',
		minLength: 7, // Umumnya 7 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{7,10})$'
	},
	{
		name: 'Mozambique',
		code: 'MZ',
		emoji: '🇲🇿',
		unicode: 'U+1F1F2 U+1F1FF',
		image: 'https://country-code-au6g.vercel.app/MZ.svg',
		dial_code: '+258',
		minLength: 8, // Umumnya 8 atau 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{8,9})$'
	},
	{
		name: 'Namibia',
		code: 'NA',
		emoji: '🇳🇦',
		unicode: 'U+1F1F3 U+1F1E6',
		image: 'https://country-code-au6g.vercel.app/NA.svg',
		dial_code: '+264',
		minLength: 7, // Umumnya 7 sampai 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{7,9})$'
	},
	{
		name: 'New Caledonia',
		code: 'NC',
		emoji: '🇳🇨',
		unicode: 'U+1F1F3 U+1F1E8',
		image: 'https://country-code-au6g.vercel.app/NC.svg',
		dial_code: '+687',
		minLength: 6, // Umumnya 6 digit
		maxLength: 6,
		regexPattern: '^(\\d{6})$'
	},
	{
		name: 'Niger',
		code: 'NE',
		emoji: '🇳🇪',
		unicode: 'U+1F1F3 U+1F1EA',
		image: 'https://country-code-au6g.vercel.app/NE.svg',
		dial_code: '+227',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Norfolk Island',
		code: 'NF',
		emoji: '🇳🇫',
		unicode: 'U+1F1F3 U+1F1EB',
		image: 'https://country-code-au6g.vercel.app/NF.svg',
		dial_code: '+672',
		minLength: 6, // Bervariasi, menggunakan 6-9
		maxLength: 9,
		regexPattern: '^(\\d{6,9})$'
	},
	{
		name: 'Nigeria',
		code: 'NG',
		emoji: '🇳🇬',
		unicode: 'U+1F1F3 U+1F1EC',
		image: 'https://country-code-au6g.vercel.app/NG.svg',
		dial_code: '+234',
		minLength: 7, // Umumnya 7 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{7,10})$'
	},
	{
		name: 'Nicaragua',
		code: 'NI',
		emoji: '🇳🇮',
		unicode: 'U+1F1F3 U+1F1EE',
		image: 'https://country-code-au6g.vercel.app/NI.svg',
		dial_code: '+505',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Netherlands',
		code: 'NL',
		emoji: '🇳🇱',
		unicode: 'U+1F1F3 U+1F1F1',
		image: 'https://country-code-au6g.vercel.app/NL.svg',
		dial_code: '+31',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'Norway',
		code: 'NO',
		emoji: '🇳🇴',
		unicode: 'U+1F1F3 U+1F1F4',
		image: 'https://country-code-au6g.vercel.app/NO.svg',
		dial_code: '+47',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Nepal',
		code: 'NP',
		emoji: '🇳🇵',
		unicode: 'U+1F1F3 U+1F1F5',
		image: 'https://country-code-au6g.vercel.app/NP.svg',
		dial_code: '+977',
		minLength: 8, // Umumnya 8 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{8,10})$'
	},
	{
		name: 'Nauru',
		code: 'NR',
		emoji: '🇳🇷',
		unicode: 'U+1F1F3 U+1F1F7',
		image: 'https://country-code-au6g.vercel.app/NR.svg',
		dial_code: '+674',
		minLength: 7, // Umumnya 7 digit
		maxLength: 7,
		regexPattern: '^(\\d{7})$'
	},
	{
		name: 'Niue',
		code: 'NU',
		emoji: '🇳🇺',
		unicode: 'U+1F1F3 U+1F1FA',
		image: 'https://country-code-au6g.vercel.app/NU.svg',
		dial_code: '+683',
		minLength: 4, // Umumnya 4 digit
		maxLength: 4,
		regexPattern: '^(\\d{4})$'
	},
	{
		name: 'New Zealand',
		code: 'NZ',
		emoji: '🇳🇿',
		unicode: 'U+1F1F3 U+1F1FF',
		image: 'https://country-code-au6g.vercel.app/NZ.svg',
		dial_code: '+64',
		minLength: 8, // Umumnya 8 sampai 10 digit
		maxLength: 10,
		regexPattern: '^(\\d{8,10})$'
	},
	{
		name: 'Oman',
		code: 'OM',
		emoji: '🇴🇲',
		unicode: 'U+1F1F4 U+1F1F2',
		image: 'https://country-code-au6g.vercel.app/OM.svg',
		dial_code: '+968',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Panama',
		code: 'PA',
		emoji: '🇵🇦',
		unicode: 'U+1F1F5 U+1F1E6',
		image: 'https://country-code-au6g.vercel.app/PA.svg',
		dial_code: '+507',
		minLength: 8, // Umumnya 8 digit
		maxLength: 8,
		regexPattern: '^(\\d{8})$'
	},
	{
		name: 'Peru',
		code: 'PE',
		emoji: '🇵🇪',
		unicode: 'U+1F1F5 U+1F1EA',
		image: 'https://country-code-au6g.vercel.app/PE.svg',
		dial_code: '+51',
		minLength: 9, // Umumnya 9 digit
		maxLength: 9,
		regexPattern: '^(\\d{9})$'
	},
	{
		name: 'French Polynesia',
		code: 'PF',
		emoji: '🇵🇫',
		unicode: 'U+1F1F5 U+1F1EB',
		image: 'https://country-code-au6g.vercel.app/PF.svg',
		dial_code: '+689',
		minLength: 6, // Umumnya 6 digit
		maxLength: 6,
		regexPattern: '^(\\d{6})$'
	}
];
/**
 * Convert phone number ke E.164 format
 * Input: "081234567890" atau "62 812 345 67890" atau "+62812345678900"
 * Output: "+6281234567890"
 */
export function toE164Format(
	phoneNumber: string,
	countryCode?: string
): { number: string; error?: string } {
	if (!phoneNumber || phoneNumber.trim() === '') {
		return { number: '', error: 'Phone number is required' };
	}

	// Remove spaces, dashes, parentheses
	let cleaned = phoneNumber.replace(/[\s\-()]/g, '');

	// If starts with +, it's international format
	if (cleaned.startsWith('+')) {
		if (cleaned.length < 10 || cleaned.length > 15) {
			return { number: '', error: 'Invalid phone number length' };
		}
		return { number: cleaned };
	}

	// If starts with 0, it's local format
	if (cleaned.startsWith('0')) {
		if (!countryCode) {
			return { number: '', error: 'Country code is required for local format' };
		}

		const country = CountryList.find((c) => c.code === countryCode.toUpperCase());
		if (!country) {
			return { number: '', error: 'Invalid country code' };
		}

		// Replace 0 with country dial code
		const dialCode = country.dial_code.replace('+', '');
		cleaned = dialCode + cleaned.slice(1);
		return { number: '+' + cleaned };
	}

	// If no + and doesn't start with 0, assume it needs country code prefix
	if (countryCode) {
		const country = CountryList.find((c) => c.code === countryCode.toUpperCase());
		if (!country) {
			return { number: '', error: 'Invalid country code' };
		}

		const dialCode = country.dial_code.replace('+', '');
		return { number: '+' + dialCode + cleaned };
	}

	return { number: '', error: 'Unable to format phone number' };
}

/**
 * Format phone number untuk display
 * E.164: "+6281234567890" → "62 812-345-67890"
 */
export function formatPhoneForDisplay(phoneNumber: string): string {
	if (!phoneNumber || !phoneNumber.startsWith('+')) {
		return phoneNumber;
	}

	// Remove +
	const cleaned = phoneNumber.slice(1);

	// For common country codes, format differently
	if (cleaned.startsWith('62')) {
		// Indonesia: +62 812-345-67890
		return `+${cleaned.slice(0, 2)} ${cleaned.slice(2, 5)}-${cleaned.slice(5, 8)}-${cleaned.slice(8)}`;
	}

	if (cleaned.startsWith('1')) {
		// USA: +1 (212) 345-6789
		return `+${cleaned.slice(0, 1)} (${cleaned.slice(1, 4)}) ${cleaned.slice(4, 7)}-${cleaned.slice(7)}`;
	}

	// Generic format: +XX XXX-XXX-XXXX
	return `+${cleaned.slice(0, 2)} ${cleaned.slice(2, 5)}-${cleaned.slice(5, 8)}-${cleaned.slice(8)}`;
}

/**
 * Detect country code dari phone number
 */
export function detectCountryFromPhone(phoneNumber: string): string | null {
	if (!phoneNumber.startsWith('+')) {
		return null;
	}

	const dialCode = phoneNumber.slice(1);

	for (const country of CountryList) {
		if (dialCode.startsWith(country.dial_code.replace('+', ''))) {
			return country.code;
		}
	}

	return null;
}

/**
 * Get country dial code
 */
export function getCountryDialCode(countryCode: string): string | null {
	const country = CountryList.find((c) => c.code === countryCode.toUpperCase());
	return country?.dial_code || null;
}

/**
 * Get country name dari code
 */
export function getCountryName(countryCode: string): string | null {
	const country = CountryList.find((c) => c.code === countryCode.toUpperCase());
	return country?.name || null;
}

/**
 * Validate E.164 format phone number
 */
export function isValidE164(phoneNumber: string): boolean {
	const e164Regex = /^\+[1-9]\d{1,14}$/;
	return e164Regex.test(phoneNumber);
}

/**
 * Extract country code dari E.164 number
 */
export function extractCountryDialCode(e164Number: string): string {
	if (!e164Number.startsWith('+')) {
		return '';
	}

	// Try to find matching country
	for (const country of CountryList) {
		const dialCode = country.dial_code.replace('+', '');
		if (e164Number.slice(1).startsWith(dialCode)) {
			return dialCode;
		}
	}

	return '';
}

/**
 * Extract phone number without country code
 */
export function extractPhoneWithoutCC(e164Number: string): string {
	if (!e164Number.startsWith('+')) {
		return e164Number;
	}

	const dialCode = extractCountryDialCode(e164Number);
	if (!dialCode) {
		return e164Number.slice(1);
	}

	return e164Number.slice(dialCode.length + 1);
}

/**
 * Get all countries untuk dropdown
 */
export function getCountriesForSelect() {
	return CountryList.map((c) => ({
		code: c.code,
		name: c.name,
		dial_code: c.dial_code,
		flag: c.image
	}));
}

/**
 * Compare two phone numbers (ignore formatting)
 */
export function comparePhoneNumbers(phone1: string, phone2: string): boolean {
	const clean1 = phone1.replace(/[\s\-()]/g, '');
	const clean2 = phone2.replace(/[\s\-()]/g, '');

	return clean1 === clean2;
}

/**
 * Generate WhatsApp link dari phone number
 */
export function generateWhatsAppLink(phoneNumber: string, text?: string): string {
	// Remove + dari E.164 format
	const cleanPhone = phoneNumber.replace(/\+/g, '');

	if (text) {
		const encodedText = encodeURIComponent(text);
		return `https://wa.me/${cleanPhone}?text=${encodedText}`;
	}

	return `https://wa.me/${cleanPhone}`;
}

/**
 * Generate WhatsApp deep link (untuk mobile)
 */
export function generateWhatsAppDeepLink(phoneNumber: string, text?: string): string {
	const cleanPhone = phoneNumber.replace(/\+/g, '');

	if (text) {
		const encodedText = encodeURIComponent(text);
		return `wa.me/${cleanPhone}?text=${encodedText}`;
	}

	return `wa.me/${cleanPhone}`;
}

/**
 * Validate phone number dengan country context
 */
export function validatePhoneForCountry(
	phoneNumber: string,
	countryCode: string
): { valid: boolean; error?: string } {
	const country = CountryList.find((c) => c.code === countryCode.toUpperCase());
	if (!country) {
		return { valid: false, error: 'Invalid country code' };
	}

	// Convert to E.164
	const { number, error } = toE164Format(phoneNumber, countryCode);
	if (error) {
		return { valid: false, error };
	}

	// Validate format
	if (!isValidE164(number)) {
		return { valid: false, error: 'Invalid phone number format' };
	}

	return { valid: true };
}

/**
 * Parse phone input (untuk form input yang accepts berbagai format)
 */
export function parsePhoneInput(
	input: string,
	countryCode?: string
): {
	valid: boolean;
	number?: string;
	countryCode?: string;
	error?: string;
} {
	if (!input || input.trim() === '') {
		return { valid: false, error: 'Phone number is required' };
	}

	// Detect country dari phone jika ada
	let detectedCC = countryCode;
	if (input.startsWith('+')) {
		const detected = detectCountryFromPhone(input);
		if (detected) {
			detectedCC = detected;
		}
	}

	if (!detectedCC) {
		return { valid: false, error: 'Country code could not be determined' };
	}

	const { number, error } = toE164Format(input, detectedCC);
	if (error) {
		return { valid: false, error };
	}

	if (!isValidE164(number)) {
		return { valid: false, error: 'Invalid phone number format' };
	}

	return { valid: true, number, countryCode: detectedCC };
}
