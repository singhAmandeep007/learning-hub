import { useRef, useState, useEffect, useCallback } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";

import { withPrefix } from "../constants";

export interface ScrollableTagItem {
  id: string;
  label: string;
  count?: number;
}

export interface ScrollableTagsProps {
  items: ScrollableTagItem[];
  selectedItemIds: string[];
  onSelectedItemIdsChange: (ids: string[]) => void;
  ariaLabel?: string;
}

const baseClass = withPrefix("scrollable-tags");

export function ScrollableTags({
  items,
  selectedItemIds,
  onSelectedItemIdsChange,
  ariaLabel = "Tag filters",
}: ScrollableTagsProps) {
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const [canScrollLeft, setCanScrollLeft] = useState(false);
  const [canScrollRight, setCanScrollRight] = useState(false);

  const checkScrollButtons = useCallback(() => {
    if (scrollContainerRef.current) {
      const { scrollLeft, scrollWidth, clientWidth } = scrollContainerRef.current;
      setCanScrollLeft(scrollLeft > 0);
      setCanScrollRight(scrollLeft < scrollWidth - clientWidth - 1);
    }
  }, []);

  useEffect(() => {
    checkScrollButtons();
    window.addEventListener("resize", checkScrollButtons);
    return () => window.removeEventListener("resize", checkScrollButtons);
  }, [items, checkScrollButtons]);

  const scroll = (direction: "left" | "right") => {
    if (scrollContainerRef.current) {
      const scrollAmount = 200;
      const newScrollLeft =
        direction === "left"
          ? scrollContainerRef.current.scrollLeft - scrollAmount
          : scrollContainerRef.current.scrollLeft + scrollAmount;

      scrollContainerRef.current.scrollTo({
        left: newScrollLeft,
        behavior: "smooth",
      });
    }
  };

  const handleItemClick = (itemId: string) => {
    if (selectedItemIds.includes(itemId)) {
      onSelectedItemIdsChange(selectedItemIds.filter((id) => id !== itemId));
      return;
    }

    onSelectedItemIdsChange([...selectedItemIds, itemId]);
  };

  return (
    <div className={`${baseClass}__container`}>
      {canScrollLeft && (
        <button
          className={`${baseClass}__scroll-button ${baseClass}__scroll-button--left`}
          onClick={() => scroll("left")}
          aria-label="Scroll left"
          type="button"
        >
          <ChevronLeft className={`${baseClass}__scroll-icon`} />
        </button>
      )}

      <div
        className={`${baseClass}__items`}
        ref={scrollContainerRef}
        onScroll={checkScrollButtons}
        role="toolbar"
        aria-label={ariaLabel}
      >
        {items.map((item) => {
          const isSelected = selectedItemIds.includes(item.id);

          return (
            <button
              key={item.id}
              onClick={() => handleItemClick(item.id)}
              className={`${baseClass}__item ${isSelected ? `${baseClass}__item--active` : ""}`}
              aria-pressed={isSelected}
              type="button"
            >
              {item.label}
              {typeof item.count === "number" && <span className={`${baseClass}__item-count`}>{item.count}</span>}
            </button>
          );
        })}
      </div>

      {canScrollRight && (
        <button
          className={`${baseClass}__scroll-button ${baseClass}__scroll-button--right`}
          onClick={() => scroll("right")}
          aria-label="Scroll right"
          type="button"
        >
          <ChevronRight className={`${baseClass}__scroll-icon`} />
        </button>
      )}
    </div>
  );
}
