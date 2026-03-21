import React, { useEffect, useState } from "react";
import { useEditor, EditorContent } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Link from "@tiptap/extension-link";
import {
  Bold,
  Italic,
  Strikethrough,
  List,
  ListOrdered,
  Undo,
  Redo,
  Heading1,
  Heading2,
  Heading3,
  Code,
  Link2,
  Unlink2,
} from "lucide-react";

import "./RichTextEditor.scss";

interface RichTextEditorProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
}

interface LinkDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (url: string, text?: string) => void;
  initialUrl?: string;
  initialText?: string;
  hasSelection: boolean;
}

const LinkDialog: React.FC<LinkDialogProps> = ({
  isOpen,
  onClose,
  onSubmit,
  initialUrl = "",
  initialText = "",
  hasSelection,
}) => {
  const [url, setUrl] = useState(initialUrl);
  const [text, setText] = useState(initialText);

  useEffect(() => {
    if (isOpen) {
      setUrl(initialUrl);
      setText(initialText);
    }
  }, [isOpen, initialUrl, initialText]);

  if (!isOpen) return null;

  const handleSubmit = () => {
    if (url.trim()) {
      onSubmit(url.trim(), hasSelection ? undefined : text.trim());
      onClose();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && url.trim() && (hasSelection || text.trim())) {
      e.preventDefault();
      handleSubmit();
    }
  };

  return (
    <div
      className="link-dialog-overlay"
      onClick={onClose}
    >
      <div
        className="link-dialog"
        onClick={(e) => e.stopPropagation()}
      >
        <h3>Insert Link</h3>
        <div>
          <div className="form-group">
            <label htmlFor="link-url">URL</label>
            <input
              id="link-url"
              type="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="https://example.com"
              autoFocus
            />
          </div>
          {!hasSelection && (
            <div className="form-group">
              <label htmlFor="link-text">Link Text</label>
              <input
                id="link-text"
                type="text"
                value={text}
                onChange={(e) => setText(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Enter link text"
              />
            </div>
          )}
          <div className="dialog-buttons">
            <button
              type="button"
              onClick={onClose}
              className="dialog-buttons-button dialog-buttons-button-secondary "
            >
              Cancel
            </button>
            <button
              type="button"
              onClick={handleSubmit}
              className="dialog-buttons-button dialog-buttons-button-primary"
            >
              Insert Link
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export const RichTextEditor: React.FC<RichTextEditorProps> = ({
  value,
  onChange,
  placeholder = "Start typing...",
  className = "",
  disabled = false,
}) => {
  const [showLinkDialog, setShowLinkDialog] = useState(false);
  const [hasSelection, setHasSelection] = useState(false);
  const [linkDialogData, setLinkDialogData] = useState({ url: "", text: "" });

  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        codeBlock: {
          HTMLAttributes: {
            class: "code-block",
          },
        },
      }),
      Link.configure({
        openOnClick: false,
      }),
    ],
    content: value,
    editable: !disabled,
    onUpdate: ({ editor }) => {
      onChange(editor.getHTML());
    },
    onCreate: ({ editor }) => {
      // Set initial content if provided
      if (value) {
        editor.commands.setContent(value);
      }
    },
    editorProps: {
      attributes: {
        class: "prose prose-sm focus:outline-none max-w-none",
        "data-placeholder": placeholder,
      },
    },
  });

  // Simple effect to update content when value changes externally
  useEffect(() => {
    if (editor && value !== editor.getHTML()) {
      const { from, to } = editor.state.selection;
      editor.commands.setContent(value, { emitUpdate: false });
      editor.commands.setTextSelection({ from, to });
    }
  }, [value, editor]);

  useEffect(() => {
    if (editor) {
      editor.setEditable(!disabled);
    }
  }, [disabled, editor]);

  const handleLinkClick = () => {
    if (!editor) return;

    const { from, to } = editor.state.selection;
    const hasSelectedText = from !== to;
    setHasSelection(hasSelectedText);
    // Get the current link attributes if cursor is on a link
    const previousUrl = editor.getAttributes("link").href || "";

    // Store the URL and selected text for the dialog in React state
    setLinkDialogData({
      url: previousUrl,
      text: hasSelectedText ? editor.state.doc.textBetween(from, to, "") : "",
    });

    setShowLinkDialog(true);
  };

  const handleUnlink = () => {
    if (!editor) return;
    editor.chain().focus().unsetLink().run();
  };

  const handleLinkSubmit = (url: string, text?: string) => {
    if (!editor) return;

    if (hasSelection) {
      // If text is selected, just add the link to it
      editor.chain().focus().setLink({ href: url }).run();
    } else {
      // If no text is selected, insert new text with link
      if (text) {
        editor
          .chain()
          .focus()
          .insertContent(`<a href="${url}" target="_blank" rel="noopener noreferrer nofollow">${text}</a>`)
          .run();
      }
    }
  };

  if (!editor) {
    return <div>Loading editor...</div>;
  }

  const isLinkActive = editor.isActive("link");

  return (
    <>
      <div className={`richtext-editor ${className}`}>
        <div className="toolbar">
          {/* Headings */}
          <button
            type="button"
            onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()}
            className={editor.isActive("heading", { level: 1 }) ? "active" : ""}
            disabled={disabled}
            title="Heading 1"
          >
            <Heading1 size={16} />
          </button>

          <button
            type="button"
            onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
            className={editor.isActive("heading", { level: 2 }) ? "active" : ""}
            disabled={disabled}
            title="Heading 2"
          >
            <Heading2 size={16} />
          </button>

          <button
            type="button"
            onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
            className={editor.isActive("heading", { level: 3 }) ? "active" : ""}
            disabled={disabled}
            title="Heading 3"
          >
            <Heading3 size={16} />
          </button>

          <div className="separator" />

          {/* Text formatting */}
          <button
            type="button"
            onClick={() => editor.chain().focus().toggleBold().run()}
            className={editor.isActive("bold") ? "active" : ""}
            disabled={disabled}
            title="Bold"
          >
            <Bold size={16} />
          </button>

          <button
            type="button"
            onClick={() => editor.chain().focus().toggleItalic().run()}
            className={editor.isActive("italic") ? "active" : ""}
            disabled={disabled}
            title="Italic"
          >
            <Italic size={16} />
          </button>

          <button
            type="button"
            onClick={() => editor.chain().focus().toggleStrike().run()}
            className={editor.isActive("strike") ? "active" : ""}
            disabled={disabled}
            title="Strikethrough"
          >
            <Strikethrough size={16} />
          </button>

          <div className="separator" />

          {/* Lists */}
          <button
            type="button"
            onClick={() => editor.chain().focus().toggleBulletList().run()}
            className={editor.isActive("bulletList") ? "active" : ""}
            disabled={disabled}
            title="Bullet List"
          >
            <List size={16} />
          </button>

          <button
            type="button"
            onClick={() => editor.chain().focus().toggleOrderedList().run()}
            className={editor.isActive("orderedList") ? "active" : ""}
            disabled={disabled}
            title="Numbered List"
          >
            <ListOrdered size={16} />
          </button>

          <div className="separator" />

          {/* Code Block */}
          <button
            type="button"
            onClick={() => editor.chain().focus().toggleCodeBlock().run()}
            className={editor.isActive("codeBlock") ? "active" : ""}
            disabled={disabled}
            title="Code Block"
          >
            <Code size={16} />
          </button>

          <div className="separator" />

          {/* Link buttons */}
          <button
            type="button"
            onClick={handleLinkClick}
            className={isLinkActive ? "active" : ""}
            disabled={disabled}
            title="Insert Link"
          >
            <Link2 size={16} />
          </button>

          <button
            type="button"
            onClick={handleUnlink}
            disabled={disabled || !isLinkActive}
            title="Remove Link"
          >
            <Unlink2 size={16} />
          </button>

          <div className="separator" />

          {/* History */}
          <button
            type="button"
            onClick={() => editor.chain().focus().undo().run()}
            disabled={disabled || !editor.can().undo()}
            title="Undo"
          >
            <Undo size={16} />
          </button>

          <button
            type="button"
            onClick={() => editor.chain().focus().redo().run()}
            disabled={disabled || !editor.can().redo()}
            title="Redo"
          >
            <Redo size={16} />
          </button>
        </div>

        <EditorContent
          editor={editor}
          className="editor-content"
        />
      </div>

      <LinkDialog
        isOpen={showLinkDialog}
        onClose={() => setShowLinkDialog(false)}
        onSubmit={handleLinkSubmit}
        hasSelection={hasSelection}
        initialUrl={linkDialogData.url}
        initialText={linkDialogData.text}
      />
    </>
  );
};
