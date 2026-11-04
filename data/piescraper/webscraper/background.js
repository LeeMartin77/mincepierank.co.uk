// Background event page: owns the storage entry store and the export, since
// the popup closes as soon as a download starts.

importScripts("lib/derive.js", "lib/zip.js");

var STORE_KEY = "entries";
var SETTINGS_KEY = "settings";
var REFERENCE_KEY = "reference";

// public json endpoints on the site (see api/website.yaml) are the only
// source of category slugs and brand ids; results are cached for a day and
// a scrape with neither a fresh fetch nor any cache fails outright.
var SITE_BASE = "https://mincepierank.co.uk";
var REFERENCE_TTL_MS = 24 * 60 * 60 * 1000;

browser.runtime.onMessage.addListener(function (msg) {
  switch (msg.type) {
    case "scrape":
      return scrapeActiveTab();
    case "getEntries":
      return getEntries();
    case "updateEntry":
      return updateEntry(msg.key, msg.patch);
    case "deleteEntry":
      return deleteEntry(msg.key);
    case "export":
      return exportZip(msg.year, msg.keys);
    default:
      return { ok: false, error: "unknown message type: " + msg.type };
  }
});

function scrapeActiveTab() {
  return browser.tabs.query({ active: true, currentWindow: true }).then(function (tabs) {
    if (!tabs.length) {
      throw new Error("no active tab");
    }
    var tabId = tabs[0].id;
    return browser.scripting
      .executeScript({
        target: { tabId: tabId },
        files: ["lib/derive.js", "lib/extract.js"]
      })
      .then(function () {
        return browser.scripting.executeScript({
          target: { tabId: tabId },
          func: function () {
            return globalThis.piescraperExtract();
          }
        });
      })
      .then(function (results) {
        var facts = results[0] && results[0].result;
        if (!facts || !facts.url) {
          throw new Error("extraction returned nothing (restricted page?)");
        }
        return getReference().then(function (ref) {
          if (!ref) {
            throw new Error("cannot reach " + SITE_BASE + " and there is no cached copy - check your connection");
          }
          var entry = PS.buildEntry(facts, ref);
          return mergeEntry(entry).then(function (merged) {
            return fillWithAi(facts, merged.key, merged.entry, ref.categories).then(function (done) {
              return { ok: true, key: done.key, entry: done.entry, ai: done.ai };
            });
          });
        });
      });
  });
}

// live category slugs and brand websites from the site, cached in storage
// so a scrape does not hit the network twice a day; failures fall back to
// the stale cache, and with no cache at all the scrape fails outright.
function getReference() {
  return browser.storage.local.get(REFERENCE_KEY).then(function (data) {
    var cached = data[REFERENCE_KEY];
    if (cached && Date.now() - cached.fetchedAt < REFERENCE_TTL_MS) {
      return cached;
    }
    return Promise.all([fetchJson(SITE_BASE + "/api/categories"), fetchJson(SITE_BASE + "/api/brands")])
      .then(function (lists) {
        var fresh = {
          fetchedAt: Date.now(),
          categories: (lists[0] || []).map(function (c) {
            return c.slug;
          }),
          brands: lists[1] || []
        };
        if (!fresh.categories.length || !fresh.brands.length) {
          throw new Error("site returned empty lists");
        }
        return browser.storage.local.set(objectWith(REFERENCE_KEY, fresh)).then(function () {
          return fresh;
        });
      })
      .catch(function (err) {
        console.warn("reference fetch failed, using " + (cached ? "stale cache" : "nothing") + ": " + err);
        return cached || null;
      });
  });
}

function fetchJson(url) {
  return fetch(url).then(function (res) {
    if (!res.ok) {
      throw new Error(url + " -> " + res.status);
    }
    return res.json();
  });
}

// storage.local.set wants a key/value object and the key here is a variable
function objectWith(key, value) {
  var obj = {};
  obj[key] = value;
  return obj;
}

// Last resort for pages with no or unusual json-ld: ask an openai-compatible
// chat endpoint to fill only what the extractor missed. A failing call never
// loses the entry that was already scraped and saved.
function fillWithAi(facts, key, entry, categories) {
  return readStore().then(function (store) {
    var ai = store.settings.openai || {};
    if (!ai.apiKey || !PS.needsAiFill(entry)) {
      return { key: key, entry: entry, ai: null };
    }
    return chatCompletion(ai, PS.aiMessages(facts, entry, categories))
      .then(function (parsed) {
        var allowed = PS.aiImageCandidates(facts, entry);
        var result = PS.aiFill(entry, parsed, allowed, categories);
        if (!result.filled.length) {
          return { key: key, entry: entry, ai: { filled: [] } };
        }
        return saveEntry(key, entry).then(function (saved) {
          return { key: saved.key, entry: saved.entry, ai: { filled: result.filled } };
        });
      })
      .catch(function (err) {
        return { key: key, entry: entry, ai: { error: String(err.message || err) } };
      });
  });
}

function chatCompletion(ai, messages) {
  var base = String(ai.baseUrl || "https://api.openai.com/v1").replace(/\/+$/, "");
  return fetch(base + "/chat/completions", {
    method: "POST",
    headers: {
      "content-type": "application/json",
      authorization: "Bearer " + ai.apiKey
    },
    body: JSON.stringify({
      model: ai.model || "gpt-4o-mini",
      messages: messages,
      temperature: 0
    })
  })
    .then(function (res) {
      if (!res.ok) {
        return res.text().then(function (text) {
          throw new Error("openai " + res.status + ": " + text.slice(0, 200));
        });
      }
      return res.json();
    })
    .then(function (data) {
      var choice = data.choices && data.choices[0];
      var content = choice && choice.message && choice.message.content;
      return PS.parseAiJson(content);
    });
}

function saveEntry(oldKey, entry) {
  return readStore().then(function (store) {
    var key = entryKey(entry);
    if (oldKey && oldKey !== key) {
      delete store.entries[oldKey];
    }
    store.entries[key] = entry;
    return browser.storage.local.set({ entries: store.entries }).then(function () {
      return { key: key, entry: entry };
    });
  });
}

function readStore() {
  return browser.storage.local.get([STORE_KEY, SETTINGS_KEY]).then(function (data) {
    return {
      entries: data[STORE_KEY] || {},
      settings: data[SETTINGS_KEY] || {}
    };
  });
}

function mergeEntry(fresh) {
  var key = fresh.makerid + "/" + fresh.id;
  return readStore().then(function (store) {
    var existing = store.entries[key];
    var merged = fresh;
    if (existing) {
      merged.editedFields = existing.editedFields || [];
      ["name", "priceInPence", "packCount", "categories", "fresh", "id", "makerid", "url"].forEach(function (f) {
        if (merged.editedFields.indexOf(f) !== -1) {
          merged[f] = existing[f];
        }
      });
      if (merged.editedFields.indexOf("id") !== -1) {
        merged.idSource = existing.idSource || "edited";
      }
      merged.image = pickChosenImage(existing, fresh);
      merged.imageCandidates = unionCandidates(existing.imageCandidates, fresh.imageCandidates);
      if (!merged.image && merged.imageCandidates.length) {
        merged.image = merged.imageCandidates[0].url;
      }
    }
    store.entries[entryKey(merged)] = merged;
    return browser.storage.local.set({ entries: store.entries }).then(function () {
      return { key: entryKey(merged), entry: merged };
    });
  });
}

function entryKey(entry) {
  return entry.makerid + "/" + entry.id;
}

function pickChosenImage(existing, fresh) {
  if (!existing.image) {
    return fresh.image;
  }
  var stillThere = fresh.imageCandidates.some(function (c) {
    return c.url === existing.image;
  });
  return stillThere ? existing.image : fresh.image;
}

function unionCandidates(a, b) {
  var byUrl = {};
  (a || []).concat(b || []).forEach(function (c) {
    if (!byUrl[c.url] || c.score > byUrl[c.url].score) {
      byUrl[c.url] = c;
    }
  });
  var out = Object.keys(byUrl).map(function (k) {
    return byUrl[k];
  });
  out.sort(function (x, y) {
    return y.score - x.score;
  });
  return out.slice(0, 8);
}

function getEntries() {
  return readStore().then(function (store) {
    return { ok: true, entries: store.entries, settings: store.settings };
  });
}

function updateEntry(key, patch) {
  return readStore().then(function (store) {
    var entry = store.entries[key];
    if (!entry) {
      throw new Error("no entry " + key);
    }
    entry.editedFields = entry.editedFields || [];
    Object.keys(patch).forEach(function (field) {
      entry[field] = patch[field];
      if (entry.editedFields.indexOf(field) === -1) {
        entry.editedFields.push(field);
      }
    });
    if (patch.id !== undefined) {
      entry.idSource = "edited";
    }
    var newKey = entryKey(entry);
    if (newKey !== key) {
      delete store.entries[key];
    }
    store.entries[newKey] = entry;
    return browser.storage.local.set({ entries: store.entries }).then(function () {
      return { ok: true, key: newKey, entry: entry };
    });
  });
}

function deleteEntry(key) {
  return readStore().then(function (store) {
    delete store.entries[key];
    return browser.storage.local.set({ entries: store.entries }).then(function () {
      return { ok: true };
    });
  });
}

function exportZip(year, keys) {
  if (!year || !/^\d{4}$/.test(String(year))) {
    throw new Error("year must be a four digit number");
  }
  return readStore().then(function (store) {
    var skipped = [];
    var rows = [];
    var images = [];
    var chain = Promise.resolve();

    keys.forEach(function (key) {
      var entry = store.entries[key];
      if (!entry) {
        return;
      }
      chain = chain
        .then(function () {
          if (!entry.image) {
            throw new Error("no image selected");
          }
          return fetch(entry.image, { credentials: "include" }).then(function (res) {
            if (!res.ok) {
              throw new Error("image fetch failed: " + res.status);
            }
            return res.arrayBuffer().then(function (buf) {
              var makerid = PS.slugify(entry.makerid) || "unknown";
              var id = PS.slugify(entry.id) || "unknown";
              var ext = extensionFor(res.headers.get("content-type"), entry.image);
              rows.push({
                year: String(year),
                makerid: makerid,
                id: id,
                displayname: entry.name,
                fresh: entry.fresh ? "yes" : "no",
                pack_count: entry.packCount === null || entry.packCount === undefined ? "" : entry.packCount,
                pack_price_in_pence: entry.priceInPence === null || entry.priceInPence === undefined ? "" : entry.priceInPence,
                category_slugs: (entry.categories || []).join(","),
                web_link: entry.url,
                image_file: "images/" + makerid + "/" + id + "." + ext
              });
              images.push({
                path: "images/" + makerid + "/" + id + "." + ext,
                bytes: new Uint8Array(buf)
              });
            });
          });
        })
        .catch(function (err) {
          skipped.push({ key: key, name: entry.name || key, reason: String(err.message || err) });
        });
    });

    return chain.then(function () {
      if (!rows.length) {
        throw new Error("nothing to export" + (skipped.length ? " (all entries failed)" : ""));
      }
      var zip = new ZipWriter();
      zip.add("piedata.csv", new TextEncoder().encode(PS.buildCsv(String(year), rows)));
      images.forEach(function (img) {
        zip.add(img.path, img.bytes);
      });
      return downloadZip(zip.finish(), "piedata-" + year + ".zip").then(function () {
        return { ok: true, exported: rows.length, skipped: skipped };
      });
    });
  });
}

function extensionFor(contentType, url) {
  var byType = {
    "image/jpeg": "jpg",
    "image/png": "png",
    "image/webp": "webp",
    "image/avif": "avif",
    "image/gif": "gif"
  };
  if (contentType) {
    var base = contentType.split(";")[0].trim().toLowerCase();
    if (byType[base]) {
      return byType[base];
    }
  }
  var m = String(url).match(/\.(jpe?g|png|webp|avif|gif)(\?|#|$)/i);
  if (m) {
    return m[1].toLowerCase() === "jpeg" ? "jpg" : m[1].toLowerCase();
  }
  return "jpg";
}

function downloadZip(bytes, filename) {
  var CHUNK = 0x8000;
  var binary = "";
  for (var i = 0; i < bytes.length; i += CHUNK) {
    binary += String.fromCharCode.apply(null, bytes.subarray(i, i + CHUNK));
  }
  var url = "data:application/zip;base64," + btoa(binary);
  return browser.downloads.download({ url: url, filename: filename, saveAs: true });
}
