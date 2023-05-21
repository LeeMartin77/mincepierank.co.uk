import type { types } from 'cassandra-driver';
import type { MakerPieRanking } from './types';

export function rowToObject<T>(row: types.Row): T {
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	const constructed: any = {};
	row.keys().forEach((key) => (constructed[key] = row.get(key)));
	return constructed as T;
}

export function calculateAverage(
	mapped: Pick<MakerPieRanking, 'filling' | 'pastry' | 'topping' | 'value' | 'looks'>
): number {
	return (mapped.filling + mapped.pastry + mapped.topping + mapped.value + mapped.looks) / 5;
}
