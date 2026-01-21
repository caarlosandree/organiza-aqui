import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";

const eslintConfig = defineConfig([
  ...nextVitals,
  ...nextTs,
  // Override default ignores of eslint-config-next.
  globalIgnores([
    // Default ignores of eslint-config-next:
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
  ]),
  // Disable react-hooks/incompatible-library warning for files using React Hook Form
  // This is a known limitation with React Compiler and libraries like React Hook Form
  {
    files: ["**/components/**/*Form.tsx"],
    rules: {
      "react-hooks/incompatible-library": "off",
    },
  },
]);

export default eslintConfig;
