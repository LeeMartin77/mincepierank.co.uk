import { describe, expect, test } from "vitest";
import { formatCurrency } from "./formatCurrency";

describe("formatCurrency", () => {
    const tests: [number, string][] = [
        [20, "£0.20"],
        [120, "£1.20"],
        [1059, "£10.59"],
        [20.0000001, "£0.20"],
        [1059.99999999, "£10.60"],
    ]
    test.each(tests)("Given %s outputs %s", (pence, expected) => {
        expect(formatCurrency(pence)).toBe(expected);
    })
})