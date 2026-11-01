import test from "node:test";
import assert from "node:assert/strict";
import { crc32 as zlibCrc32 } from "node:zlib";
import { loadLibs } from "./helpers.mjs";

const ZipWriter = loadLibs(["lib/zip.js"]).ZipWriter;

// Independent reader (signature scan instead of relying on the EOCD) so the
// writer's offsets end up cross-checked rather than self-consistent.
function readZip(buf) {
  const files = [];
  for (let i = 0; i + 4 <= buf.length; i++) {
    if (buf.readUInt32LE(i) !== 0x04034b50) {
      continue;
    }
    const method = buf.readUInt16LE(i + 8);
    const crc = buf.readUInt32LE(i + 14);
    const compSize = buf.readUInt32LE(i + 18);
    const uncompSize = buf.readUInt32LE(i + 22);
    const nameLen = buf.readUInt16LE(i + 26);
    const extraLen = buf.readUInt16LE(i + 28);
    const name = buf.subarray(i + 30, i + 30 + nameLen).toString("utf8");
    const data = buf.subarray(i + 30 + nameLen + extraLen, i + 30 + nameLen + extraLen + compSize);
    files.push({ name, method, crc, compSize, uncompSize, data });
    i = i + 30 + nameLen + extraLen + compSize - 1;
  }
  return files;
}

test("crc32 matches node zlib", function () {
  const bytes = new TextEncoder().encode("mince pies");
  assert.equal(ZipWriter.crc32(bytes), zlibCrc32(bytes) >>> 0);
});

test("zip round trip with independent reader", function () {
  const zip = new ZipWriter();
  const csv = new TextEncoder().encode("year,makerid,id\n2026,greggs,mince-pie\n");
  const fakeJpeg = Uint8Array.from([0xff, 0xd8, 0xff, 0xe0, 1, 2, 3, 250]);
  zip.add("piedata.csv", csv);
  zip.add("images/greggs/mince-pie.jpg", fakeJpeg);
  const out = Buffer.from(zip.finish());

  assert.equal(out.readUInt32LE(out.length - 22), 0x06054b50, "EOCD at end");
  assert.equal(out.readUInt16LE(out.length - 22 + 10), 2, "entry count");

  const files = readZip(out);
  assert.equal(files.length, 2);
  assert.deepEqual(
    files.map((f) => f.name),
    ["piedata.csv", "images/greggs/mince-pie.jpg"]
  );
  for (const f of files) {
    assert.equal(f.method, 0, "stored");
    assert.equal(f.compSize, f.uncompSize);
    const expected = f.name.endsWith(".jpg") ? fakeJpeg : csv;
    assert.equal(Buffer.compare(Buffer.from(f.data), Buffer.from(expected)), 0);
    assert.equal(f.crc, zlibCrc32(f.data) >>> 0);
  }

  const centralCount = out.readUInt16LE(out.length - 22 + 8);
  assert.equal(centralCount, 2);
});

test("rejects non-Uint8Array data", function () {
  const zip = new ZipWriter();
  assert.throws(() => zip.add("x", "string"), (err) => err.name === "TypeError");
});
