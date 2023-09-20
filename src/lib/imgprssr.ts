export const imgprssrPrefix = (link: string, root = 'https://static.mincepierank.co.uk') => {
  if (link.startsWith('http://') || link.startsWith('https://')) {
    return link;
  }
  return `${root}${link}`;
};
