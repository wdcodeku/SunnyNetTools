import { resolve } from "node:path";
import { rcedit } from "rcedit";

const [exePathArg, iconPathArg, manifestPathArg] = process.argv.slice(2);
const version = process.env.APP_VERSION;

if (!exePathArg || !iconPathArg || !manifestPathArg || !version) {
  throw new Error(
    "Usage: APP_VERSION=x.y.z node set-windows-resources.mjs <exe> <icon> <manifest>",
  );
}

await rcedit(resolve(exePathArg), {
  "file-version": version,
  "product-version": version,
  "version-string": {
    CompanyName: "qtgolang",
    FileDescription: "SunnyNetTools",
    FileVersion: `v${version}`,
    LegalCopyright: "© 2026, 秦天",
    OriginalFilename: "SunnyNetTools.exe",
    ProductName: "SunnyNetTools",
    ProductVersion: `v${version}`,
  },
  icon: resolve(iconPathArg),
  "application-manifest": resolve(manifestPathArg),
});
