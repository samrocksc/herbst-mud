import js from "@eslint/js";
import globals from "globals";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import functional from "eslint-plugin-functional";
import tseslint from "typescript-eslint";
import { defineConfig, globalIgnores } from "eslint/config";

export default defineConfig([
  globalIgnores(["dist", "**/*.gen.ts", "**/*.gen.tsx"]),

  // 1. BASE CONFIG: Applied to all TypeScript files
  {
    files: ["**/*.{ts,tsx}"],
    extends: [
      js.configs.recommended,
      ...tseslint.configs.recommended,
      reactHooks.configs.flat.recommended,
      functional.configs.lite,
      reactRefresh.configs.vite,
    ],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
    rules: {
      "quotes": ["error", "double"],
      "semi": ["error", "always"],
      "@typescript-eslint/no-explicit-any": "error",
      "@typescript-eslint/consistent-type-imports": ["error", { prefer: "type-imports" }],
      "@typescript-eslint/no-unused-vars": ["error", { argsIgnorePattern: "^_" }],
      "@typescript-eslint/explicit-module-boundary-types": "off",

      "no-restricted-syntax": [
        "error",
        {
          "selector": "TSInterfaceDeclaration",
          "message": "Use 'type' instead of 'interface' to maintain a functional-first codebase."
        },
        {
          "selector": "ClassDeclaration",
          "message": "Prefer functional components and pure functions over classes."
        }
      ],
    },
  },

  // 2. STRICT LOGIC: Applied to pure TypeScript utility/lib files
  {
    files: ["src/lib/**/*.ts", "src/utils/**/*.ts"],
    rules: {
      "functional/no-return-void": [
        "error",
        {
          "allowNull": false,
          "allowUndefined": false,
          "ignoreInferredTypes": false
        }
      ],
      "functional/no-expression-statements": "error",
    }
  },

  // 3. REACT OVERRIDES: Relaxes rules for UI components and Hooks
  {
    files: ["**/*.tsx", "src/hooks/**/*.ts"],
    rules: {
      "functional/no-return-void": "off",
      "functional/no-expression-statements": "off",
    }
  },

  // 4. WEB-CLIENT OVERRIDES: Additional relaxations for web-client patterns
  {
    files: ["**/*.tsx"],
    rules: {
      "functional/no-mixed-types": "off",
      "functional/prefer-immutable-types": "off",
      "functional/immutable-data": "off",
      "react-hooks/set-state-in-effect": "off",
    }
  },
  {
    files: ["src/hooks/**/*.ts"],
    rules: {
      "functional/no-let": "off",
      "functional/immutable-data": "off",
      "functional/prefer-immutable-types": "off",
      "functional/no-mixed-types": "off",
    }
  },
  {
    files: ["src/lib/websocket.ts", "src/lib/theme.tsx", "src/lib/ThemeProvider.tsx"],
    rules: {
      "functional/no-expression-statements": "off",
      "functional/no-return-void": "off",
      "functional/immutable-data": "off",
      "functional/no-throw-statements": "off",
      "functional/no-mixed-types": "off",
      "functional/prefer-immutable-types": "off",
      "react-refresh/only-export-components": "off",
    }
  }
]);
