import { useRef, forwardRef, useImperativeHandle, useCallback } from "react";
import { Button } from "../ui";

export type InputBarHandle = { focus: () => void };

type Props = {
  onSubmit: (cmd: string) => void
  history: string[]
  historyIndex: number
  setHistoryIndex: (i: number) => void
};

const InputBar = forwardRef<InputBarHandle, Props>(function InputBar({ onSubmit, history, historyIndex, setHistoryIndex }, ref) {
  const inputRef = useRef<HTMLInputElement>(null);

  useImperativeHandle(ref, () => ({
    focus: () => inputRef.current?.focus(),
  }));

  const handleSubmit = useCallback(() => {
    const val = inputRef.current?.value || "";
    if (val.trim()) {
      onSubmit(val);
      if (inputRef.current) inputRef.current.value = "";
    }
  }, [onSubmit]);

  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    const val = inputRef.current?.value || "";

    if (e.key === "Enter") {
      e.preventDefault();
      handleSubmit();
      return;
    }

    // Movement when input is empty
    if (val === "" && ["ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight"].includes(e.key)) {
      e.preventDefault();
      const dirMap: Record<string, string> = {
        ArrowUp: "n",
        ArrowDown: "s",
        ArrowLeft: "w",
        ArrowRight: "e",
      };
      onSubmit(dirMap[e.key]);
      return;
    }

    if (history.length === 0) return;
    if (e.key === "ArrowUp") {
      e.preventDefault();
      const newIndex = Math.min(historyIndex + 1, history.length - 1);
      setHistoryIndex(newIndex);
      if (inputRef.current) inputRef.current.value = history[newIndex] || "";
      return;
    }
    if (e.key === "ArrowDown") {
      e.preventDefault();
      const newIndex = Math.max(historyIndex - 1, -1);
      setHistoryIndex(newIndex);
      if (inputRef.current) {
        inputRef.current.value = newIndex >= 0 ? history[newIndex] : "";
      }
      return;
    }
  }, [handleSubmit, history, historyIndex, setHistoryIndex, onSubmit]);

  return (
    <div className="shrink-0 bg-surface border-t border-border px-3 py-2 flex gap-2 items-center">
      <span className="text-accent font-mono">&gt;</span>
      <input
        ref={inputRef}
        type="text"
        autoFocus
        autoComplete="off"
        autoCorrect="off"
        autoCapitalize="off"
        spellCheck={false}
        onKeyDown={handleKeyDown}
        className="flex-1 bg-surface text-foreground font-mono text-sm outline-none placeholder-muted"
        placeholder="Type a command..." />
      <Button variant="secondary" size="sm" onClick={handleSubmit}>SEND</Button>
    </div>
  );
});

export default InputBar;