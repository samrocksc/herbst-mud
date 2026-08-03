import { useCallback, useEffect, useState } from "react";
import type { ShopItem } from "../lib/types";
import { getShopItems, buyItem, sellItem } from "../lib/api";

type Props = {
  shopID: number;
  vendorName: string;
  currencyName: string;
  initialGold: number;
  onTransactionLog: (text: string) => void;
  onShopClose: () => void;
  onGoldChange: (gold: number) => void;
};

export function useShopEngine({
  shopID,
  vendorName,
  currencyName,
  initialGold,
  onTransactionLog,
  onShopClose,
  onGoldChange,
}: Props) {
  const [shopItems, setShopItems] = useState<readonly ShopItem[]>([]);
  const [gold, setGold] = useState(initialGold);
  const [loading, setLoading] = useState(false);
  const [selectedItem, setSelectedItem] = useState<ShopItem | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [message, setMessage] = useState<string | null>(null);

  const fetchItems = useCallback(async () => {
    setLoading(true);
    try {
      const items = await getShopItems(shopID);
      setShopItems(items);
      setMessage(null);
    } catch (err) {
      setMessage(`Failed to load shop items: ${err instanceof Error ? err.message : "Unknown error"}`);
    } finally {
      setLoading(false);
    }
  }, [shopID]);

  const buySelected = useCallback(async () => {
    if (!selectedItem || selectedItem.quantity <= 0) {
      setMessage("Item is out of stock!");
      return;
    }
    if (quantity > selectedItem.quantity) {
      setMessage(`Only ${selectedItem.quantity} available!`);
      return;
    }
    if (selectedItem.price * quantity > gold) {
      setMessage("Not enough gold!");
      return;
    }

    setLoading(true);
    try {
      const result = await buyItem(shopID, selectedItem.id, quantity);
      setMessage(result.message);
      if (result.success) {
        setGold(result.gold);
        onTransactionLog(`Bought ${selectedItem.name} x${quantity}`);
        await fetchItems();
      }
    } catch (err) {
      setMessage(`Purchase failed: ${err instanceof Error ? err.message : "Unknown error"}`);
    } finally {
      setLoading(false);
    }
  }, [selectedItem, quantity, gold, shopID, onTransactionLog, fetchItems]);

  const sellSelected = useCallback(async () => {
    if (!selectedItem) {
      setMessage("Select an item to sell first!");
      return;
    }
    if (quantity <= 0) {
      setMessage("Quantity must be at least 1!");
      return;
    }

    setLoading(true);
    try {
      const result = await sellItem(shopID, selectedItem.id, quantity);
      setMessage(result.message);
      if (result.success) {
        setGold(result.gold);
        onTransactionLog(`Sold ${selectedItem.name} x${quantity}`);
        await fetchItems();
      }
    } catch (err) {
      setMessage(`Sale failed: ${err instanceof Error ? err.message : "Unknown error"}`);
    } finally {
      setLoading(false);
    }
  }, [selectedItem, quantity, shopID, onTransactionLog, fetchItems]);

  const selectItem = useCallback(
    (item: ShopItem) => {
      if (item.is_enabled && item.quantity > 0) {
        setSelectedItem(item);
        setQuantity(1);
        setMessage(null);
      } else {
        setMessage("This item is unavailable.");
      }
    },
    []
  );

  const refreshShop = useCallback(() => {
    setMessage(null);
    fetchItems();
  }, [fetchItems]);

  // Load items on mount
  useEffect(() => {
    fetchItems();
  }, [fetchItems]);

  useEffect(() => {
    onGoldChange(gold);
  }, [gold, onGoldChange]);

  return {
    shopItems,
    gold,
    currencyName,
    vendorName,
    loading,
    selectedItem,
    quantity,
    message,
    selectItem,
    setQuantity,
    buySelected,
    sellSelected,
    refreshShop,
    fetchItems,
    onShopClose,
  };
}
