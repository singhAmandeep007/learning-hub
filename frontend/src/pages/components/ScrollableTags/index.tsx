import React, { useRef, useState, useEffect } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";

import { type Tag } from "../../../types";

import "./ScrollableTags.scss";

interface ScrollableTagsProps {
  tags: Tag[];
  selectedTags: string[];
  setSelectedTags: (tags: string[]) => void;
}

export const ScrollableTags: React.FC<ScrollableTagsProps> = ({ tags, selectedTags, setSelectedTags }) => {
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const [canScrollLeft, setCanScrollLeft] = useState(false);
  const [canScrollRight, setCanScrollRight] = useState(false);

  const checkScrollButtons = () => {
    if (scrollContainerRef.current) {
      const { scrollLeft, scrollWidth, clientWidth } = scrollContainerRef.current;
      setCanScrollLeft(scrollLeft > 0);
      setCanScrollRight(scrollLeft < scrollWidth - clientWidth - 1);
    }
  };

  useEffect(() => {
    checkScrollButtons();
    window.addEventListener("resize", checkScrollButtons);
    return () => window.removeEventListener("resize", checkScrollButtons);
  }, [tags]);

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

  const handleTagClick = (tagName: Tag["name"]) => {
    if (selectedTags.includes(tagName)) {
      setSelectedTags(selectedTags.filter((t) => t !== tagName));
    } else {
      setSelectedTags([...selectedTags, tagName]);
    }
  };

  return (
    <div className="scrollable-tags-container">
      {canScrollLeft && (
        <button
          className="scrollable-tags-scroll-button learning-hub-scroll-button-left"
          onClick={() => scroll("left")}
          aria-label="Scroll left"
        >
          <ChevronLeft className="scrollable-tags-scroll-icon" />
        </button>
      )}

      <div
        className="scrollable-tags-tag-filters"
        ref={scrollContainerRef}
        onScroll={checkScrollButtons}
      >
        {tags.map((tag) => (
          <button
            key={tag.name}
            onClick={() => handleTagClick(tag.name)}
            className={`scrollable-tags-tag-filter ${
              selectedTags.includes(tag.name) ? "scrollable-tags-tag-filter-active" : ""
            }`}
          >
            {tag.name}
            <span className="scrollable-tags-tag-filter-icon-usage">{tag.usageCount}</span>
          </button>
        ))}
      </div>

      {canScrollRight && (
        <button
          className="scrollable-tags-scroll-button scrollable-tags-scroll-button-right"
          onClick={() => scroll("right")}
          aria-label="Scroll right"
        >
          <ChevronRight className="scrollable-tags-scroll-icon" />
        </button>
      )}
    </div>
  );
};
