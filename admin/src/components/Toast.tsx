/* eslint-disable react-refresh/only-export-components, functional/no-let, functional/immutable-data, functional/prefer-immutable-types, react-hooks/globals */
import { createContext, useCallback, useContext, useState } from "react";

type ToastVariant = "success" | "error" | "info"

type Toast = Readonly<{
  id: number
  message: string
  variant: ToastVariant
}>

type ToastContextValue = Readonly<{
  addToast: (message: string, variant?: ToastVariant) => void
}>;

const ToastContext = createContext<ToastContextValue>({ addToast: () => {} });

const nextId = { value: 0 };
const inc = () => nextId.value++;

let globalAddToast: ((message: string, variant?: ToastVariant) => void) | null = null;

export function showToast(message: string, variant: ToastVariant = "error") {
  if (globalAddToast) {
    globalAddToast(message, variant);
  } else {
    console.error("[Toast]", message);
  }
}

const VARIANT_CLASSES: Record<ToastVariant, string> = {
  success: "bg-success/10 border-success text-success",
  error: "bg-danger/10 border-danger text-danger",
  info: "bg-primary/10 border-primary text-primary",
};

export function useToast() {
  return useContext(ToastContext);
}

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<ReadonlyArray<Toast>>([]);

  const addToast = useCallback((message: string, variant: ToastVariant = "info") => {
    const id = inc();
    setToasts((prev) => [...prev, { id, message, variant }] as unknown as ReadonlyArray<Toast>);
    setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id) as unknown as ReadonlyArray<Toast>);
    }, 4000);
  }, []);

  globalAddToast = addToast;

  return (
    <ToastContext.Provider value={{ addToast }}>
      {children}
      {toasts.length > 0 && (
        <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
          {toasts.map((t) => (
            <div key={t.id} className={`px-4 py-2 rounded border text-sm ${VARIANT_CLASSES[t.variant]}`}>
              {t.message}
            </div>
          ))}
        </div>
      )}
    </ToastContext.Provider>
  );
}
