import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import { useMutation, useQuery, type UseMutationOptions, type UseQueryOptions } from "@tanstack/react-query";
import { useSearchParams } from "react-router";

import { useReactQueryFlash } from "../components/Flash";

import { type ResourcesFilters, type Tag, RESOURCE_TYPES } from "../types";

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

  // Track previous states to avoid duplicate notifications
  const prevStatusRef = useRef<"idle" | "pending" | "error" | "success">("idle");
  const prevErrorRef = useRef<TError | null>(null);

  const query = useQuery(queryOptions);

  // Handle success notifications
  useEffect(() => {
    if (query.isSuccess && showSuccessFlash && prevStatusRef.current !== "success") {
      const message =
        typeof successMessage === "function"
          ? successMessage(query.data)
          : successMessage || "Query completed successfully";
      flash.showQuerySuccess(message);
    }
    prevStatusRef.current = query.status;
  }, [query.isSuccess, query.status, query.data, showSuccessFlash, successMessage, flash]);

  // Handle error notifications
  useEffect(() => {
    if (query.isError && showErrorFlash && prevErrorRef.current !== query.error) {
      const message = typeof errorMessage === "function" ? errorMessage(query.error) : errorMessage;
      flash.showQueryError(query.error, message);
    }
    prevErrorRef.current = query.error;
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
      mutationOptions.onSuccess?.(data, variables, context);

      if (showSuccessFlash) {
        const message =
          typeof successMessage === "function"
            ? successMessage(data, variables)
            : successMessage || "Operation completed successfully";

        flash.showMutationSuccess(message);
      }
    },
    // Handle error notifications
    onError: (error, variables, context) => {
      mutationOptions.onError?.(error, variables, context);

      if (showErrorFlash) {
        const message = typeof errorMessage === "function" ? errorMessage(error, variables) : errorMessage;

        flash.showMutationError(error, message);
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

const ITEMS_PER_PAGE = 20;

export function useResourceFilters({
  loadedTags = [],
  hasFetchedTags,
}: {
  loadedTags?: Tag["name"][];
  hasFetchedTags: boolean;
}) {
  const [searchParams, setSearchParams] = useSearchParams();

  // Get initial values from URL
  const getUrlSearch = () => searchParams.get("search") || "";
  const getUrlTags = () => {
    const tags = searchParams.get("tags");
    return tags ? tags.split(",").filter(Boolean) : [];
  };
  const getUrlType = () => {
    const type = searchParams.get("type") as ResourcesFilters["type"];
    return type && [...Object.values(RESOURCE_TYPES), "all"].includes(type) ? type : "all";
  };

  // Search state - searchInput is what user types, activeSearch is what's searched
  const [searchInput, setSearchInput] = useState(() => getUrlSearch());
  const [activeSearch, setActiveSearch] = useState(() => getUrlSearch());

  // Filter state
  const [selectedTags, setSelectedTagsState] = useState<Tag["name"][]>([]);
  const [selectedType, setSelectedTypeState] = useState<ResourcesFilters["type"]>(() => getUrlType());
  const [currentPage, setCurrentPageState] = useState(1);

  const [hasInitializedTags, setHasInitializedTags] = useState(false);

  // Initialize selected tags from URL once tags are loaded
  useEffect(
    () => {
      if (hasFetchedTags && !hasInitializedTags) {
        const urlTags = getUrlTags();
        const validTags = urlTags.filter((tag) => loadedTags.includes(tag));
        setSelectedTagsState(validTags);
        setHasInitializedTags(true);
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [hasFetchedTags, hasInitializedTags, loadedTags]
  );

  // Update selected tags when loadedTags change (after initialization)
  // This handles cases where tags are deleted/renamed after resources are updated
  useEffect(() => {
    if (hasFetchedTags && hasInitializedTags) {
      setSelectedTagsState((prevTags) => {
        return prevTags.filter((tag) => loadedTags.includes(tag));
      });
    }
  }, [loadedTags, hasFetchedTags, hasInitializedTags]);

  // Calculate cursor for pagination
  const calculateCursor = useCallback((page: number) => {
    return page <= 1 ? null : (page - 1) * ITEMS_PER_PAGE;
  }, []);

  // Sync URL with state changes - only when activeSearch, filters, or page change
  const updateUrlParams = useCallback(() => {
    const params = new URLSearchParams();

    if (activeSearch) params.set("search", activeSearch);
    if (selectedType && selectedType !== "all") params.set("type", selectedType);
    if (selectedTags.length > 0) params.set("tags", selectedTags.join(","));

    setSearchParams(params, { replace: true });
  }, [activeSearch, selectedTags, selectedType, setSearchParams]);

  useEffect(
    () => {
      // Only sync to URL after tags have been initialized to prevent clearing URL tags
      if (hasInitializedTags || !getUrlTags().length) {
        updateUrlParams();
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [updateUrlParams, hasInitializedTags]
  );

  // Reset Page when filters change
  useEffect(() => {
    if (currentPage !== 1) {
      setCurrentPageState(1);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeSearch, selectedTags, selectedType]);

  // Query params for API calls
  const queryParams = useMemo(() => {
    const calculatedCursor = calculateCursor(currentPage);

    return {
      ...(activeSearch ? { search: activeSearch } : {}),
      ...(selectedType && selectedType !== "all" ? { type: selectedType } : {}),
      ...(selectedTags.length > 0 ? { tags: selectedTags } : {}),
      ...(calculatedCursor !== null ? { cursor: String(calculatedCursor), limit: String(ITEMS_PER_PAGE) } : {}),
    };
  }, [activeSearch, selectedType, selectedTags, currentPage, calculateCursor]);

  // Check if any filters are active
  const hasActiveFilters = useMemo(() => {
    return Boolean(activeSearch || selectedType !== "all" || selectedTags.length > 0);
  }, [activeSearch, selectedType, selectedTags]);

  // Handle clear filters
  const handleClearFilters = useCallback(() => {
    setSearchInput("");
    setActiveSearch("");
    setSelectedTagsState([]);
    setSelectedTypeState("all");
    setCurrentPageState(1);
  }, []);

  // Wrapper functions to ensure proper state updates
  const setSelectedTags = useCallback((tags: Tag["name"][]) => {
    setSelectedTagsState(tags);
  }, []);

  const setSelectedType = useCallback((type: ResourcesFilters["type"]) => {
    setSelectedTypeState(type);
  }, []);

  const setCurrentPage = useCallback((page: number) => {
    setCurrentPageState(page);
  }, []);

  return {
    searchInput,
    setSearchInput,

    selectedTags,
    selectedType,
    currentPage,

    setSelectedTags,
    setSelectedType,
    setCurrentPage,
    setActiveSearch,
    handleClearFilters,

    updateUrlParams,

    queryParams,
    hasActiveFilters,
  };
}
