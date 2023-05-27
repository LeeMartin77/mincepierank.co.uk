import type { types } from 'cassandra-driver';
import { rowToObject } from './utilities';
import { describe, test, expect } from 'vitest';

describe('rowToObject', () => {
	test('Maps row to object', () => {
		const fakeObject: { [key: string]: any } = {
			propOne: 'something',
			propTheSecond: 2,
			propThird: 'hello'
		};
		const testRow = {
			keys: () => {
				return Object.keys(fakeObject);
			},
			get: (key: string) => {
				return fakeObject[key];
			}
		};

		const result = rowToObject(testRow as types.Row);

		expect(result).toStrictEqual(fakeObject);
	});
});
