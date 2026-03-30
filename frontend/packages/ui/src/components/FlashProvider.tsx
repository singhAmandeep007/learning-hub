import { useState, useEffect, useCallback, useMemo, type ReactNode } from "react";
import { X, CheckCircle, AlertCircle, Info, AlertTriangle } from "lucide-react";

import { withPrefix } from "../constants";
import { FlashContext, type FlashNotificationType } from "./FlashContext";

interface Notification {
  id: string;
  message: string;
  type: FlashNotificationType;
  duration?: number;
}

interface ApiErrorShape {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

const extractErrorMessage = (error: unknown, fallbackMessage: string): string => {
  if (!error || typeof error !== "object") {
    return fallbackMessage;
  }

  const typedError = error as ApiErrorShape;
  return typedError.response?.data?.message || typedError.message || fallbackMessage;
};

const baseClass = withPrefix("flash");

export function FlashProvider({ children }: { children: ReactNode }) {
  const [notifications, setNotifications] = useState<Notification[]>([]);

  const addNotification = useCallback((message: string, type: FlashNotificationType = "info", duration = 4000) => {
    const newNotification: Notification = {
      id: Date.now().toString() + Math.random().toString(36).substring(2, 9),
      message,
      type,
      duration,
    };
    setNotifications((prev) => [...prev, newNotification]);
  }, []);

  const removeNotification = useCallback((id: string) => {
    setNotifications((prev) => prev.filter((notification) => notification.id !== id));
  }, []);

  const showSuccess = useCallback(
    (message: string, duration = 2000) => {
      addNotification(message, "success", duration);
    },
    [addNotification]
  );

  const showError = useCallback(
    (message: string, duration = 5000) => {
      addNotification(message, "error", duration);
    },
    [addNotification]
  );

  const showInfo = useCallback(
    (message: string, duration = 4000) => {
      addNotification(message, "info", duration);
    },
    [addNotification]
  );

  const showWarning = useCallback(
    (message: string, duration = 4000) => {
      addNotification(message, "warning", duration);
    },
    [addNotification]
  );

  const showQuerySuccess = useCallback(
    (message = "Data loaded successfully", duration = 2000) => {
      addNotification(message, "success", duration);
    },
    [addNotification]
  );

  const showQueryError = useCallback(
    (error: unknown, customMessage?: string, duration = 5000) => {
      const errorMessage = customMessage || extractErrorMessage(error, "Failed to load data");
      addNotification(errorMessage, "error", duration);
    },
    [addNotification]
  );

  const showMutationSuccess = useCallback(
    (message = "Operation completed successfully", duration = 4000) => {
      addNotification(message, "success", duration);
    },
    [addNotification]
  );

  const showMutationError = useCallback(
    (error: unknown, customMessage?: string, duration = 5000) => {
      const errorMessage = customMessage || extractErrorMessage(error, "Operation failed");
      addNotification(errorMessage, "error", duration);
    },
    [addNotification]
  );

  const contextValue = useMemo(
    () => ({
      addNotification,
      showSuccess,
      showError,
      showInfo,
      showWarning,
      showQuerySuccess,
      showQueryError,
      showMutationSuccess,
      showMutationError,
    }),
    [
      addNotification,
      showSuccess,
      showError,
      showInfo,
      showWarning,
      showQuerySuccess,
      showQueryError,
      showMutationSuccess,
      showMutationError,
    ]
  );

  return (
    <FlashContext.Provider value={contextValue}>
      {children}
      <div className={`${baseClass}__container`}>
        {notifications.map((notification) => (
          <FlashNotification
            key={notification.id}
            notification={notification}
            onClose={() => removeNotification(notification.id)}
          />
        ))}
      </div>
    </FlashContext.Provider>
  );
}

interface FlashNotificationProps {
  notification: Notification;
  onClose: () => void;
}

function FlashNotification({ notification, onClose }: FlashNotificationProps) {
  const { message, type, duration } = notification;

  useEffect(() => {
    if (duration) {
      const timer = setTimeout(onClose, duration);
      return () => clearTimeout(timer);
    }

    return undefined;
  }, [duration, onClose]);

  const getIcon = () => {
    switch (type) {
      case "success":
        return <CheckCircle className={`${baseClass}__icon`} />;
      case "error":
        return <AlertCircle className={`${baseClass}__icon`} />;
      case "warning":
        return <AlertTriangle className={`${baseClass}__icon`} />;
      case "info":
      default:
        return <Info className={`${baseClass}__icon`} />;
    }
  };

  return (
    <div
      className={`${baseClass}__notification ${baseClass}__notification--${type}`}
      role="alert"
      aria-live="polite"
    >
      {getIcon()}
      <p className={`${baseClass}__message`}>{message}</p>
      <button
        onClick={onClose}
        className={`${baseClass}__close`}
        aria-label="Close notification"
        type="button"
      >
        <X size={16} />
      </button>
    </div>
  );
}

export const ReactQueryFlashProvider = FlashProvider;
