/* eslint-disable @typescript-eslint/no-explicit-any */
import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  useMemo,
} from "react";
import { X, CheckCircle, AlertCircle, Info, AlertTriangle } from "lucide-react";

import "./Flash.scss";

interface Notification {
  id: string;
  message: string;
  type: "success" | "error" | "info" | "warning";
  duration?: number;
}

interface ReactQueryFlashContextType {
  addNotification: (
    message: string,
    type?: Notification["type"],
    duration?: number,
  ) => void;
  showSuccess: (message: string, duration?: number) => void;
  showError: (message: string, duration?: number) => void;
  showInfo: (message: string, duration?: number) => void;
  showWarning: (message: string, duration?: number) => void;
  // React Query specific helpers
  showQuerySuccess: (message?: string, duration?: number) => void;
  showQueryError: (
    error: any,
    customMessage?: string,
    duration?: number,
  ) => void;
  showMutationSuccess: (message?: string, duration?: number) => void;
  showMutationError: (
    error: any,
    customMessage?: string,
    duration?: number,
  ) => void;
}

const ReactQueryFlashContext = createContext<
  ReactQueryFlashContextType | undefined
>(undefined);

export const ReactQueryFlashProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [notifications, setNotifications] = useState<Notification[]>([]);

  const addNotification = useCallback(
    (
      message: string,
      type: Notification["type"] = "info",
      duration: number = 4000,
    ) => {
      const newNotification: Notification = {
        id: Date.now().toString() + Math.random().toString(36).substring(2, 9),
        message,
        type,
        duration,
      };
      setNotifications((prev) => [...prev, newNotification]);
    },
    [],
  );

  const removeNotification = useCallback((id: string) => {
    setNotifications((prev) => prev.filter((n) => n.id !== id));
  }, []);

  const showSuccess = useCallback(
    (message: string, duration = 4000) => {
      addNotification(message, "success", duration);
    },
    [addNotification],
  );

  const showError = useCallback(
    (message: string, duration = 5000) => {
      addNotification(message, "error", duration);
    },
    [addNotification],
  );

  const showInfo = useCallback(
    (message: string, duration = 4000) => {
      addNotification(message, "info", duration);
    },
    [addNotification],
  );

  const showWarning = useCallback(
    (message: string, duration = 4000) => {
      addNotification(message, "warning", duration);
    },
    [addNotification],
  );

  const showQuerySuccess = useCallback(
    (message = "Data loaded successfully", duration = 3000) => {
      addNotification(message, "success", duration);
    },
    [addNotification],
  );

  const showQueryError = useCallback(
    (error: any, customMessage?: string, duration = 5000) => {
      const errorMessage =
        customMessage ||
        error?.response?.data?.message ||
        error?.message ||
        "Failed to load data";
      addNotification(errorMessage, "error", duration);
    },
    [addNotification],
  );

  const showMutationSuccess = useCallback(
    (message = "Operation completed successfully", duration = 4000) => {
      addNotification(message, "success", duration);
    },
    [addNotification],
  );

  const showMutationError = useCallback(
    (error: any, customMessage?: string, duration = 5000) => {
      const errorMessage =
        customMessage ||
        error?.response?.data?.message ||
        error?.message ||
        "Operation failed";
      addNotification(errorMessage, "error", duration);
    },
    [addNotification],
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
    ],
  );

  return (
    <ReactQueryFlashContext.Provider value={contextValue}>
      {children}
      <div className="rq-flash-container">
        {notifications.map((notification) => (
          <ReactQueryFlashNotification
            key={notification.id}
            notification={notification}
            onClose={() => removeNotification(notification.id)}
          />
        ))}
      </div>
    </ReactQueryFlashContext.Provider>
  );
};

export const useReactQueryFlash = () => {
  const context = useContext(ReactQueryFlashContext);
  if (context === undefined) {
    throw new Error(
      "useReactQueryFlash must be used within a ReactQueryFlashProvider",
    );
  }
  return context;
};

interface ReactQueryFlashNotificationProps {
  notification: Notification;
  onClose: () => void;
}

const ReactQueryFlashNotification: React.FC<
  ReactQueryFlashNotificationProps
> = ({ notification, onClose }) => {
  const { message, type, duration } = notification;

  useEffect(() => {
    if (duration) {
      const timer = setTimeout(onClose, duration);
      return () => clearTimeout(timer);
    }
  }, [duration, onClose]);

  const getIcon = () => {
    switch (type) {
      case "success":
        return <CheckCircle className="rq-flash-icon" />;
      case "error":
        return <AlertCircle className="rq-flash-icon" />;
      case "warning":
        return <AlertTriangle className="rq-flash-icon" />;
      case "info":
      default:
        return <Info className="rq-flash-icon" />;
    }
  };

  return (
    <div
      className={`rq-flash-notification rq-flash-notification-${type}`}
      role="alert"
      aria-live="polite"
    >
      {getIcon()}
      <p className="rq-flash-message">{message}</p>
      <button
        onClick={onClose}
        className="rq-flash-close"
        aria-label="Close notification"
      >
        <X size={16} />
      </button>
    </div>
  );
};
