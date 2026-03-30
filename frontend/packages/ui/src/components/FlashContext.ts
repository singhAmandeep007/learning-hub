import { createContext } from "react";

export type FlashNotificationType = "success" | "error" | "info" | "warning";

export interface FlashContextType {
  addNotification: (message: string, type?: FlashNotificationType, duration?: number) => void;
  showSuccess: (message: string, duration?: number) => void;
  showError: (message: string, duration?: number) => void;
  showInfo: (message: string, duration?: number) => void;
  showWarning: (message: string, duration?: number) => void;
  showQuerySuccess: (message?: string, duration?: number) => void;
  showQueryError: (error: unknown, customMessage?: string, duration?: number) => void;
  showMutationSuccess: (message?: string, duration?: number) => void;
  showMutationError: (error: unknown, customMessage?: string, duration?: number) => void;
}

export const FlashContext = createContext<FlashContextType | undefined>(undefined);
