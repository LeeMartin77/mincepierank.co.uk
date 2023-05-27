import { ppCategory } from './formatCategory';
import { describe, test, expect } from 'vitest';

describe('ppCategory', () => {
	const testcases = [
		['something', 'Something'],
		['another', 'Another'],
		['why-not-me', 'Why Not Me'],
		['another-one', 'Another One'],
		['yes-no-maybe', 'Yes No Maybe']
	];

	test.each(testcases)('%s outputs %s', (input, expected) => {
		expect(ppCategory(input)).toBe(expected);
	});
});
