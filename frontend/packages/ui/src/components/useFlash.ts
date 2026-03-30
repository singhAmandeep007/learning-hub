import { useContext } from "react";

import { FlashContext } from "./FlashContext";

export const useFlash = () => {
  const context = useContext(FlashContext);

  if (context === undefined) {
    throw new Error("useFlash must be used within a FlashProvider");
  }

  return context;
};

export const useReactQueryFlash = useFlash;
