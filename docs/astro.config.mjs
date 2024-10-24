import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";

import sitemap from "@astrojs/sitemap";

// https://astro.build/config
export default defineConfig({
  site: "https://armintalaie.github.io",
  base: "/docs",
  integrations: [
    sitemap(),
    starlight({
      title: "",
      head: [
        {
          tag: "script",
          attrs: {
            src: "/clarity.js",
            defer: true,
          },
        },
      ],
      logo: {
        light: "./src/assets/cli.svg",
        dark: "./src/assets/dark-logo.svg",
      },
      social: {
        github: "https://github.com/armintalaie/vigilant",
      },
      sidebar: [],
    }),
  ],
});
