// Minimal zip writer (STORE, no compression) so the extension can produce
// the bulk upload archive with no bundled dependencies. Product images are
// already compressed formats, so deflating them buys nothing; the server
// reads the archive with Go's archive/zip.

var ZipWriter = (function () {
  "use strict";

  var CRC_TABLE = (function () {
    var table = new Uint32Array(256);
    for (var n = 0; n < 256; n++) {
      var c = n;
      for (var k = 0; k < 8; k++) {
        c = c & 1 ? 0xedb88320 ^ (c >>> 1) : c >>> 1;
      }
      table[n] = c >>> 0;
    }
    return table;
  })();

  function crc32(bytes) {
    var c = 0xffffffff;
    for (var i = 0; i < bytes.length; i++) {
      c = CRC_TABLE[(c ^ bytes[i]) & 0xff] ^ (c >>> 8);
    }
    return (c ^ 0xffffffff) >>> 0;
  }

  function dosDateTime(d) {
    var year = d.getFullYear();
    if (year < 1980) {
      return { time: 0, date: (1 << 5) | 1 };
    }
    var time = (d.getHours() << 11) | (d.getMinutes() << 5) | Math.floor(d.getSeconds() / 2);
    var date = ((year - 1980) << 9) | ((d.getMonth() + 1) << 5) | d.getDate();
    return { time: time, date: date };
  }

  function ZipWriter() {
    this.entries = [];
  }

  ZipWriter.prototype.add = function (name, data) {
    if (!(data instanceof Uint8Array)) {
      throw new TypeError("zip data must be a Uint8Array");
    }
    var nameBytes = new TextEncoder().encode(name);
    if (nameBytes.length > 0xffff) {
      throw new Error("zip entry name too long: " + name);
    }
    this.entries.push({ name: nameBytes, data: data, crc: crc32(data) });
  };

  ZipWriter.prototype.finish = function () {
    var now = dosDateTime(new Date());
    var localSize = 0;
    var centralSize = 0;
    var self = this;

    this.entries.forEach(function (e) {
      localSize += 30 + e.name.length + e.data.length;
      centralSize += 46 + e.name.length;
    });

    var out = new Uint8Array(localSize + centralSize + 22);
    var view = new DataView(out.buffer);
    var offset = 0;
    var centralOffset = 0;

    function put16(v) { view.setUint16(offset, v, true); offset += 2; }
    function put32(v) { view.setUint32(offset, v, true); offset += 4; }
    function putBytes(b) { out.set(b, offset); offset += b.length; }

    this.entries.forEach(function (e) {
      e.localOffset = offset;
      put32(0x04034b50); // local file header magic
      put16(20);         // version needed
      put16(0x0800);     // flags: UTF-8 names
      put16(0);          // method: store
      put16(now.time);
      put16(now.date);
      put32(e.crc);
      put32(e.data.length);
      put32(e.data.length);
      put16(e.name.length);
      put16(0);
      putBytes(e.name);
      putBytes(e.data);
    });

    centralOffset = offset;
    self.entries.forEach(function (e) {
      put32(0x02014b50); // central directory entry magic
      put16(20);         // version made by
      put16(20);         // version needed
      put16(0x0800);     // flags: UTF-8 names
      put16(0);          // method: store
      put16(now.time);
      put16(now.date);
      put32(e.crc);
      put32(e.data.length);
      put32(e.data.length);
      put16(e.name.length);
      put16(0);          // extra length
      put16(0);          // comment length
      put16(0);          // disk number start
      put16(0);          // internal attributes
      put32(0);          // external attributes
      put32(e.localOffset);
      putBytes(e.name);
    });

    put32(0x06054b50); // end of central directory magic
    put16(0);
    put16(0);
    put16(self.entries.length);
    put16(self.entries.length);
    put32(centralSize);
    put32(centralOffset);
    put16(0);

    return out;
  };

  ZipWriter.crc32 = crc32;
  return ZipWriter;
})();
