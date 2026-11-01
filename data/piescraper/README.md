# Piescraper

## The problem

Most supermarkets and pie-sellers have some form of bot protection on their website - totally valid.

If this was a big commercial project, we'd pay for access, but for a hobby project like this, it makes sense to just copy down the data to a CSV and upload it - which is at least progress against the edit-each-pie-individually model of the past.

Problem is this is massively boring gruntwork, so the goal is to simplify it a bit.

## The extension

`webscraper/` is a Firefox extension that turns product pages you browse by hand into rows for the bulk importer.

Load it from `about:debugging#/runtime/this-firefox` -> "Load Temporary Add-on" -> pick `webscraper/manifest.json` (Firefox drops it again on restart, so reload it each session).

Then:

- Browse to a product page on a retailer site.
- Open the Piescraper toolbar popup and hit "Scrape current tab". It reads the JSON-LD product data, `og:` meta tags and page images, and guesses id, makerid, price, pack count and categories.
- Entries accumulate in extension storage keyed by `makerid/id`, so scraping the same pie again just merges new image candidates in without losing your edits.
- Fix anything wrong in the popup table (id and makerid form the key, price is in pence, categories are slugs from `data/generic/categories.sql`).
- Pick the year, hit "Export zip". The images are fetched at this point (in your browsing session, so cookies get past the image cdn), packed with `piedata.csv` into `piedata-<year>.zip`, and you can upload that straight to the admin bulk import page.

The pure logic (ids, prices, pack counts, category guesses, csv, zip) lives in `webscraper/lib/derive.js` and `webscraper/lib/zip.js` and is unit tested with `node --test webscraper/test/*.test.mjs` (node >= 22). `lib/extract.js` only collects raw DOM facts.

If a page has no json-ld, or json-ld in a shape the extension does not understand, open "AI fill" in the popup and point it at an OpenAI-compatible endpoint (base url + api key + model). It only runs when scraping has left gaps, and the model sees the page's raw json-ld and visible text but can only fill the missing fields: categories are filtered to the real slugs, images must be urls the page itself referenced, a real product code from the data replaces a weak guessed id, and ids you have edited are left alone. A failing api call never loses the scraped entry.

## Legacy flow

The older dump-everything approach (`organiser/`) is kept around but the extension replaced it: structured extraction gets us most of the way there, and the optional AI fill only needs the json-ld and page text for the rest.

## Considerations

- We want a script that doesn't tie us to a one-and-done upload - the prompts and structure should focus on the idea we might be coming back later to add more (seasonal, after all)
