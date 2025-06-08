import { useCallback, useEffect, useRef, useState } from "react";

import { useMutation, useQuery, type UseMutationOptions, type UseQueryOptions } from "@tanstack/react-query";
import { useReactQueryFlash } from "../components/Flash";

export function useQueryWithFlash<TData, TError = Error>(
  options: UseQueryOptions<TData, TError> & {
    successMessage?: string | ((data: TData) => string);
    errorMessage?: string | ((error: TError) => string);
    showSuccessFlash?: boolean;
    showErrorFlash?: boolean;
  }
) {
  const flash = useReactQueryFlash();
  const { successMessage, errorMessage, showSuccessFlash = false, showErrorFlash = true, ...queryOptions } = options;

  const query = useQuery(queryOptions);

  // Handle success notifications
  useEffect(() => {
    if (query.isSuccess && showSuccessFlash && query.data) {
      const message =
        typeof successMessage === "function"
          ? successMessage(query.data)
          : successMessage || "Data loaded successfully";
      flash.showQuerySuccess(message);
    }
  }, [query.isSuccess, query.data, showSuccessFlash, successMessage, flash]);

  // Handle error notifications
  useEffect(() => {
    if (query.isError && showErrorFlash && query.error) {
      const message = typeof errorMessage === "function" ? errorMessage(query.error) : errorMessage;
      flash.showQueryError(query.error, message);
    }
  }, [query.isError, query.error, showErrorFlash, errorMessage, flash]);

  return query;
}

// Custom hook for mutations with automatic flash notifications
export function useMutationWithFlash<TData, TError = Error, TVariables = void, TContext = unknown>(
  options: UseMutationOptions<TData, TError, TVariables, TContext> & {
    successMessage?: string | ((data: TData, variables: TVariables) => string);
    errorMessage?: string | ((error: TError, variables: TVariables) => string);
    showSuccessFlash?: boolean;
    showErrorFlash?: boolean;
  }
) {
  const flash = useReactQueryFlash();
  const { successMessage, errorMessage, showSuccessFlash = true, showErrorFlash = true, ...mutationOptions } = options;

  const mutation = useMutation({
    ...mutationOptions,
    // Handle success notifications
    onSuccess: (data, variables, context) => {
      if (showSuccessFlash) {
        const message =
          typeof successMessage === "function"
            ? successMessage(data, mutation.variables!)
            : successMessage || "Operation completed successfully";

        flash.showMutationSuccess(message);
      }
      if (mutationOptions.onSuccess) {
        mutationOptions.onSuccess(data, variables, context);
      }
    },
    // Handle error notifications
    onError: (error, variables, context) => {
      if (showErrorFlash) {
        const message = typeof errorMessage === "function" ? errorMessage(error, mutation.variables!) : errorMessage;

        flash.showMutationError(mutation.error, message);
      }
      if (mutationOptions.onError) {
        mutationOptions.onError(error, variables, context);
      }
    },
  });

  return mutation;
}

// Returns the previous value of the given variable.
export function usePrevious<T>(value: T, initialValue?: T): T | undefined {
  const ref = useRef<T | undefined>(initialValue);

  useEffect(() => {
    ref.current = value;
  }, [value]);

  return ref.current;
}

export function useDebouncedInputState<T extends HTMLInputElement | HTMLTextAreaElement>(
  initialValue: string = "",
  delay: number = 500
): [string, (event: React.ChangeEvent<T>) => void, React.RefObject<T | undefined>] {
  const inputRef = useRef<T | undefined>(undefined);

  // immediate value of the input as the user types.
  const [immediateValue, setImmediateValue] = useState<string>(initialValue);

  // debounced value, which updates after the delay.
  const [debouncedValue, setDebouncedValue] = useState<string>(initialValue);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(immediateValue);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [immediateValue, delay]);

  const handleChange = useCallback((event: React.ChangeEvent<T>) => {
    // Update the immediate value state whenever the input value changes.
    setImmediateValue(event.target.value);
  }, []);

  return [debouncedValue, handleChange, inputRef];
}
