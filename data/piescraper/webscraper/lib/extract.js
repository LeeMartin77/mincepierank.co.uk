// Runs inside the page (injected by background.js after lib/derive.js).
// Only gathers raw facts from the DOM; all the decision-making lives in
// PS (derive.js) where it can be unit tested.

globalThis.piescraperExtract = function () {
  "use strict";

  function metaContent(names) {
    for (var i = 0; i < names.length; i++) {
      var el =
        document.querySelector('meta[property="' + names[i] + '"]') ||
        document.querySelector('meta[name="' + names[i] + '"]');
      if (el && el.content) {
        return el.content.trim();
      }
    }
    return null;
  }

  // nodes: every object anywhere in the json-ld, for normalizeProducts.
  // rawDocs: the whole documents as text, so the openai fallback can make
  // sense of shapes normalizeProducts does not understand.
  function collectLd() {
    var nodes = [];
    var rawDocs = [];
    var rawSize = 0;
    document.querySelectorAll('script[type="application/ld+json"]').forEach(function (script) {
      var data;
      try {
        data = JSON.parse(script.textContent);
      } catch (e) {
        return;
      }
      if (rawDocs.length < 10) {
        var text = JSON.stringify(data);
        if (text && rawSize < 8000) {
          rawDocs.push(text.slice(0, 8000 - rawSize));
          rawSize += Math.min(text.length, 8000);
        }
      }
      walk(data);
    });
    function walk(node) {
      if (Array.isArray(node)) {
        node.forEach(walk);
        return;
      }
      if (node && typeof node === "object") {
        nodes.push(node);
        walk(node["@graph"]);
      }
    }
    return { nodes: nodes, rawDocs: rawDocs };
  }

  function pageText() {
    var text = document.body ? document.body.innerText || "" : "";
    return text.replace(/\s+/g, " ").trim().slice(0, 8000);
  }

  function asText(v) {
    if (typeof v === "string") {
      return v;
    }
    if (v && typeof v === "object") {
      var keys = Object.keys(v);
      if (keys.length && typeof v[keys[0]] === "string") {
        return v[keys[0]];
      }
    }
    return null;
  }

  function imageUrls(image) {
    var urls = [];
    function one(v) {
      if (typeof v === "string") {
        urls.push(v);
      } else if (Array.isArray(v)) {
        v.forEach(one);
      } else if (v && typeof v === "object" && typeof v.url === "string") {
        urls.push(v.url);
      }
    }
    one(image);
    return urls;
  }

  function normalizeProducts(ldNodes) {
    var products = [];
    ldNodes.forEach(function (node) {
      var type = node["@type"];
      var isProduct =
        (typeof type === "string" && type.toLowerCase() === "product") ||
        (Array.isArray(type) && type.some(function (t) { return String(t).toLowerCase() === "product"; }));
      if (!isProduct) {
        return;
      }
      var offers = node.offers;
      if (Array.isArray(offers)) {
        offers = offers[0];
      }
      var price = null;
      if (offers && typeof offers === "object") {
        price = offers.price !== undefined ? offers.price :
          offers.lowPrice !== undefined ? offers.lowPrice : null;
      }
      products.push({
        name: asText(node.name),
        description: asText(node.description),
        imageUrls: imageUrls(node.image),
        price: price,
        gtin: node.gtin13 || node.gtin12 || node.gtin || node.gtin8 || node.ean || null,
        sku: node.sku || node.mpn || null
      });
    });
    return products;
  }

  function bestSrcset(srcset) {
    var best = null;
    var bestSize = 0;
    String(srcset || "").split(",").forEach(function (part) {
      var bits = part.trim().split(/\s+/);
      if (bits.length < 2) {
        return;
      }
      var size = parseInt(bits[1], 10);
      if (isFinite(size) && size > bestSize) {
        bestSize = size;
        best = bits[0];
      }
    });
    return best;
  }

  function collectImages() {
    var out = [];
    document.querySelectorAll("img").forEach(function (img) {
      var url = (img.currentSrc && img.currentSrc.indexOf("data:") !== 0 ? img.currentSrc : img.src) || "";
      var w = img.naturalWidth || img.width || 0;
      var h = img.naturalHeight || img.height || 0;
      var fromSrcset = bestSrcset(img.getAttribute("srcset"));
      if (fromSrcset) {
        url = fromSrcset.indexOf("data:") !== 0 && fromSrcset.indexOf("http") === 0 ? fromSrcset : url;
      }
      if (!url || url.indexOf("data:") === 0 || url.indexOf("http") !== 0) {
        return;
      }
      if (!img.complete && !fromSrcset && w < 150 && h < 150) {
        return;
      }
      var productish = !!(
        img.closest('[class*="product" i], [id*="product" i], [itemprop*="image" i], [class*="gallery" i], [class*="media" i], figure, main')
      );
      out.push({ url: url, alt: img.alt || "", w: w, h: h, productish: productish });
    });
    return out.slice(0, 200);
  }

  var priceMeta =
    (document.querySelector('meta[itemprop="price"]') || {}).content ||
    metaContent(["product:price:amount", "og:price:amount"]);

  var ld = collectLd();

  return {
    url: window.location.href,
    title: document.title,
    ogTitle: metaContent(["og:title"]),
    ogDescription: metaContent(["og:description"]),
    ogImage: metaContent(["og:image:secure_url", "og:image"]),
    priceMeta: priceMeta || null,
    products: normalizeProducts(ld.nodes),
    rawLd: ld.rawDocs,
    bodyText: pageText(),
    images: collectImages()
  };
};
