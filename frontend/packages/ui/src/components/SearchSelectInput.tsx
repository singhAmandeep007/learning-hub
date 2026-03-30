import { useState, useRef, useEffect, useCallback } from "react";
import { X } from "lucide-react";

import { withPrefix } from "../constants";

export interface SearchSelectItem {
  id: string;
  name: string;
  isNew?: boolean;
}

export interface SearchSelectInputProps {
  items: SearchSelectItem[];
  onSelectedItemsChange: (selectedItems: SearchSelectItem[]) => void;
  initialSelectedItems?: SearchSelectItem[];
  placeholder?: string;
  allowNewItems?: boolean; // allow adding new tags that are not in the list
}

const wrapperClass = withPrefix("search-select-input");

export function SearchSelectInput({
  items,
  onSelectedItemsChange,
  initialSelectedItems = [],
  placeholder = "Search and select items...",
  allowNewItems = false,
}: SearchSelectInputProps) {
  const [inputValue, setInputValue] = useState<string>("");
  const [selectedItems, setSelectedItems] = useState<SearchSelectItem[]>(initialSelectedItems);
  const [filteredItems, setFilteredItems] = useState<SearchSelectItem[]>([]);
  const [isDropdownOpen, setIsDropdownOpen] = useState<boolean>(false);

  const wrapperRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const listId = `${wrapperClass}-listbox`;

  const getNewItemId = useCallback(
    (itemName: string) => `new-${itemName.toLowerCase().replace(/\s+/g, "-")}-${Date.now()}`,
    []
  );

  // filter items based on input value
  useEffect(() => {
    const normalizedInput = inputValue.toLowerCase().trim();
    let newFilteredItems: SearchSelectItem[] = [];

    if (inputValue.length > 0) {
      // filter existing items based on search input
      newFilteredItems = items.filter(
        (item) =>
          item.name.toLowerCase().includes(normalizedInput) &&
          !selectedItems.some((selected) => selected.id === item.id)
      );

      // check if input value exactly matches any existing (non-selected) item
      const exactMatchExists = items.some(
        (item) =>
          item.name.toLowerCase() === normalizedInput && !selectedItems.some((selected) => selected.id === item.id)
      );
      // Check if input value exactly matches any already selected item (to avoid duplicate 'add new' option)
      const alreadySelectedExactMatch = selectedItems.some((item) => item.name.toLowerCase() === normalizedInput);

      // if allowNewTags is true AND no exact match found among all items (available or selected)
      // AND the input is not just whitespace
      // Ensure not just whitespace
      if (allowNewItems && inputValue.trim().length > 0 && !exactMatchExists && !alreadySelectedExactMatch) {
        // add a "create new" option if no existing item matches exactly
        // and ensure it's not a duplicate of an already selected 'new' tag
        const isNewItemAlreadySelected = selectedItems.some(
          (item) => item.isNew && item.name.toLowerCase() === normalizedInput
        );

        if (!isNewItemAlreadySelected) {
          newFilteredItems.unshift({
            id: getNewItemId(inputValue.trim()), // use a temporary ID
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
  }, [inputValue, items, selectedItems, allowNewItems, getNewItemId]);

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

  const handleItemSelect = useCallback(
    (item: SearchSelectItem) => {
      let itemToAdd: SearchSelectItem;

      if (item.isNew) {
        // If it's the "Add new" option, create a new item object
        itemToAdd = {
          id: getNewItemId(inputValue.trim()),
          name: inputValue.trim(),
          isNew: true, // Mark it as a genuinely new, selected
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
    [selectedItems, onSelectedItemsChange, inputValue, getNewItemId]
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
      className={`${wrapperClass}__wrapper`}
      ref={wrapperRef}
    >
      <div className={`${wrapperClass}__selected-items`}>
        {selectedItems.map((item) => (
          <div
            key={item.id}
            className={`${wrapperClass}__selected-item-tag ${item.isNew ? `${wrapperClass}__selected-item-tag--new` : ""}`}
          >
            <span>{item.name}</span>
            <button
              type="button"
              className={`${wrapperClass}__remove-item-button`}
              onClick={() => handleRemoveItem(item.id)}
              aria-label={`Remove ${item.name}`}
            >
              <X size={12} />
            </button>
          </div>
        ))}

        <input
          ref={inputRef}
          type="text"
          className={`${wrapperClass}__search-input`}
          role="combobox"
          aria-expanded={isDropdownOpen}
          aria-controls={listId}
          aria-autocomplete="list"
          placeholder={selectedItems.length === 0 ? placeholder : ""}
          value={inputValue}
          onChange={(event) => setInputValue(event.target.value)}
          onFocus={() => setIsDropdownOpen(true)}
          onKeyDown={(event) => {
            if (event.key === "Backspace" && inputValue === "" && selectedItems.length > 0) {
              // Remove last selected item on backspace if input is empty
              handleRemoveItem(selectedItems[selectedItems.length - 1].id);
            } else if (event.key === "Enter" && isDropdownOpen && filteredItems.length > 0) {
              // If Enter is pressed and dropdown is open, select the first item
              // This is typically the "Add new" option if it's visible, or the first filtered item
              event.preventDefault(); // Prevent form submission if applicable
              handleItemSelect(filteredItems[0]);
            }
          }}
        />
      </div>

      {isDropdownOpen && filteredItems.length > 0 && (
        <ul
          className={`${wrapperClass}__dropdown-list`}
          role="listbox"
          id={listId}
        >
          {filteredItems.map((item) => (
            <li
              key={item.id}
              className={`${wrapperClass}__dropdown-item ${item.isNew ? `${wrapperClass}__dropdown-item--new` : ""}`}
              role="option"
              onMouseDown={(event) => event.preventDefault()}
              onClick={() => handleItemSelect(item)}
            >
              {item.name}
            </li>
          ))}
        </ul>
      )}

      {isDropdownOpen && filteredItems.length === 0 && inputValue.length > 0 && !allowNewItems && (
        <div className={`${wrapperClass}__no-results`}>No matching items found.</div>
      )}
      {/* If allowNewTags is true, "No matching items" message is replaced by the "Add new" option */}
    </div>
  );
}
