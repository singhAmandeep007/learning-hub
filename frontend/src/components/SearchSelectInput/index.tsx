import React, { useState, useRef, useEffect, useCallback } from "react";

import { X } from "lucide-react";

import "./SearchSelectInput.scss";

export interface Item {
  id: string;
  name: string;
  isNew?: boolean; // flag to identify newly added tags
}

interface SearchSelectInputProps {
  items: Item[];
  onSelectedItemsChange: (selectedItems: Item[]) => void;
  initialSelectedItems?: Item[];
  placeholder?: string;
  allowNewTags?: boolean; // allow adding new tags that are not in the list
}

export const SearchSelectInput: React.FC<SearchSelectInputProps> = ({
  items,
  onSelectedItemsChange,
  initialSelectedItems = [],
  placeholder = "Search and select items...",
  allowNewTags = false, // default to false
}) => {
  const [inputValue, setInputValue] = useState<string>("");
  const [selectedItems, setSelectedItems] = useState<Item[]>(initialSelectedItems);
  const [filteredItems, setFilteredItems] = useState<Item[]>([]);
  const [isDropdownOpen, setIsDropdownOpen] = useState<boolean>(false);

  const wrapperRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const getNewTagId = useCallback(
    (tagName: string) => `new-${tagName.toLowerCase().replace(/\s+/g, "-")}-${Date.now()}`,
    []
  );

  // filter items based on input value
  useEffect(() => {
    const lowerCaseInput = inputValue.toLowerCase().trim();
    let newFilteredItems: Item[] = [];

    if (inputValue.length > 0) {
      // filter existing items based on search input
      newFilteredItems = items.filter(
        (item) =>
          item.name.toLowerCase().includes(lowerCaseInput) && !selectedItems.some((selected) => selected.id === item.id)
      );

      // check if input value exactly matches any existing (non-selected) item
      const exactMatchExists = items.some(
        (item) =>
          item.name.toLowerCase() === lowerCaseInput && !selectedItems.some((selected) => selected.id === item.id)
      );
      // Check if input value exactly matches any already selected item (to avoid duplicate 'add new' option)
      const alreadySelectedExactMatch = selectedItems.some((item) => item.name.toLowerCase() === lowerCaseInput);

      // if allowNewTags is true AND no exact match found among all items (available or selected)
      // AND the input is not just whitespace
      if (
        allowNewTags &&
        inputValue.trim().length > 0 && // Ensure not just whitespace
        !exactMatchExists &&
        !alreadySelectedExactMatch
      ) {
        // add a "create new" option if no existing item matches exactly
        // and ensure it's not a duplicate of an already selected 'new' tag
        const isNewTagAlreadySelected = selectedItems.some(
          (item) => item.isNew && item.name.toLowerCase() === lowerCaseInput
        );

        if (!isNewTagAlreadySelected) {
          newFilteredItems.unshift({
            id: getNewTagId(inputValue.trim()), // use a temporary ID
            name: `Add "${inputValue.trim()}"`, // display text for adding new tag
            isNew: true, // mark it as a new tag option
          });
        }
      }
    } else {
      // Show all unselected items when input is empty
      newFilteredItems = items.filter((item) => !selectedItems.some((selected) => selected.id === item.id));
    }

    setFilteredItems(newFilteredItems);
  }, [inputValue, items, selectedItems, allowNewTags, getNewTagId]);

  // Handle clicks outside the component to close the dropdown
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target as Node)) {
        setIsDropdownOpen(false);
        setInputValue(""); // Clear input when clicking outside
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
  };

  const handleItemSelect = useCallback(
    (item: Item) => {
      let itemToAdd: Item;

      if (item.isNew) {
        // If it's the "Add new" option, create a new item object
        itemToAdd = {
          id: getNewTagId(inputValue.trim()),
          name: inputValue.trim(),
          isNew: true, // Mark it as a genuinely new, selected tag
        };
      } else {
        itemToAdd = item;
      }

      // Ensure the item isn't already selected
      if (!selectedItems.some((selected) => selected.id === itemToAdd.id)) {
        const newSelectedItems = [...selectedItems, itemToAdd];
        setSelectedItems(newSelectedItems);
        onSelectedItemsChange(newSelectedItems);
      }

      setInputValue(""); // Clear input after selection
      inputRef.current?.focus(); // Keep focus on the input
    },
    [selectedItems, onSelectedItemsChange, inputValue, getNewTagId]
  );

  const handleRemoveItem = useCallback(
    (itemId: string) => {
      const newSelectedItems = selectedItems.filter((item) => item.id !== itemId);
      setSelectedItems(newSelectedItems);
      onSelectedItemsChange(newSelectedItems);
      inputRef.current?.focus(); // Keep focus on the input
    },
    [selectedItems, onSelectedItemsChange]
  );

  return (
    <div
      className="search-select-input-wrapper"
      ref={wrapperRef}
    >
      <div className="selected-items-container">
        {selectedItems.map((item) => (
          <div
            key={item.id}
            className={`selected-item-tag ${item.isNew ? "selected-item-tag--new" : ""}`}
          >
            <span>{item.name}</span>
            <button
              type="button"
              className="remove-item-button"
              onClick={() => handleRemoveItem(item.id)}
              aria-label={`Remove ${item.name}`}
            >
              <X size={16} />
            </button>
          </div>
        ))}
        <input
          ref={inputRef}
          type="text"
          className="search-input"
          placeholder={selectedItems.length === 0 ? placeholder : ""}
          value={inputValue}
          onChange={handleInputChange}
          onFocus={() => setIsDropdownOpen(true)}
          onKeyDown={(e) => {
            if (e.key === "Backspace" && inputValue === "" && selectedItems.length > 0) {
              // Remove last selected item on backspace if input is empty
              handleRemoveItem(selectedItems[selectedItems.length - 1].id);
            } else if (e.key === "Enter" && isDropdownOpen && filteredItems.length > 0) {
              // If Enter is pressed and dropdown is open, select the first item
              // This is typically the "Add new" option if it's visible, or the first filtered item
              e.preventDefault(); // Prevent form submission if applicable
              handleItemSelect(filteredItems[0]);
            }
          }}
        />
      </div>

      {isDropdownOpen && filteredItems.length > 0 && (
        <ul className="dropdown-list">
          {filteredItems.map((item) => (
            <li
              key={item.id}
              className={`dropdown-item ${item.isNew ? "dropdown-item--new" : ""}`}
              onClick={() => handleItemSelect(item)}
            >
              {item.name}
            </li>
          ))}
        </ul>
      )}

      {isDropdownOpen && filteredItems.length === 0 && inputValue.length > 0 && !allowNewTags && (
        <div className="no-results">No matching items found.</div>
      )}
      {/* If allowNewTags is true, "No matching items" message is replaced by the "Add new" option */}
    </div>
  );
};
