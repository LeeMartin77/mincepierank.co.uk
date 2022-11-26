
export function formatPrice(pencePrice: number): string {
    if (pencePrice > 99) {
        return `£${(pencePrice / 100).toFixed(2)}`;
    }
    return `£0.${(pencePrice).toFixed(0)}`;
}