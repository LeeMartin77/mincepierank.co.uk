// Loads the extension's classic scripts into a vm context, the way the
// background page (importScripts) and executeScript load them.
import { readFileSync } from "node:fs";
import { createContext, runInContext } from "node:vm";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";

const here = dirname(fileURLToPath(import.meta.url));

export function loadLibs(files) {
  // Host intrinsics are passed in so objects created inside the vm (typed
  // arrays, encoders) are plain host objects the tests can assert on.
  const context = createContext({
    console,
    URL,
    Date,
    btoa: globalThis.btoa,
    TextEncoder,
    Uint8Array,
    DataView
  });
  context.globalThis = context;
  for (const file of files) {
    runInContext(readFileSync(join(here, "..", file), "utf8"), context, { filename: file });
  }
  return context;
}
