import { useContext } from "react";

import { ReactQueryFlashContext } from "./flashContext";

export const useReactQueryFlash = () => {
  const context = useContext(ReactQueryFlashContext);

  if (context === undefined) {
    throw new Error("useReactQueryFlash must be used within a ReactQueryFlashProvider");
  }

  return context;
};
