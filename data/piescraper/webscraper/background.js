browser.action.onClicked.addListener(async (tab) => {
  await browser.scripting.executeScript({
    target: { tabId: tab.id },
    files: ["scraper.js"]
  });
});