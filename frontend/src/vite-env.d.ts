/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_VALID_PRODUCTS: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
