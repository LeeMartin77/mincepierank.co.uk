// Pure scraping logic, shared by the page extractor, the background page
// (importScripts) and the node tests. Attaches to the global scope so it
// can also be loaded as a classic script via scripting.executeScript.

var PS = (function () {
  "use strict";

  var IMAGE_JUNK = /logo|icon|sprite|banner|cookie|avatar|nav|footer|payment|trustpilot|review|rating|social|flag|badge|app-?store/i;

  function slugify(s) {
    return String(s || "")
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, "-")
      .replace(/^-+|-+$/g, "")
      .slice(0, 64);
  }

  function parseUrl(url) {
    try {
      return new URL(url);
    } catch (e) {
      return null;
    }
  }

  // brand websites in the db are free text ("https://www.tesco.com",
  // "tesco.com"), so reduce both sides to a bare hostname before comparing.
  function hostOf(url) {
    var s = String(url || "").trim();
    if (!s) {
      return "";
    }
    var u = parseUrl(/^[a-z][a-z0-9+.-]*:\/\//i.test(s) ? s : "https://" + s);
    var host = u ? u.hostname : s;
    return host.toLowerCase().replace(/^www\./, "");
  }

  // brands is the live list from /api/brands ([{id, website}]); a host match
  // there gives the canonical makerid. Without one the makerid is just a slug
  // derived from the host, which the import will create as a new brand.
  function makeridFromUrl(url, brands) {
    var u = parseUrl(url);
    if (!u) {
      return "unknown";
    }
    var host = u.hostname.toLowerCase().replace(/^www\./, "");
    if (brands && brands.length) {
      var probes = [host];
      if (host.startsWith("groceries.")) {
        probes.push(host.slice("groceries.".length));
      }
      for (var i = 0; i < brands.length; i++) {
        var b = brands[i] || {};
        var bh = hostOf(b.website);
        if (!b.id || !bh) {
          continue;
        }
        if (probes.indexOf(bh) !== -1) {
          return b.id;
        }
      }
    }
    if (host.startsWith("groceries.")) {
      host = host.slice("groceries.".length);
    }
    return slugify(host.replace(/\.(co\.uk|com|uk|ie|net)$/g, "")) || "unknown";
  }

  // id is what converges on re-scrape (it is part of the primary key in the
  // bulk csv), so prefer stable numeric product codes over anything else.
  function deriveId(facts) {
    return deriveIdDetailed(facts).id;
  }

  function deriveIdDetailed(facts) {
    var code = firstDigits(facts.gtin) || firstDigits(facts.sku);
    if (code && code.length >= 8) {
      return { id: code, source: "code" };
    }
    var u = parseUrl(facts.url);
    if (u) {
      var segs = u.pathname.split("/").filter(Boolean);
      for (var i = segs.length - 1; i >= 0; i--) {
        var num = firstDigits(decodeURIComponent(segs[i]));
        if (num && num.length >= 8) {
          return { id: num, source: "url-code" };
        }
      }
      for (var j = segs.length - 1; j >= 0; j--) {
        var seg = slugify(decodeURIComponent(segs[j]));
        if (seg && !/^(p|pd|pdp|product|products|item|sku|c)$/g.test(seg) && !/^\d+$/.test(seg)) {
          return { id: seg, source: "url-slug" };
        }
      }
    }
    var nameSlug = slugify(facts.name);
    return { id: nameSlug || "unknown", source: nameSlug ? "name-slug" : "unknown" };
  }

  function firstDigits(s) {
    var m = String(s || "").match(/\d{8,}/);
    return m ? m[0] : null;
  }

  function priceToPence(v) {
    if (v === null || v === undefined || v === "") {
      return null;
    }
    if (typeof v === "number" && isFinite(v)) {
      return Math.round(v * 100);
    }
    var s = String(v).replace(/[£€\s]/g, "").replace(/,/g, ".");
    var m = s.match(/^\d+(\.\d+)?$/);
    if (!m) {
      return null;
    }
    return Math.round(parseFloat(m[0]) * 100);
  }

  function packCountFromText(text) {
    var t = String(text || "");
    var patterns = [
      /(\d+)\s*[-\s]?pack\b/i,
      /\b(?:pack|box|tray|bag|multipack)\s*of\s*(\d+)\b/i,
      /\b(\d+)\s+(?:mince\s+)?pies\b/i,
      /\b(\d+)\s*(?:x|×)\s*\d*\s*[a-z]*\s*(?:mince\s+)?pie/i
    ];
    for (var i = 0; i < patterns.length; i++) {
      var m = t.match(patterns[i]);
      if (m) {
        var n = parseInt(m[1], 10);
        if (n > 0 && n <= 100) {
          return n;
        }
      }
    }
    return null;
  }

  var CATEGORY_PATTERNS = [
    ["salted-caramel", /salted[\s-]?caramel/i],
    ["caramel", /\bcaramel\b/i],
    ["gluten-free", /gluten[\s-]?free/i],
    ["brown-butter", /brown(?:ed)?\s+butter/i],
    ["puff-pastry", /puff\s+pastry/i],
    ["shortcrust", /short\s?crust/i],
    ["frangipane", /frangipane/i],
    ["gingerbread", /gingerbread/i],
    ["cognac", /\b(cognac|brandy)\b/i],
    ["rum", /\brum\b/i],
    ["alcohol", /\b(alcohol|whisky|baileys|port|sherry)\b/i],
    ["vegan", /vegan/i],
    ["organic", /organic/i],
    ["mini", /\bmini(\s|$|-)/i],
    ["giant", /\b(giant|jumbo)\b/i],
    ["premium", /(specially selected|signature|collected|taste the difference|finest|premium|luxury|p[aâ]tisserie)/i],
    ["iced", /\b(iced|glazed|icing)\b/i],
    ["maple", /maple/i],
    ["orange", /\borange\b/i],
    ["chocolate", /chocolate|\bchoc\b/i],
    ["toffee", /toffee/i],
    ["apple", /\bapple\b/i],
    ["almond", /\balmond\b/i],
    ["pecan", /pecan/i],
    ["flapjack", /flapjack/i],
    ["crumble", /crumble/i],
    ["lattice", /lattice/i],
    ["bakery", /\bbakery\b|baked\s+fresh/i],
    ["spiced", /\bspiced\b/i]
  ];

  function guessCategories(text) {
    var t = String(text || "");
    var found = [];
    for (var i = 0; i < CATEGORY_PATTERNS.length; i++) {
      var slug = CATEGORY_PATTERNS[i][0];
      if (CATEGORY_PATTERNS[i][1] && CATEGORY_PATTERNS[i][1].test(t) && found.indexOf(slug) === -1) {
        found.push(slug);
      }
    }
    if (found.length === 0 && /mince\s+pie/i.test(t)) {
      found.push("classic");
    }
    if (found.indexOf("salted-caramel") !== -1 && found.indexOf("caramel") !== -1) {
      found.splice(found.indexOf("caramel"), 1);
    }
    return found.sort();
  }

  // live is the slug list from /api/categories; without it the pattern slugs
  // go through unfiltered (the pattern list mirrors what the site carries).
  function filterCategories(slugs, live) {
    if (!live || !live.length) {
      return slugs;
    }
    return slugs.filter(function (s) {
      return live.indexOf(s) !== -1;
    });
  }

  function pickProduct(products) {
    var best = null;
    var bestScore = -1;
    for (var i = 0; i < products.length; i++) {
      var p = products[i];
      var score = 0;
      if (p.name) {
        score += 2;
      }
      if (p.price !== null && p.price !== undefined) {
        score += 2;
      }
      if (p.imageUrls && p.imageUrls.length) {
        score += 1;
      }
      if ((p.gtin || p.sku) && /mince\s+pie/i.test(p.name || "")) {
        score += 3;
      }
      if (score > bestScore) {
        best = p;
        bestScore = score;
      }
    }
    return best;
  }

  // facts.images: [{url, alt, w, h, productish}] collected from the DOM;
  // facts.ogImage and facts.ldImageUrls seed synthetic candidates too so a
  // lazy-loaded hero image still ends up downloadable from the background.
  function scoreImages(facts) {
    var byUrl = {};
    function add(url, cand) {
      if (!url || !/^https?:/i.test(url)) {
        return;
      }
      var prev = byUrl[url];
      if (!prev || cand.score > prev.score) {
        byUrl[url] = cand;
      }
    }

    var ogPath = pathOf(facts.ogImage);
    var ldPaths = {};
    (facts.ldImageUrls || []).forEach(function (u) {
      var p = pathOf(u);
      if (p) {
        ldPaths[p] = true;
      }
    });

    (facts.images || []).forEach(function (img) {
      var p = pathOf(img.url);
      var score = 0;
      if (ogPath && p === ogPath) {
        score += 1000;
      }
      if (p && ldPaths[p]) {
        score += 900;
      }
      if (img.productish) {
        score += 200;
      }
      if (/mince|pie/i.test(img.alt || "")) {
        score += 150;
      }
      score += Math.min(img.w || 0, 800) + Math.min(img.h || 0, 800);
      if (IMAGE_JUNK.test((img.url || "") + " " + (img.alt || ""))) {
        score -= 600;
      }
      if ((img.w || 0) < 150 && (img.h || 0) < 150) {
        score -= 300;
      }
      add(img.url, { url: img.url, score: score });
    });

    [facts.ogImage].concat(facts.ldImageUrls || []).forEach(function (u) {
      if (u) {
        add(u, { url: u, score: 800 });
      }
    });

    var out = Object.keys(byUrl).map(function (k) {
      return byUrl[k];
    });
    out.sort(function (a, b) {
      return b.score - a.score;
    });
    return out.slice(0, 8);
  }

  function pathOf(url) {
    var u = parseUrl(url);
    return u ? u.hostname.toLowerCase() + u.pathname : null;
  }

  // ref is the site data fetched by background.js: {brands: [{id, website}],
  // categories: [slugs]}. Guessed slugs are filtered through the live list so
  // a slug the site has renamed can never poison the bulk import.
  function buildEntry(facts, ref) {
    ref = ref || {};
    var product = pickProduct(facts.products || []);
    var name = (product && product.name) || facts.ogTitle || facts.title || "";
    var derived = deriveIdDetailed({
      url: facts.url,
      name: name,
      gtin: product && product.gtin,
      sku: product && product.sku
    });
    var id = derived.id;
    var price = null;
    if (product && product.price !== null && product.price !== undefined) {
      price = priceToPence(product.price);
    }
    if (price === null && facts.priceMeta) {
      price = priceToPence(facts.priceMeta);
    }
    var text = [name, product && product.description, facts.ogDescription].join(" ");
    var candidates = scoreImages({
      images: facts.images || [],
      ogImage: facts.ogImage,
      ldImageUrls: product && product.imageUrls
    });
    return {
      makerid: makeridFromUrl(facts.url, ref.brands),
      id: id,
      idSource: derived.source,
      url: facts.url,
      name: name,
      description: (product && product.description) || facts.ogDescription || "",
      priceInPence: price,
      packCount: packCountFromText(text),
      fresh: false,
      categories: filterCategories(guessCategories(text), ref.categories),
      image: candidates.length ? candidates[0].url : (facts.ogImage || ""),
      imageCandidates: candidates,
      capturedAt: new Date().toISOString(),
      editedFields: []
    };
  }

  // ---------------------------------------------------------------------
  // openai fallback: used when the rules above left gaps, or when the page
  // has json-ld in a shape we do not understand. The model sees the raw
  // json-ld documents plus the visible page text, and may only fill fields
  // that are still missing.
  // ---------------------------------------------------------------------

  function urlsInText(text, limit) {
    var m = String(text || "").match(/https?:\/\/[^\s"'\\<>]+/g) || [];
    var out = [];
    for (var i = 0; i < m.length && out.length < (limit || 16); i++) {
      var u = m[i].replace(/[),.;:]+$/, "");
      if (out.indexOf(u) === -1) {
        out.push(u);
      }
    }
    return out;
  }

  function imageLikeUrl(u) {
    return /\.(jpe?g|png|webp|avif|gif)([?#]|$)/i.test(u) || /image/i.test(u);
  }

  function missingFields(entry) {
    var missing = [];
    if (!entry.name) missing.push("name");
    if (entry.priceInPence === null || entry.priceInPence === undefined) missing.push("price");
    if (entry.packCount === null || entry.packCount === undefined) missing.push("pack_count");
    if (!(entry.categories && entry.categories.length)) missing.push("categories");
    if (!entry.image) missing.push("image");
    if (!/^(code|url-code|url-slug|edited)$/.test(entry.idSource || "")) missing.push("id");
    return missing;
  }

  function needsAiFill(entry) {
    return missingFields(entry).length > 0;
  }

  var AI_SYSTEM =
    "You extract structured product data for a mince pie price ranking site from " +
    "supermarket web pages. Pages may have no JSON-LD, or JSON-LD in unusual shapes, " +
    "so use whatever structure the input has. Reply with ONLY a JSON object, no prose:\n" +
    '{"name": string, "price_gbp": number like 2.49 or null, "pack_count": integer or null, ' +
    '"categories": [strings], "fresh": true or false, "image": string or null, "id": string or null}\n' +
    "Rules: CATEGORY_RULE " +
    "image must be one of the supplied image_urls. id: the product's numeric code if one " +
    "appears in the data, otherwise a short lowercase slug. " +
    "Prefer values found in jsonld over guesses from page_text. Use null instead of guessing.";

  // image urls the model is allowed to pick: the scored candidates plus any
  // image-ish urls that appear in the page's own json-ld. background.js passes
  // the same list to aiFill so the model cannot invent urls to download.
  function aiImageCandidates(facts, entry) {
    var candidates = (entry.imageCandidates || []).slice(0, 8).map(function (c) {
      return c.url;
    });
    if (entry.image && candidates.indexOf(entry.image) === -1) {
      candidates.unshift(entry.image);
    }
    urlsInText((facts.rawLd || []).join(" "), 16).forEach(function (u) {
      if (imageLikeUrl(u) && candidates.indexOf(u) === -1) {
        candidates.push(u);
      }
    });
    return candidates.slice(0, 16);
  }

  // The chat payload: raw json-ld first (the point is to understand shapes we
  // could not parse ourselves), then the visible text as backstop. categories
  // is the live slug list from /api/categories; without it the model is told
  // to leave categories alone, since unknown slugs abort the bulk import.
  function aiMessages(facts, entry, categories) {
    var rule = categories && categories.length
      ? "categories may only use slugs from this list, [] if none fit: " + categories.join(", ") + "."
      : "no category list is available, so categories must be [].";
    var system = AI_SYSTEM.replace("CATEGORY_RULE", rule);
    var user = JSON.stringify({
      url: entry.url || facts.url,
      title: facts.title || null,
      og_title: facts.ogTitle || null,
      og_description: facts.ogDescription || null,
      jsonld: facts.rawLd || [],
      known: {
        name: entry.name || null,
        id: entry.id,
        price_gbp: entry.priceInPence == null ? null : entry.priceInPence / 100,
        pack_count: entry.packCount == null ? null : entry.packCount,
        categories: entry.categories || []
      },
      missing: missingFields(entry),
      image_urls: aiImageCandidates(facts, entry),
      page_text: String(facts.bodyText || "").slice(0, 6000)
    });
    return [
      { role: "system", content: system },
      { role: "user", content: user }
    ];
  }

  // chat completions wrap answers in prose or ``` fences, so slice out the object
  function parseAiJson(content) {
    var text = String(content || "");
    var start = text.indexOf("{");
    var end = text.lastIndexOf("}");
    if (start === -1 || end <= start) {
      throw new Error("ai response contained no json object");
    }
    return JSON.parse(text.slice(start, end + 1));
  }

  // Only fills what is still missing, and only with values that can be checked:
  // categories must be real slugs (an empty live list means no categories can
  // be accepted), images must be urls seen on the page, and a
  // suggested id replaces a weak one (or any id when it is a real product code).
  function aiFill(entry, parsed, allowedImageUrls, categories) {
    var filled = [];
    if (!entry.name && typeof parsed.name === "string" && parsed.name.trim()) {
      entry.name = parsed.name.trim().slice(0, 200);
      filled.push("name");
    }
    if (entry.priceInPence == null) {
      var pence = priceToPence(parsed.price_gbp);
      if (pence !== null && pence >= 1 && pence <= 100000) {
        entry.priceInPence = pence;
        filled.push("priceInPence");
      }
    }
    if (entry.packCount == null) {
      var n = parseInt(parsed.pack_count, 10);
      if (isFinite(n) && n >= 1 && n <= 100) {
        entry.packCount = n;
        filled.push("packCount");
      }
    }
    if (!(entry.categories && entry.categories.length) && Array.isArray(parsed.categories)) {
      var cats = [];
      var live = categories || [];
      parsed.categories.forEach(function (c) {
        c = String(c).toLowerCase().trim();
        if (live.indexOf(c) !== -1 && cats.indexOf(c) === -1) {
          cats.push(c);
        }
      });
      if (cats.length) {
        entry.categories = cats.sort();
        filled.push("categories");
      }
    }
    if (!entry.fresh && parsed.fresh === true) {
      entry.fresh = true;
      filled.push("fresh");
    }
    if (!entry.image && parsed.image) {
      var url = String(parsed.image);
      if (/^https?:\/\//i.test(url) && (allowedImageUrls || []).indexOf(url) !== -1) {
        entry.image = url;
        entry.imageCandidates = (entry.imageCandidates || []).concat([{ url: url, score: -1 }]);
        filled.push("image");
      }
    }
    if (parsed.id !== null && parsed.id !== undefined && String(parsed.id).trim()) {
      var rawId = String(parsed.id).trim();
      var digits = rawId.replace(/\D/g, "");
      var id;
      if (/^\d{8,}$/.test(rawId) || (digits.length >= 8 && /^[0-9 -]{8,}$/.test(rawId))) {
        id = /^\d{8,}$/.test(rawId) ? rawId : digits;
      } else {
        id = slugify(rawId);
      }
      var weakId = !/^(code|url-code|edited)$/.test(entry.idSource || "");
      if (id && id !== "unknown" && (/^\d{8,}$/.test(id) || (weakId && id !== entry.id))) {
        entry.id = id;
        entry.idSource = "ai";
        filled.push("id");
      }
    }
    return { entry: entry, filled: filled };
  }

  function csvField(s) {
    s = String(s === null || s === undefined ? "" : s);
    if (/[",\n]/.test(s)) {
      return '"' + s.replace(/"/g, '""') + '"';
    }
    return s;
  }

  var CSV_HEADER = [
    "year", "makerid", "id", "displayname", "fresh", "pack_count",
    "pack_price_in_pence", "category_slugs", "web_link", "image_file"
  ];

  function buildCsv(year, rows) {
    var lines = [CSV_HEADER.join(",")];
    rows.forEach(function (r) {
      lines.push(CSV_HEADER.map(function (h) {
        return csvField(r[h]);
      }).join(","));
    });
    return lines.join("\n") + "\n";
  }

  return {
    slugify: slugify,
    parseUrl: parseUrl,
    hostOf: hostOf,
    makeridFromUrl: makeridFromUrl,
    deriveId: deriveId,
    priceToPence: priceToPence,
    packCountFromText: packCountFromText,
    guessCategories: guessCategories,
    pickProduct: pickProduct,
    scoreImages: scoreImages,
    buildEntry: buildEntry,
    buildCsv: buildCsv,
    csvField: csvField,
    missingFields: missingFields,
    needsAiFill: needsAiFill,
    aiImageCandidates: aiImageCandidates,
    aiMessages: aiMessages,
    parseAiJson: parseAiJson,
    aiFill: aiFill
  };
})();
