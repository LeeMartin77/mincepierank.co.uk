export const imgprssrPrefix = (link: string) => {
  if (link.startsWith('http://') || link.startsWith('https://')) {
    return link;
  }
  return `${process.env.IMGPRSSR_ROOT || 'https://static.mincepierank.co.uk'}${link}`;
};
