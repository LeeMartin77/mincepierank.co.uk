export const formatCurrency = (pence: number): string => {
    const inPounds = pence / 100;
    return `£${inPounds.toFixed(2)}`
}