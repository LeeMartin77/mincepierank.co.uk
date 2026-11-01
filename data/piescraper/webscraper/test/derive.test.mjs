import test from "node:test";
import assert from "node:assert/strict";
import { loadLibs } from "./helpers.mjs";

const PS = loadLibs(["lib/derive.js"]).PS;

// Arrays built inside the vm have the vm's Array prototype; copy to a host
// array so deepEqual compares plain values.
function eqArray(actual, expected) {
  assert.deepEqual(Array.from(actual), expected);
}

test("slugify", function () {
  assert.equal(PS.slugify("Specially Selected Mince Pies!"), "specially-selected-mince-pies");
  assert.equal(PS.slugify("  ---  "), "");
});

test("makeridFromUrl", function () {
  assert.equal(PS.makeridFromUrl("https://www.tesco.com/groceries/en-GB/products/906005527"), "tesco");
  assert.equal(PS.makeridFromUrl("https://groceries.aldi.co.uk/EN-Online/shop/p-specially-selected/4061462560620"), "aldi");
  assert.equal(PS.makeridFromUrl("https://www.marksandspencer.com/p/snowballs/p200804242"), "marks-and-spencer");
  assert.equal(PS.makeridFromUrl("https://www.morrisons.com/product/christmas-mince-pies-256647233/"), "morrisons");
  assert.equal(PS.makeridFromUrl("https://www.piefarm.co.uk/pies"), "piefarm");
});

test("deriveId prefers codes then url then name", function () {
  assert.equal(PS.deriveId({ url: "https://x.com/p/whatever", name: "Mince Pies", gtin: "0123456789012" }), "0123456789012");
  assert.equal(PS.deriveId({ url: "https://x.com/products/mince-pies/906005527", name: "Mince Pies" }), "906005527");
  assert.equal(PS.deriveId({ url: "https://x.com/products/christmas-mince-pies", name: "Mince Pies" }), "christmas-mince-pies");
  assert.equal(PS.deriveId({ url: "https://x.com/", name: "Giant Mince Pie" }), "giant-mince-pie");
});

test("priceToPence", function () {
  assert.equal(PS.priceToPence(2.49), 249);
  assert.equal(PS.priceToPence("£2.49"), 249);
  assert.equal(PS.priceToPence("1,99"), 199);
  assert.equal(PS.priceToPence("abc"), null);
  assert.equal(PS.priceToPence(""), null);
});

test("packCountFromText", function () {
  assert.equal(PS.packCountFromText("Mince Pies 6 Pack"), 6);
  assert.equal(PS.packCountFromText("Pack of 4 mince pies"), 4);
  assert.equal(PS.packCountFromText("12 x 25g mince pies"), 12);
  assert.equal(PS.packCountFromText("8 mince pies"), 8);
  assert.equal(PS.packCountFromText("Mince Pie"), null);
});

test("guessCategories", function () {
  eqArray(PS.guessCategories("Specially Selected Snowballs Mini Mince Pies"), ["mini", "premium"]);
  eqArray(PS.guessCategories("Salted Caramel Mince Pies"), ["salted-caramel"]);
  eqArray(PS.guessCategories("Caramel Mince Pies"), ["caramel"]);
  eqArray(PS.guessCategories("Vegan Mince Pies"), ["vegan"]);
  eqArray(PS.guessCategories("All Butter Mince Pies"), ["classic"]);
  eqArray(PS.guessCategories("dog food"), []);
});

test("guessCategories never emits unknown slugs", function () {
  const samples = [
    "Signature Maple & Pecan Frangipane Mini Mince Pies",
    "Giant Spiced Mince Pie with Brandy, vegan",
    "Specially Selected Gluten Free Salted Caramel & Chocolate Mince Pies, iced",
    "Collected Brown Butter Mince Pies with Cognac, Organic, puff pastry, orange, apple, gingerbread, almond, toffee, flapjack, crumble, lattice, rum, caramel"
  ];
  for (const s of samples) {
    for (const slug of PS.guessCategories(s)) {
      assert.ok(PS.CATEGORIES.includes(slug), "unknown slug " + slug);
    }
  }
});

test("pickProduct prefers the product with name and price", function () {
  const picked = PS.pickProduct([
    { name: "Breadcrumb", price: null, imageUrls: [] },
    { name: "Mince Pies", price: 249, imageUrls: ["https://x/y.jpg"] }
  ]);
  assert.equal(picked.name, "Mince Pies");
});

test("scoreImages ranks og and json-ld images first", function () {
  const candidates = PS.scoreImages({
    ogImage: "https://cdn.shop/hero.jpg",
    ldImageUrls: ["https://cdn.shop/ld.jpg"],
    images: [
      { url: "https://cdn.shop/hero.jpg", alt: "pie", w: 400, h: 400, productish: true },
      { url: "https://cdn.shop/ld.jpg", alt: "", w: 300, h: 300, productish: true },
      { url: "https://cdn.shop/trustpilot-logo.png", alt: "reviews", w: 500, h: 500, productish: false },
      { url: "https://cdn.shop/side.jpg", alt: "mince pies", w: 350, h: 250, productish: true }
    ]
  });
  assert.equal(candidates[0].url, "https://cdn.shop/hero.jpg");
  assert.equal(candidates[1].url, "https://cdn.shop/ld.jpg");
  assert.ok(candidates[candidates.length - 1].url === "https://cdn.shop/trustpilot-logo.png");
});

test("buildEntry assembles a bulk-ready row", function () {
  const entry = PS.buildEntry({
    url: "https://www.greggs.co.uk/products/mince-pie/123456789",
    title: "Mince Pie | Greggs",
    ogTitle: "Mince Pie",
    ogImage: "https://www.greggs.co.uk/img/mince-pie.jpg",
    products: [
      {
        name: "Mince Pie",
        description: "Our classic mince pie",
        imageUrls: ["https://www.greggs.co.uk/img/mince-pie.jpg"],
        price: "1.25",
        gtin: null,
        sku: null
      }
    ],
    images: [],
    priceMeta: null
  });
  assert.equal(entry.makerid, "greggs");
  assert.equal(entry.id, "123456789");
  assert.equal(entry.idSource, "url-code");
  assert.equal(entry.name, "Mince Pie");
  assert.equal(entry.priceInPence, 125);
  assert.equal(entry.fresh, false);
  eqArray(entry.categories, ["classic"]);
  assert.equal(entry.image, "https://www.greggs.co.uk/img/mince-pie.jpg");
});

test("buildCsv quotes fields and matches the importer header", function () {
  const csv = PS.buildCsv("2026", [
    {
      year: "2026",
      makerid: "greggs",
      id: "mince-pie",
      displayname: 'The "Best" Mince Pie, 6 pack',
      fresh: "no",
      pack_count: 6,
      pack_price_in_pence: 299,
      category_slugs: "classic,mini",
      web_link: "https://example.com/p/1",
      image_file: "images/greggs/mince-pie.jpg"
    }
  ]);
  const lines = csv.trim().split("\n");
  assert.equal(
    lines[0],
    "year,makerid,id,displayname,fresh,pack_count,pack_price_in_pence,category_slugs,web_link,image_file"
  );
  assert.ok(lines[1].includes('"The ""Best"" Mince Pie, 6 pack"'));
});

function baseEntry(over) {
  return Object.assign(
    {
      makerid: "weird",
      id: "festive-mince-pies",
      idSource: "name-slug",
      url: "https://www.weird-supermarket.co.uk/p/99",
      name: "Festive Mince Pies",
      description: "",
      priceInPence: null,
      packCount: null,
      fresh: false,
      categories: [],
      image: "",
      imageCandidates: [],
      editedFields: []
    },
    over || {}
  );
}

test("missingFields flags gaps but not strong or edited ids", function () {
  eqArray(PS.missingFields(baseEntry({
    priceInPence: 399, packCount: 6, categories: ["classic"],
    image: "https://x/y.jpg", idSource: "url-code"
  })), []);
  eqArray(PS.missingFields(baseEntry({
    priceInPence: 399, packCount: 6, categories: ["classic"],
    image: "https://x/y.jpg", id: "user-chosen", idSource: "edited"
  })), []);
  eqArray(PS.missingFields(baseEntry()), ["price", "pack_count", "categories", "image", "id"]);
  assert.ok(PS.needsAiFill(baseEntry()));
});

test("aiMessages sends raw json-ld, gaps and page-ish image urls", function () {
  const facts = {
    url: "https://www.weird-supermarket.co.uk/p/99",
    title: "Festive Mince Pies",
    rawLd: ['{"@type":"WeirdThing","item":{"price":"3.99","imgs":["https://cdn.weird/img_1.jpg"]}}'],
    bodyText: "Festive Mince Pies 6 pack £3.99 " + "x".repeat(7000)
  };
  const msgs = PS.aiMessages(facts, baseEntry());
  assert.equal(msgs[0].role, "system");
  assert.ok(msgs[0].content.includes("salted-caramel"));
  const user = JSON.parse(msgs[1].content);
  assert.deepEqual(user.jsonld, facts.rawLd);
  eqArray(user.missing, ["price", "pack_count", "categories", "image", "id"]);
  assert.ok(user.image_urls.includes("https://cdn.weird/img_1.jpg"));
  assert.ok(user.page_text.length <= 6000);
});

test("parseAiJson copes with prose and code fences", function () {
  // parsed objects come from the vm realm, so compare the json text instead
  assert.equal(JSON.stringify(PS.parseAiJson('{"name": "x"}')), '{"name":"x"}');
  assert.equal(JSON.stringify(PS.parseAiJson('Sure!\n```json\n{"name": "x"}\n```')), '{"name":"x"}');
  assert.throws(function () { PS.parseAiJson("{name: not json}"); }, function (err) {
    return err.name === "SyntaxError";
  });
  assert.throws(function () { PS.parseAiJson(""); }, function (err) {
    return err.name === "Error" && err.message.includes("no json object");
  });
});

test("aiFill only fills gaps and validates every value", function () {
  const entry = baseEntry();
  const res = PS.aiFill(entry, {
    name: "Ignored Renaming",
    price_gbp: 3.99,
    pack_count: 6,
    categories: ["NotACategory", "Mini", "classic", "classic"],
    fresh: true,
    image: "https://evil.example.com/made-up.jpg",
    id: "festive-mince-pies-6pk"
  }, ["https://cdn.weird/img_1.jpg"]);
  assert.equal(entry.name, "Festive Mince Pies");
  assert.equal(entry.priceInPence, 399);
  assert.equal(entry.packCount, 6);
  eqArray(entry.categories, ["classic", "mini"]);
  assert.equal(entry.fresh, true);
  assert.equal(entry.image, "");
  assert.equal(entry.id, "festive-mince-pies-6pk");
  assert.equal(entry.idSource, "ai");
  eqArray(res.filled, ["priceInPence", "packCount", "categories", "fresh", "id"]);
});

test("aiFill accepts allowed images and product codes over weak ids", function () {
  const entry = baseEntry({ id: "christmas-mince-pies", idSource: "url-slug" });
  PS.aiFill(entry, {
    image: "https://cdn.weird/img_1.jpg",
    id: "0 123 4567 8901",
    price_gbp: 0
  }, ["https://cdn.weird/img_1.jpg"]);
  assert.equal(entry.image, "https://cdn.weird/img_1.jpg");
  assert.equal(entry.id, "012345678901");
  assert.equal(entry.idSource, "ai");
  assert.equal(entry.priceInPence, null);

  const byCode = baseEntry({ id: "906005527", idSource: "code" });
  PS.aiFill(byCode, { id: "mince-pies" }, []);
  assert.equal(byCode.id, "906005527");

  const edited = baseEntry({ id: "user-chosen", idSource: "edited" });
  PS.aiFill(edited, { id: "ai-guess" }, []);
  assert.equal(edited.id, "user-chosen");
});
