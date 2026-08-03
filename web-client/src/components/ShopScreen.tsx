import { useState } from "react";
import type { ShopItem } from "../lib/types";

type Props = {
  shop: { shop_id: number; name: string; vendor_name: string; currency_name: string; gold: number; items: readonly ShopItem[] };
  onTransactionLog: (text: string) => void;
  onShopClose: () => void;
  onGoldChange: (gold: number) => void;
};

export default function ShopScreen({
  shop,
  onTransactionLog,
  onShopClose,
  onGoldChange,
}: Props) {
  void onTransactionLog;
  void onGoldChange;

  const [selectedItem, setSelectedItem] = useState<ShopItem | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [message, setMessage] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleBuy = async () => {
    if (!selectedItem || selectedItem.quantity <= 0) {
      setMessage("Item is out of stock!");
      return;
    }
    if (quantity > selectedItem.quantity) {
      setMessage(`Only ${selectedItem.quantity} available!`);
      return;
    }
    if (selectedItem.price * quantity > shop.gold) {
      setMessage("Not enough gold!");
      return;
    }
    setLoading(true);
    setMessage(null);
    // TODO: Call buyItem API and update shop state
    setLoading(false);
  };

  const handleSell = async () => {
    if (!selectedItem) {
      setMessage("Select an item to sell first!");
      return;
    }
    setLoading(true);
    setMessage(null);
    // TODO: Call sellItem API and update shop state
    setLoading(false);
  };

  const selectItem = (item: ShopItem) => {
    if (item.is_enabled && item.quantity > 0) {
      setSelectedItem(item);
      setQuantity(1);
      setMessage(null);
    } else {
      setMessage("This item is unavailable.");
    }
  };

  return (
    <div className="fixed inset-0 z-50 bg-background flex flex-col">
      {/* Header */}
      <header className="shrink-0 flex items-center justify-between px-4 py-3 border-b border-border bg-surface">
        <div className="flex items-center gap-2">
          <span className="font-bold text-accent text-sm">🏪 {shop.name}</span>
          <span className="text-xs text-muted">({shop.vendor_name})</span>
        </div>
        <button
          type="button"
          onClick={onShopClose}
          className="text-muted hover:text-foreground text-xl leading-none px-1"
          aria-label="Close shop"
        >
          &times;
        </button>
      </header>

      {/* Body */}
      <div className="flex-1 min-h-0 flex flex-col">
        {/* Top bar: Gold, Currency */}
        <div className="shrink-0 px-4 py-3 border-b border-border/50 bg-surface flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted uppercase tracking-wider">Gold:</span>
            <span className="font-mono font-bold text-accent text-sm">{shop.gold}</span>
            <span className="text-xs text-muted">{shop.currency_name}</span>
          </div>
          <div className="text-xs text-muted">
            {shop.items.filter((i) => i.is_enabled && i.quantity > 0).length} items available
          </div>
        </div>

        {/* Shop content area */}
        <div className="flex-1 min-h-0 flex overflow-hidden">
          {/* Items Grid */}
          <div className="flex-1 min-w-0 p-4 overflow-y-auto bg-surface/50">
            <div className="text-xs text-muted mb-2 uppercase tracking-wider">Items</div>
            <div className="grid grid-cols-2 gap-2">
              {shop.items.filter((i) => i.is_enabled).length === 0 ? (
                <div className="col-span-2 text-center py-8 text-muted text-sm italic">
                  No items available for sale.
                </div>
              ) : (
                shop.items.filter((i) => i.is_enabled).map((item) => (
                  <button
                    key={item.id}
                    type="button"
                    onClick={() => selectItem(item)}
                    disabled={!item.is_enabled}
                    className={`relative flex flex-col p-3 rounded border transition-all text-left ${
                      selectedItem?.id === item.id
                        ? "border-accent bg-accent/10"
                        : "border-border hover:border-accent/50"
                    } ${!item.is_enabled ? "opacity-50" : ""}`}
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-sm font-bold text-foreground truncate max-w-full">
                        {item.name}
                      </span>
                      <span className="text-xs font-mono text-success">
                        {item.price} {shop.currency_name}
                      </span>
                    </div>
                    {item.description && (
                      <p className="text-[10px] text-muted line-clamp-2 leading-tight">
                        {item.description}
                      </p>
                    )}
                    <div className="mt-2 flex items-center justify-between text-[10px]">
                      <span className="text-muted">
                        Qty: {item.quantity}
                      </span>
                      {item.category && (
                        <span className="text-xs text-accent/70 uppercase tracking-wide">
                          {item.category}
                        </span>
                      )}
                    </div>
                  </button>
                ))
              )}
            </div>
          </div>

          {/* Sidebar - Selected Item Details */}
          <div className="w-64 bg-surface border-l border-border p-4 overflow-y-auto">
            <div className="text-xs text-muted mb-2 uppercase tracking-wider">
              Selected Item
            </div>
            {selectedItem ? (
              <div className="space-y-4">
                <div>
                  <div className="font-bold text-foreground text-sm mb-1">
                    {selectedItem.name}
                  </div>
                  {selectedItem.description && (
                    <p className="text-[10px] text-muted mb-2">{selectedItem.description}</p>
                  )}
                </div>

                {selectedItem.category && (
                  <div className="text-xs text-accent uppercase tracking-wide">
                    Category: {selectedItem.category}
                  </div>
                )}

                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted">Price:</span>
                  <span className="font-mono text-success">{selectedItem.price} {shop.currency_name}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted">In Stock:</span>
                  <span className="font-mono">{selectedItem.quantity}</span>
                </div>

                {/* Quantity Input */}
                <div>
                  <label className="text-[10px] text-muted block mb-1">
                    Quantity
                  </label>
                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={() => setQuantity(Math.max(1, quantity - 1))}
                      disabled={quantity <= 1}
                      className="w-6 h-6 flex items-center justify-center rounded border border-border hover:bg-surface-alt disabled:opacity-50"
                    >
                      -
                    </button>
                    <input
                      type="number"
                      value={quantity}
                      onChange={(e) => setQuantity(Math.max(1, parseInt(e.target.value) || 1))}
                      className="flex-1 text-center text-sm font-mono bg-transparent border border-border rounded py-1"
                      min={1}
                      max={selectedItem.quantity}
                    />
                    <button
                      type="button"
                      onClick={() => setQuantity(Math.min(selectedItem.quantity, quantity + 1))}
                      disabled={quantity >= selectedItem.quantity}
                      className="w-6 h-6 flex items-center justify-center rounded border border-border hover:bg-surface-alt disabled:opacity-50"
                    >
                      +
                    </button>
                  </div>
                </div>

                {/* Buy/Sell Actions */}
                <div className="space-y-2">
                  <button
                    type="button"
                    onClick={handleBuy}
                    disabled={loading || !selectedItem.is_enabled}
                    className="w-full py-2 rounded bg-success hover:bg-success/80 text-background font-bold text-sm disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                  >
                    Buy for {selectedItem.price * quantity} {shop.currency_name}
                  </button>
                  <button
                    type="button"
                    onClick={handleSell}
                    disabled={loading}
                    className="w-full py-2 rounded bg-accent hover:bg-accent/80 text-background font-bold text-sm disabled:opacity-50 transition-colors"
                  >
                    Sell Item
                  </button>
                </div>
              </div>
            ) : (
              <div className="text-center py-8 text-muted text-sm italic">
                Select an item to view details
              </div>
            )}

            {/* Message display */}
            {message && (
              <div className={`mt-4 p-3 rounded text-[10px] font-mono ${message.startsWith("Success") ? "text-success" : "text-danger"}`}>
                {message}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
