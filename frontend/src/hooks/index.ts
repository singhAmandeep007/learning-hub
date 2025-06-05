import { useMutation, useQuery, type UseMutationOptions, type UseQueryOptions } from "@tanstack/react-query";
import { useReactQueryFlash } from "../components/Flash";
import { useEffect } from "react";

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

  const mutation = useMutation(mutationOptions);

  // Handle success notifications
  useEffect(() => {
    if (mutation.isSuccess && showSuccessFlash && mutation.data) {
      const message =
        typeof successMessage === "function"
          ? successMessage(mutation.data, mutation.variables!)
          : successMessage || "Operation completed successfully";
      flash.showMutationSuccess(message);
    }
  }, [mutation.isSuccess, mutation.data, mutation.variables, showSuccessFlash, successMessage, flash]);

  // Handle error notifications
  useEffect(() => {
    if (mutation.isError && showErrorFlash && mutation.error) {
      const message =
        typeof errorMessage === "function" ? errorMessage(mutation.error, mutation.variables!) : errorMessage;
      flash.showMutationError(mutation.error, message);
    }
  }, [mutation.isError, mutation.error, mutation.variables, showErrorFlash, errorMessage, flash]);

  return mutation;
}
