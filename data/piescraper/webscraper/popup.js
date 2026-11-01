// Popup: scrape button, editable entry list, export. Reads and writes the
// entry store through the background so there is one owner of the data.

var entries = {};
var settings = {};

function send(type, extra) {
  return browser.runtime.sendMessage(Object.assign({ type: type }, extra || {}));
}

function setStatus(message, cls) {
  var el = document.getElementById("status");
  el.textContent = message || "";
  el.className = cls || "";
}

function yearValue() {
  var v = parseInt(document.getElementById("year").value, 10);
  return isFinite(v) ? String(v) : "";
}

function refresh() {
  return send("getEntries").then(function (res) {
    entries = res.entries || {};
    settings = res.settings || {};
    render();
  });
}

function render() {
  document.getElementById("count").textContent = Object.keys(entries).length + " pies collected";

  var keys = Object.keys(entries).sort();
  var list = document.getElementById("entries");
  list.textContent = "";

  keys.forEach(function (key) {
    var entry = entries[key];
    var row = document.createElement("div");
    row.className = "entry" + (entry.image ? "" : " missing");

    var img = document.createElement("img");
    if (entry.image) {
      img.src = entry.image;
      img.addEventListener("error", function () {
        img.removeAttribute("src");
      });
    }
    row.appendChild(img);

    var fields = document.createElement("div");
    fields.className = "fields";
    var nameLabel = label("display name", input("text", entry.name, key, "name"));
    nameLabel.classList.add("wide");
    fields.appendChild(nameLabel);
    fields.appendChild(label("pack count", input("number", numeric(entry.packCount), key, "packCount")));
    fields.appendChild(label("price (pence)", input("number", numeric(entry.priceInPence), key, "priceInPence")));
    fields.appendChild(label("makerid", input("text", entry.makerid, key, "makerid")));
    fields.appendChild(label("id", input("text", entry.id, key, "id")));
    var cats = input("text", (entry.categories || []).join(","), key, "categories");
    cats.setAttribute("list", "categoryOptions");
    fields.appendChild(label("categories", cats));
    var fresh = document.createElement("input");
    fresh.type = "checkbox";
    fresh.checked = !!entry.fresh;
    fresh.addEventListener("change", function () {
      patch(key, { fresh: fresh.checked }, true);
    });
    var freshLabel = label("fresh", fresh);
    freshLabel.classList.add("check");
    fields.appendChild(freshLabel);
    row.appendChild(fields);

    var meta = document.createElement("div");
    meta.className = "meta";
    var keySpan = document.createElement("span");
    keySpan.className = "key";
    keySpan.textContent = key;
    meta.appendChild(keySpan);

    var link = document.createElement("a");
    link.href = entry.url;
    link.textContent = "page";
    meta.appendChild(link);

    var picker = document.createElement("select");
    var candidates = entry.imageCandidates || [];
    if (!candidates.length) {
      var option = document.createElement("option");
      option.value = "";
      option.textContent = "(no image candidates)";
      picker.appendChild(option);
    }
    candidates.forEach(function (candidate) {
      var option = document.createElement("option");
      option.value = candidate.url;
      option.textContent = candidate.url.replace(/^https?:\/\//, "").slice(0, 70) + " [" + candidate.score + "]";
      if (candidate.url === entry.image) {
        option.selected = true;
      }
      picker.appendChild(option);
    });
    if (entry.image && !(candidates.length && candidates.some(function (c) { return c.url === entry.image; }))) {
      var option = document.createElement("option");
      option.value = entry.image;
      option.textContent = entry.image.replace(/^https?:\/\//, "").slice(0, 70) + " [current]";
      option.selected = true;
      picker.appendChild(option);
    }
    picker.addEventListener("change", function () {
      patch(key, { image: picker.value }, true);
    });
    meta.appendChild(picker);

    var del = document.createElement("button");
    del.textContent = "delete";
    del.addEventListener("click", function () {
      send("deleteEntry", { key: key }).then(refresh);
    });
    meta.appendChild(del);

    row.appendChild(meta);
    list.appendChild(row);
  });
}

function label(text, control) {
  var l = document.createElement("label");
  l.textContent = text;
  l.appendChild(control);
  return l;
}

function numeric(v) {
  return v === null || v === undefined ? "" : String(v);
}

function input(kind, value, key, field) {
  var el = document.createElement("input");
  el.type = kind;
  el.value = value;
  if (kind === "number" && value === "") {
    el.classList.add("invalid");
  }
  el.addEventListener("change", function () {
    var patchValue = kind === "number" ? (el.value === "" ? null : parseInt(el.value, 10)) : el.value;
    if (kind === "number" && isNaN(patchValue)) {
      el.value = "";
      patchValue = null;
    }
    patch(key, makePatch(field, patchValue), true);
  });
  return el;
}

function makePatch(field, value) {
  var p = {};
  p[field] = value;
  return p;
}

function patch(key, fields, rerender) {
  send("updateEntry", { key: key, patch: fields }).then(function (res) {
    delete entries[key];
    entries[res.key] = res.entry;
    if (rerender) {
      render();
    }
  });
}

document.getElementById("scrape").addEventListener("click", function () {
  setStatus("scraping...");
  send("scrape").then(function (res) {
    if (!res || !res.ok) {
      setStatus("scrape failed: " + ((res && res.error) || "unknown error"), "error");
      return;
    }
    var message = "captured " + res.key;
    if (res.ai && res.ai.error) {
      message += "\nopenai fill failed: " + res.ai.error;
    } else if (res.ai && res.ai.filled && res.ai.filled.length) {
      message += "\nopenai filled: " + res.ai.filled.join(", ");
    }
    setStatus(message, res.ai && res.ai.error ? "error" : "ok");
    return refresh();
  }).catch(function (err) {
    setStatus("scrape failed: " + err.message, "error");
  });
});

document.getElementById("export").addEventListener("click", function () {
  var year = yearValue();
  if (!year) {
    setStatus("set a year first", "error");
    return;
  }
  setStatus("exporting " + Object.keys(entries).length + " pies...");
  send("export", { year: year, keys: Object.keys(entries) }).then(function (res) {
    if (!res || !res.ok) {
      setStatus("export failed: " + ((res && res.error) || "unknown error"), "error");
      return;
    }
    var message = "exported " + res.exported + " pies";
    if (res.skipped && res.skipped.length) {
      message += "\nskipped:\n" + res.skipped.map(function (s) {
        return "  " + s.name + " - " + s.reason;
      }).join("\n");
    }
    setStatus(message, res.skipped && res.skipped.length ? "error" : "ok");
  }).catch(function (err) {
    setStatus("export failed: " + err.message, "error");
  });
});

document.getElementById("year").addEventListener("change", function () {
  settings.year = yearValue();
  browser.storage.local.set({ settings: settings });
});

var AI_FIELDS = { "ai-base": "baseUrl", "ai-key": "apiKey", "ai-model": "model" };

function initAiSettings() {
  var openai = settings.openai || {};
  Object.keys(AI_FIELDS).forEach(function (id) {
    var el = document.getElementById(id);
    el.value = openai[AI_FIELDS[id]] || "";
    el.addEventListener("change", function () {
      var fresh = settings.openai || {};
      Object.keys(AI_FIELDS).forEach(function (other) {
        fresh[AI_FIELDS[other]] = document.getElementById(other).value.trim();
      });
      settings.openai = fresh;
      browser.storage.local.set({ settings: settings });
    });
  });
}

function init() {
  var datalist = document.getElementById("categoryOptions");
  PS.CATEGORIES.forEach(function (slug) {
    var option = document.createElement("option");
    option.value = slug;
    datalist.appendChild(option);
  });
  return refresh().then(function () {
    initAiSettings();
    var year = document.getElementById("year");
    year.value = settings.year || String(new Date().getFullYear());
  });
}

init();
