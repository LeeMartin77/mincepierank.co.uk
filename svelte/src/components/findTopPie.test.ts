import { MakerPie, PieRankingSummary } from '../system/storage';
import { findTopPie } from './findTopPie';
describe('findTopPie', () => {
	test('Empty Arrays :: Returns Undefined', () => {
		expect(findTopPie([], [])).toBeUndefined();
		expect(findTopPie([{} as MakerPie], [])).toBeUndefined();
		expect(findTopPie([], [{} as PieRankingSummary])).toBeUndefined();
	});
	test('Returns the average top pie', () => {
		const makerPies: MakerPie[] = [
			{
				makerid: 'third-makerid',
				id: 'third-pieid',
				displayname: 'Third Pie',
				fresh: false,
				labels: [],
				image_file: 'somefile.jpg',
				web_link: 'somelink.html',
				pack_count: 6,
				pack_price_in_pence: 100
			},
			{
				makerid: 'second-makerid',
				id: 'second-pieid',
				displayname: 'Second Pie',
				fresh: false,
				labels: [],
				image_file: 'somefile.jpg',
				web_link: 'somelink.html',
				pack_count: 6,
				pack_price_in_pence: 100
			},
			{
				makerid: 'first-makerid',
				id: 'first-pieid',
				displayname: 'First Pie',
				fresh: false,
				labels: [],
				image_file: 'somefile.jpg',
				web_link: 'somelink.html',
				pack_count: 6,
				pack_price_in_pence: 100
			}
		];

		const pieRankings: PieRankingSummary[] = makerPies.map((mp, i) => {
			return {
				makerid: mp.makerid,
				pieid: mp.id,
				value: i,
				filling: i,
				pastry: i,
				topping: i,
				looks: i,
				count: 1,
				average: i
			};
		});

		const [topPie, topPieRanking] = findTopPie(makerPies, pieRankings) ?? [];

		expect(topPie?.id).toBe('first-pieid');
		expect(topPieRanking?.pieid).toBe('first-pieid');
	});
});
