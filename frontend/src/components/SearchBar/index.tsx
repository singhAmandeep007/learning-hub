import React, { useState } from "react";
import { Search, Tag } from "lucide-react";

import "./SearchBar.scss";

interface SearchBarProps {
  onSearch: (searchTerm: string) => void;
  onTagSearch?: (tag: string) => void;
  placeholder?: string;
}

export const SearchBar: React.FC<SearchBarProps> = ({ onSearch, onTagSearch, placeholder = "Search resources..." }) => {
  const [searchTerm, setSearchTerm] = useState("");
  const [isTagSearch, setIsTagSearch] = useState(false);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearchTerm(value);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (isTagSearch && searchTerm.trim() && onTagSearch) {
      onTagSearch(searchTerm.trim());
    } else {
      onSearch(searchTerm);
    }
  };

  const toggleTagSearch = () => {
    setIsTagSearch(!isTagSearch);
    // Clear search when switching modes
    setSearchTerm("");
    onSearch("");
  };

  // Debounced search
  React.useEffect(() => {
    if (searchTerm.trim() === "") {
      onSearch("");
      return;
    }

    const timer = setTimeout(() => {
      if (isTagSearch && onTagSearch) {
        onTagSearch(searchTerm);
      } else {
        onSearch(searchTerm);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [searchTerm, onSearch, onTagSearch, isTagSearch]);

  return (
    <form
      className="search-bar"
      onSubmit={handleSubmit}
    >
      <div className="search-input-container">
        {isTagSearch ? (
          <Tag
            className="search-icon"
            size={20}
          />
        ) : (
          <Search
            className="search-icon"
            size={20}
          />
        )}
        <input
          type="text"
          value={searchTerm}
          onChange={handleInputChange}
          placeholder={isTagSearch ? "Search by tag..." : placeholder}
          className="search-input"
          aria-label={isTagSearch ? "Search by tag" : "Search resources"}
        />
        <button
          type="button"
          className={`tag-toggle ${isTagSearch ? "active" : ""}`}
          onClick={toggleTagSearch}
          aria-pressed={isTagSearch}
          title={isTagSearch ? "Switch to normal search" : "Switch to tag search"}
        >
          <Tag size={16} />
        </button>
      </div>
    </form>
  );
};
