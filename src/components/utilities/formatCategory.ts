export function ppCategory(catid: string): string {
  return catid.split('-').reduce((prev, curr) => {
    if (prev.length > 0) {
      prev += ' ';
    }
    return prev + curr.substring(0, 1).toLocaleUpperCase() + curr.substring(1);
  }, '');
}
