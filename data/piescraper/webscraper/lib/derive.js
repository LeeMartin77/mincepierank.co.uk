// Pure scraping logic, shared by the page extractor, the background page
// (importScripts) and the node tests. Attaches to the global scope so it
// can also be loaded as a classic script via scripting.executeScript.

var PS = (function () {
  "use strict";

  var CATEGORIES = [
    "alcohol", "almond", "alternative", "apple", "bakery", "brown-butter",
    "caramel", "chocolate", "cognac", "classic", "crumble", "flapjack",
    "frangipane", "fresh", "gingerbread", "giant", "gluten-free", "iced",
    "lattice", "maple", "mini", "orange", "organic", "pecan", "premium",
    "puff-pastry", "rum", "salted-caramel", "shortcrust", "spiced",
    "standard", "toffee", "vegan"
  ];

  var MAKERID_BY_HOST = {
    "aldi.co.uk": "aldi",
    "asda.com": "asda",
    "coop.co.uk": "coop",
    "farmfoods.co.uk": "farmfoods",
    "greggs.co.uk": "greggs",
    "iceland.co.uk": "iceland",
    "icelandfoods.co.uk": "iceland",
    "lidl.co.uk": "lidl",
    "marksandspencer.com": "marks-and-spencer",
    "morrisons.com": "morrisons",
    "sainsburys.co.uk": "sainsburys",
    "tesco.com": "tesco",
    "waitrose.com": "waitrose",
    "waitrosecellar.com": "waitrose-cellar"
  };

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

  function makeridFromUrl(url) {
    var u = parseUrl(url);
    if (!u) {
      return "unknown";
    }
    var host = u.hostname.toLowerCase().replace(/^www\./, "");
    if (MAKERID_BY_HOST[host]) {
      return MAKERID_BY_HOST[host];
    }
    if (host.startsWith("groceries.")) {
      host = host.slice("groceries.".length);
    }
    return slugify(host.replace(/\.(co\.uk|com|uk|ie|net)$/g, "")) || "unknown";
  }

  // id is what converges on re-scrape (it is part of the primary key in the
  // bulk csv), so prefer stable numeric product codes over anything else.
  function deriveId(facts) {
    var code = firstDigits(facts.gtin) || firstDigits(facts.sku);
    if (code && code.length >= 8) {
      return code;
    }
    var u = parseUrl(facts.url);
    if (u) {
      var segs = u.pathname.split("/").filter(Boolean);
      for (var i = segs.length - 1; i >= 0; i--) {
        var num = firstDigits(decodeURIComponent(segs[i]));
        if (num && num.length >= 8) {
          return num;
        }
      }
      for (var j = segs.length - 1; j >= 0; j--) {
        var seg = slugify(decodeURIComponent(segs[j]));
        if (seg && !/^(p|pd|pdp|product|products|item|sku|c)$/g.test(seg) && !/^\d+$/.test(seg)) {
          return seg;
        }
      }
    }
    return slugify(facts.name) || "unknown";
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

  function buildEntry(facts) {
    var product = pickProduct(facts.products || []);
    var name = (product && product.name) || facts.ogTitle || facts.title || "";
    var id = deriveId({
      url: facts.url,
      name: name,
      gtin: product && product.gtin,
      sku: product && product.sku
    });
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
      makerid: makeridFromUrl(facts.url),
      id: id,
      url: facts.url,
      name: name,
      description: (product && product.description) || facts.ogDescription || "",
      priceInPence: price,
      packCount: packCountFromText(text),
      fresh: false,
      categories: guessCategories(text),
      image: candidates.length ? candidates[0].url : (facts.ogImage || ""),
      imageCandidates: candidates,
      capturedAt: new Date().toISOString(),
      editedFields: []
    };
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
    CATEGORIES: CATEGORIES,
    MAKERID_BY_HOST: MAKERID_BY_HOST,
    slugify: slugify,
    parseUrl: parseUrl,
    makeridFromUrl: makeridFromUrl,
    deriveId: deriveId,
    priceToPence: priceToPence,
    packCountFromText: packCountFromText,
    guessCategories: guessCategories,
    pickProduct: pickProduct,
    scoreImages: scoreImages,
    buildEntry: buildEntry,
    buildCsv: buildCsv,
    csvField: csvField
  };
})();
